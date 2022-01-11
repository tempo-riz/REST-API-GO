# Part 1 - Building a RESTful API in Golang


All REST endpoints are successfully avalaible for students and teachers as follow at 

0.0.0.0:8080 reachable using HTTP : 

- /students :

        GET – Get a list of all students,returned as JSON
    
        POST – Add a new student from request data sent as JSON.
    
- /students/:id

        GET – Get a student by its ID, returning the student data as JSON.

        DELETE – Delete a student by its ID.

- /teachers

        GET – Get a list of all teachers, returned as JSON.

        POST – Add a new teacher from request data sent as JSON.
- /teachers/:id

        GET – Get a teacher by its ID, returning the student data as JSON.

        DELETE – Delete a teacher by its ID.

<br>

## How to run

with this simple line in app's directory

    go run . 

## How to test

for example :

    curl http://0.0.0.0:8080/students


## Requirements

    golang, docker, gin-gonic framework

<br>

# Containerization

## Build docker image 

    docker build --tag app-sec .
## Run image
    docker run --publish 8080:8080 app-sec

Now app is running inside of container and you can still test it with ports "published"


# Part2 - TLS Protection

## Requirements

    nginx, docker-compose 

## Goal

Implementing a reverse proxy that will terminate the TLS
connection.

Using nginx as a http gateway (
redirect traffic on
HTTP port 80 to the HTTP listener of the application )

nginx.conf is updated with our subdomain, key and certificate 

## Run all of it in docker compose

...