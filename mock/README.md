## Build and run with docker

`docker build . -t blind/mock`

`docker run -p 8001:8001 --name blind_mock blind/mock`

## Run with golang

`go run main.go`

you can watch index.html to see how to upload files

### Routes
http://localhost:8001/api/skin/evaluate/image   param: "file": file<br>
http://localhost:8001/api/lunghj/evaluate/image param: "file": file<br>
http://localhost:8001/api/fall/evaluate/video   param: "file": file<br>

### Response format
{"Probability":60.466026,"ProcessingTime":"4.094966ms"}
