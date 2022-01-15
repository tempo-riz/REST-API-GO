package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

type student struct {
	ID       string `json:"id"`
	Lastname string `json:"Lastname"`
	Name     string `json:"Name"`
	Filiere  string `json:"filiere"`
}

type teacher struct {
	ID       string `json:"id"`
	Lastname string `json:"Lastname"`
	Name     string `json:"Name"`
	Class    string `json:"class"`
}

var students = []student{
	{ID: "1", Lastname: "Montandon", Name: "Philippe", Filiere: "ISC"},
	{ID: "2", Lastname: "Chatillon", Name: "Thibault", Filiere: "ISC"},
}

var teachers = []teacher{
	{ID: "1", Lastname: "Pfeiffer", Name: "Ludovic", Class: "A401"},
	{ID: "2", Lastname: "Ouafi", Name: "Khaled", Class: "A402"},
}

// getStudents responds with the list of all students as JSON.
func getStudents(c *gin.Context) {
	if !check_auth(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	c.IndentedJSON(http.StatusOK, students)
}

// postStudents adds an student from JSON received in the request body.
func postStudents(c *gin.Context) {
	if !check_auth(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	var newStudent student

	// Call BindJSON to bind the received JSON to
	// newStudent.
	if err := c.BindJSON(&newStudent); err != nil {
		fmt.Println("error post student json ?")
		return
	}

	// Add the new student to the slice.
	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
}

// getStudentByID locates the student whose ID value matches the id
// parameter sent by the client, then returns that student as a response.
func getStudentByID(c *gin.Context) {
	if !check_auth(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	id := c.Param("id")

	// Loop through the list of students, looking for
	// an student whose ID value matches the parameter.
	for _, a := range students {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

func deleteStudentByID(c *gin.Context) {
	if !check_auth(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}
	id := c.Param("id")

	// Loop through the list of students, looking for
	// an student whose ID value matches the parameter.
	for i := 0; i < len(students); i++ {
		if students[i].ID == id {
			c.IndentedJSON(http.StatusOK, students[i])
			copy(students[i:], students[i+1:])    // Shift a[i+1:] left one index.
			students = students[:len(students)-1] // Truncate slice.

			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

// getTeachers responds with the list of all teachers as JSON.
func getTeachers(c *gin.Context) {
	if !verify(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	c.IndentedJSON(http.StatusOK, teachers)
}

// postTeachers adds an teacher from JSON received in the request body.
func postTeachers(c *gin.Context) {
	if !verify(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	var newTeacher teacher

	// Call BindJSON to bind the received JSON to
	// newTeacher.
	if err := c.BindJSON(&newTeacher); err != nil {
		return
	}

	// Add the new teacher to the slice.
	teachers = append(teachers, newTeacher)
	c.IndentedJSON(http.StatusCreated, newTeacher)
}

// getTeacherByID locates the teacher whose ID value matches the id
// parameter sent by the client, then returns that teacher as a response.
func getTeacherByID(c *gin.Context) {
	if !verify(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	id := c.Param("id")

	// Loop through the list of teachers, looking for
	// an teacher whose ID value matches the parameter.
	for _, a := range teachers {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "teacher not found"})
}

func deleteTeacherByID(c *gin.Context) {
	if !verify(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "acces denied"})
		return
	}

	id := c.Param("id")

	// Loop through the list of teachers, looking for
	// an teacher whose ID value matches the parameter.
	for i := 0; i < len(teachers); i++ {
		if teachers[i].ID == id {
			c.IndentedJSON(http.StatusOK, teachers[i])
			copy(teachers[i:], teachers[i+1:])    // Shift a[i+1:] left one index.
			teachers = teachers[:len(teachers)-1] // Truncate slice.

			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "teacher not found"})
}

func check_auth(c *gin.Context) bool {
	user, _, _ := c.Request.BasicAuth()
	if c.Request.Method != "GET" && user != "aristote" {
		return false
	}
	return true
}

var toValidate = map[string]string{
	"aud": "api://default",
	"cid": os.Getenv("0oa3lp6i6zXA2yMvp5d7"),
}

func verify(c *gin.Context) bool {
	status := true
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		verifierSetup := jwtverifier.JwtVerifier{
			Issuer:           "https://" + os.Getenv("dev-58917141.okta.com") + "/oauth2/default",
			ClaimsToValidate: toValidate,
		}
		verifier := verifierSetup.New()
		_, err := verifier.VerifyAccessToken(token)
		if err != nil {
			c.String(http.StatusForbidden, err.Error())
			print(err.Error())
			status = false
		}
	} else {
		c.String(http.StatusUnauthorized, "Unauthorized")
		status = false
	}
	return status
}

func main() {
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{"foo": "bar", "aristote": "Eucl1de"}))

	authorized.GET("/students", getStudents)
	authorized.GET("/students/:id", getStudentByID)
	authorized.GET("/teachers", getTeachers)
	authorized.GET("/teachers/:id", getTeacherByID)

	authorized.POST("/students", postStudents)
	authorized.POST("/teachers", postTeachers)

	authorized.DELETE("/students/:id", deleteStudentByID)
	authorized.DELETE("/teachers/:id", deleteTeacherByID)

	router.Run("0.0.0.0:8080")

}

// curl -X POST -H "Content-Type: application/json" \
//     -d '{ID: "1", Lastname: "Pfeiffer", Name: "Ludovic", Class: "A401"}' \
//     http://aristote:Eucl1de@0.0.0.0:8080/students

//token 00FOcy87-zhZJqUVGLztyy740UlnIUwZBf_aOQg0x8
