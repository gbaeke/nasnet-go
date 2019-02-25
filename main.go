package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/mholt/certmagic"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type ResultPageData struct {
	Probability float32
	MaxIndex    int
	Category    string
	Picture     string
}

var model *tf.SavedModel
var categories map[int][]string

func init() {
	// load tensorflow model from disk
	var err error
	model, err = tf.LoadSavedModel("nasnet",
		[]string{"atag"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// get the categories
	categories, err = getCategories()
	if err != nil {
		log.Fatal("Error getting categories", err)
	}

}

func main() {
	log.Println("TensorFlow version: ", tf.Version())

	//check environment variable to enable SSL
	sslEnabled := getEnv("ssl", "false")
	hostName := getEnv("hostname", "")
	if hostName == "" && sslEnabled == "true" {
		log.Fatalln("Specify hostname environment variable when SSL is on")
	}

	if sslEnabled == "true" {
		log.Println("SSL is on")
		// certmagic
		certmagic.Agreed = true
		certmagic.Email = "mail@mail.com"
		certmagic.CA = certmagic.LetsEncryptStagingCA

		mux := http.NewServeMux()
		mux.HandleFunc("/", upload)

		err := certmagic.HTTPS([]string{hostName}, mux)
		if err != nil {
			log.Println(err)
		}
	} else {
		http.HandleFunc("/", upload)
		http.ListenAndServe(":9090", nil)
	}

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func getTensor(m image.Image) (*tf.Tensor, error) {
	var BCHW [1][224][224][3]float32

	//get bounds of image 0 0 224 224
	bounds := m.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// height = y and width = x
			BCHW[0][y][x][0] = convertColor(r)
			BCHW[0][y][x][1] = convertColor(g)
			BCHW[0][y][x][2] = convertColor(b)
		}
	}

	return tf.NewTensor(BCHW)
}

func convertColor(value uint32) float32 {
	return (float32(value>>8) - float32(127.5)) / float32(127.5)
}

func readImage(imgPath string, width, height int) (image.Image, error) {
	// read the image file
	reader, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// decode the image
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// resize image
	m = imaging.Resize(m, width, height, imaging.Linear)

	return m, nil
}

func getCategories() (map[int][]string, error) {
	// open categories file
	reader, err := os.Open("imagenet_class_index.json")
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// read JSON categories
	catJSON, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// unmarshal into map of int to array of string
	var categories map[int][]string
	err = json.Unmarshal(catJSON, &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// GET
		t, _ := template.ParseFiles("upload.gtpl")

		t.Execute(w, nil)

	} else if r.Method == "POST" {
		// Post
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		//we want to save the file
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

		//not efficient because we read again
		m, err := readImage("./test/"+handler.Filename, 224, 224)
		if err != nil {
			log.Fatal("Cannot read image")
		}

		//get tensor from image
		tensor, err := getTensor(m)
		if err != nil {
			log.Fatal("Cannot get tensor")
		}

		//run a session
		output, err := model.Session.Run(
			map[tf.Output]*tf.Tensor{
				model.Graph.Operation("input_1").Output(0): tensor,
			},
			[]tf.Output{
				model.Graph.Operation("predictions/Softmax").Output(0),
			},
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		predictions, ok := output[0].Value().([][]float32)

		if !ok {
			log.Fatal(fmt.Sprintf("output has unexpected type %T", output[0].Value()))
		}

		// highest result
		maxProb := float32(0.0)
		maxIndex := 0
		for index, prob := range predictions[0] {
			if prob > maxProb {
				maxProb = prob
				maxIndex = index
			}
		}

		//base64 encode image data to show image in response
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, m, nil)
		encodedImage := base64.StdEncoding.EncodeToString(buf.Bytes())

		//response data
		data := ResultPageData{
			Probability: maxProb,
			MaxIndex:    maxIndex,
			Category:    categories[maxIndex][1],
			Picture:     encodedImage,
		}

		//respond with template
		tmpl := template.Must(template.ParseFiles("response.gtpl"))
		tmpl.Execute(w, data)

	} else {
		fmt.Println("Unknown HTTP " + r.Method + "  Method")
	}
}
