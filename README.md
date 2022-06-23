# delay-server

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/nexus-uw/delay-server)

## purpose
super simple http server that will allow you 'mock' delays in HTTP requires (ie: dont know what the request will be, just that you need to have a little delay and a real connection)

## running (only works on amd64 docker arch)
First will be to build an image from the Dockerfile included by running:
```
docker build -t delay-server .
```
On successful build, you can go ahead to create a running container from the image
```
docker run -p 8080:8080 delay-server
```
then vist [localhost:8080](http://localhost:8080)


