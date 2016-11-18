package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"html/template"
	"log"

)

var db *gorm.DB
var store = sessions.NewCookieStore([]byte("UsermailId"))

type User struct {
	ID        int `gorm:AUTO_INCREMENT`
	FirstName string
	LastName  string
	Email     string `gorm:"not null;unique"`
	Password  string
	Country   string
	CreatedAt string
	UpdatedAt string
}

//noinspection ALL
func UserMainLogin(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	session, err := store.Get(req, "email")

	if err != nil {
		log.Print(err)
	}

	session.Values["email"] = ""

	render(response, "index.html")

}


func getUser(req *http.Request) string {
	session, err := store.Get(req, "store_email")

	if err != nil {
		log.Print(err)
	}

	email := session.Values["email"].(string)

	return email

}

func getHeader(res http.ResponseWriter,req *http.Request)  {
	email := getUser(req)
	if email == "" {
		http.Redirect(res, req, "/login", http.StatusFound)
	}

	data := struct {
		UserEmail string

	}{
		UserEmail: email,

	}

	render(res,"headerFile.html",data)


}
/*func render(res http.ResponseWriter,fname string){
	tmpl := fmt.Sprintf("template/%s",fname)
	t, err := template.ParseFiles(tmpl)

	if err != nil{
		log.Print("Template Parsing error ",err)
	}

	err = t.Execute(res,nil)

	if err != nil{
		log.Print("Parsing error ",err)
	}

}*/

func render(param ...interface{}) {

	fname := param[1].(string)
	res := param[0].(http.ResponseWriter)

	if len(param) == 2 {

		tmpl := fmt.Sprintf("template/%s", fname)
		t, err := template.ParseFiles(tmpl)

		if err != nil {
			log.Print("Template Parsing error ", err)
		}

		err = t.Execute(res, nil)

		if err != nil {
			log.Print("Parsing error ", err)
		}
	}

	if len(param) == 3 {

		//	data := param[2].(struct{})
		tmpl := fmt.Sprintf("template/%s", fname)
		t, err := template.ParseFiles(tmpl)

		if err != nil {
			log.Print("Template Parsing error ", err)
		}

		err = t.Execute(res, param[2])

		if err != nil {
			log.Print("Parsing error ", err)
		}

	}

}

func main() {

	router := httprouter.New()
	router.GET("/", UserMainLogin)

	var err error
	db, err = gorm.Open("mysql", "root:password@/userInfo?charset=UTF8")

	if err != nil {
		fmt.Println("Fail to connect to database-->", err)
	}

	if db.HasTable(&User{}) == false {

		db.AutoMigrate(&User{})
	}

	router.GET("/login", Login)
	router.POST("/login", Login)

	router.GET("/registerUser", NewUser)
	router.POST("/registerUser", NewUser)

	router.GET("/profile", ViewProfile)
	router.POST("/profile", UpdateProfile)

	router.GET("/updatePswd", UpdatePassword)
	router.POST("/updatePswd", UpdatePassword)

	router.GET("/userList",ViewUsers)
	router.ServeFiles("/assets/*filepath", http.Dir("template/assets"))

	log.Fatal(http.ListenAndServe(":8000", router))

	defer db.Close()
}
