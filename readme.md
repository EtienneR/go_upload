# Upload & listing files


"goupload" is a simple manager images files with these functionalities:

- upload a file
	- only "gif", "jpg" & "png" extensions
	- limitation size to 2 Mb
- display file informations
	- name
	- size
	- date of last modifications
	- dimensions: width & height
- delete a file


## Classic installation

1. Install: `go get github.com/EtienneR/go_upload`
2. Launch the server since the folder application: `go run main.go`
3. Access to the server with the URL address: 127.0.0.1:3000


## Docker installation

1. Put you in the application folder
2. Build the image: `sudo docker build -t go_upload .`
3. Run the application in a temporarly container: `sudo docker run --publish 3000:3000 --name go_upload_test --rm go_upload`
4. Access to the server with the URL address: 127.0.0.1:3000