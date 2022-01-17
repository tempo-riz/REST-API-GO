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

    docker build --tag appsec .
## Run image
    docker run --publish 8080:8080 appsec

Now app is running inside of container and you can still test it with ports "published"

<br><br>
# Part 2 - TLS Protection

## Requirements

    nginx, docker-compose 

## Goal

Implementing a reverse proxy that will terminate the TLS
connection.

Using nginx as a http gateway (
redirect traffic on
HTTP port 80 to the HTTP listener of the application )

nginx.conf is updated with our subdomain, key and certificate 

openssl was used to create those key and certificate

## Run  in docker compose

    docker-compose up 

if we try to acces http://0.0.0.0/students via browser we can see 

that it's getting redirected to it's secured version : https://0.0.0.0/students

I used the following ciphers to harden the tls config :

    ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-SHA384;

TLS analyzer

I failed to compile the library with maven getting :

    [INFO] BUILD FAILURE
    [ERROR] Failed to execute goal com.mycila:license-maven-plugin:4.2.rc2:format (default) on project TLS-Scanner
<br><br>

# Partie 3 - Authentication

Here we want to add authentication and authorization on our endpoints using basic authentication and Oauth2.0/OIDC

## Adding Authorization Logic

- foo can only perform GET requests 
- aristote can perform all HTTP requests

I implemented the logic this way :

    if c.Request.Method != "GET" && user != "aristote" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}
<br>

## Implementingan OIDC Consumerusing Okta

I finally connected postman to okta, using a fresh generated token, and could test and validate requests with postman

<br>

## Implement authorization

For this part I used a hashmap to create authorization logic depending on who is sending a request, sending correct HTTP responses codes for denied requests :

    if val, ok := m[user]; ok { //if dico contains key (user)
        for _, v := range val {
            if v == c.Request.Method { //if they have correct acces right for query
                c.Next() //continue routing 
            }
        }
		//if they dont
        c.AbortWithStatus(http.StatusForbidden)

<br><br>

# Part 4 - SecretManagement