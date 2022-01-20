package main

import (
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

func getStudents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, students)
}

func postStudents(c *gin.Context) {

	var newStudent student

	if err := c.BindJSON(&newStudent); err != nil {
		print(err.Error())
		return
	}

	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
}

func getStudentByID(c *gin.Context) {

	id := c.Param("id")

	for _, a := range students {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

func deleteStudentByID(c *gin.Context) {
	id := c.Param("id")

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

func getTeachers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teachers)
}

func postTeachers(c *gin.Context) {

	var newTeacher teacher

	if err := c.BindJSON(&newTeacher); err != nil {
		return
	}

	teachers = append(teachers, newTeacher)
	c.IndentedJSON(http.StatusCreated, newTeacher)
}

func getTeacherByID(c *gin.Context) {

	id := c.Param("id")

	for _, a := range teachers {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "teacher not found"})
}

func deleteTeacherByID(c *gin.Context) {

	id := c.Param("id")

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

// - foo can only perform GET requests
// - aristote can perform all HTTP requests
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

func check_teacher_authentication(c *gin.Context) {
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

	if status {
		c.Next() //continue routing
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func getSpliceFromEnv(key string) []string { //return [user1, user2] from env key=user1 user2
	return strings.Split(os.Getenv(key), " ")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func check_student_authorization(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string) // get username

	post := getSpliceFromEnv("POST")
	delete := getSpliceFromEnv("DELETE")

	switch c.Request.Method {
	case "POST":
		if contains(post, user) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	case "DELETE":
		if contains(delete, user) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	case "GET":
		c.Next() //all authentified users can use GET
	default:
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func main() {
	//retrieve logins from env variable
	users := make(map[string]string)

	for _, e := range getSpliceFromEnv("USERS") {
		if i := strings.Index(e, ":"); i >= 0 {
			users[e[:i]] = e[i+1:]
		}
	}

	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(users))

	authorized.GET("/students", check_student_authorization, getStudents)
	authorized.GET("/students/:id", check_student_authorization, getStudentByID)
	authorized.POST("/students", check_student_authorization, postStudents)
	authorized.DELETE("/students/:id", check_student_authorization, deleteStudentByID)

	authorized.GET("/teachers", check_teacher_authentication, getTeachers)
	authorized.GET("/teachers/:id", check_teacher_authentication, getTeacherByID)
	authorized.POST("/teachers", check_teacher_authentication, postTeachers)
	authorized.DELETE("/teachers/:id", check_teacher_authentication, deleteTeacherByID)

	router.Run("0.0.0.0:8080")

}
