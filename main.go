package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"

	"github.com/disintegration/imaging"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	imagePath := flag.String("image", "", "path to image")
	flag.Parse()

	if *imagePath == "" {
		log.Fatal("Please specify image path. Use --help for help.")
	}

	fmt.Println("TensorFlow version: ", tf.Version())

	//turn off logging

	// load tensorflow model from disk
	model, err := tf.LoadSavedModel("nasnet",
		[]string{"atag"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// read cat image
	m, err := readImage(*imagePath, 224, 224)
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

	// get the categories
	categories, err := getCategories()
	if err != nil {
		log.Fatal("Error getting categories", err)
	}

	fmt.Println("Highest prob is", maxProb, "at", maxIndex)
	fmt.Println("Probably ", categories[maxIndex])

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
