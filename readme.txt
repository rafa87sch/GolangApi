#runnning the api alone
go run main.go

#there are two post will show a string
/api/encrypt
/api/decrypt


example
/api/encrypt
{
    "password":"123"
}

/api/decrypt
/api/decrypt
{
    "password":"4521ef5e25909739a43938d60064ebf6159e3ab90183d9c5d859435f599cf4"
}

#run the test
go test -v

#dockers
docker build -t my-app .

#run with the port 5000
docker run -d -p 5000:5000 my-app


#stop container
docker stop <container id>
