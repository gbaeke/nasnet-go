# nasnet-go

Web app to classify images with nasnetmobile

Use the Docker image gbaeke/nasnet to try it:

docker run -p 80:9090 -d gbaeke/nasnet

Afterwards, point a browser at http://localhost and try it out

If you want to build and run the app locally, you will need Linux or MacOS. Check the Dockerfile for prerequisites like the TensorFlow C API.

Note: to enable SSL set environment variable ssl equal to true (lowercase) and hostname to the hostname for the certificate. Make sure that the hostname resolves to the container on port 80! The code in main.go uses Let's Encrypt staging CAs

Note: if you deploy to Azure Container Instances, remember that the certificate is stored in ephemeral storage. You might hit Let's Encrypt rate limits if you switch to the production CAs and create the container many times. One solution is to mount an Azure File Share on  /root/.local/share/certmagic

Note: with the release of Google Cloud Run, the code was updated to use the PORT environment variable that Cloud Run injects.