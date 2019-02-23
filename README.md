# nasnet-go

Web app to classify images with nasnetmobile

Use the Docker image gbaeke/nasnet to try it:

docker run -p 80:9090 -d gbaeke/nasnet

Afterwards, point a browser at http://localhost and try it out

If you want to build and run the app locally, you will need Linux or MacOS. Check the Dockerfile for prerequisites like the TensorFlow C API.