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

if we try to acces http://foo:bar@0.0.0.0/ via browser we can see 

that it's getting redirected to it's secured version : https://0.0.0.0/students

I used the following ciphers to harden the tls config :

    ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-SHA384;

TLS analyzer

I failed to compile the library with maven getting :

    [INFO] BUILD FAILURE
    [ERROR] Failed to execute goal com.mycila:license-maven-plugin:4.2.rc2:format (default) on project TLS-Scanner
after some unsuccesfulls attemps to fix it i decided to continue with part 3

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

I finally connected postman to okta, using a fresh generated token, and could test and validate requests with postman,
from here i'm not sure how to link existing API with okta and postman...


So i created logic based on existing user in db like this :

- foo can only perform POST and GET queries

- Aristote can only perform DELETE and GET queries 

- All other users can only perform GET queries

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

Until now user credentials were hardcoded into our code, lets use environment variables instead !

Using os package we can easily set and get environnement variables :        

    os.getEnv("key")

In order to make it work in docker
I created the file variables.env containting the key-values to inject in docker-compose, also specifying in docker-compose.yml :

    env_file:
    - variables.env
This is what's inside it in format key=value

    POST=foo user1
    DELETE=aristote
    USERS= foo:bar artistote:Eucl1de user1:pass1

Now in our main.go we can retrieve them, and convert those values into a Map splice

    getSpliceFromEnv("KEY")

And then, for users we can parse them into a map:

    users := make(map[string]string)
	for _, e := range getSpliceFromEnv("USERS") {
		if i := strings.Index(e, ":"); i >= 0 {
			users[e[:i]] = e[i+1:]
		}
	}

Now we can use our creditentials without them beeing readable by anyone accessing our code


Before :

    gin.BasicAuth(gin.Accounts{"foo": "bar", "aristote": "Eucl1de"})
    
Now : 

    authorized := router.Group("/", gin.BasicAuth(users))

Same concept is used for permission logic in check_student_authorization function :

I retrieve value containing a list with allowed users, check if current user is in it :

    switch c.Request.Method {
	case "POST":
		if contains(getSpliceFromEnv("POST"), user) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
    ...

<br><br>

# Conlusion

So this is it, we have a (mostly) functional and secured Rest API !

Sadly I couldn't succeed with okta part but the rest seems to work so far.

### What I learned :
- golang
- how an api is really working, how to create one
- A lot of security concepts, how to apply them incrementally 
- purpose and utiliy of docker and nginx
