package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	//"text/template"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"html/template"
	"log"
	"strings"

	"time"
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
	CreatedAt time.Time
	UpdatedAt time.Time
}

//noinspection ALL
func Index(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	session, err := store.Get(req, "email")

	if err != nil {
		log.Print(err)
	}

	session.Values["email"] = ""

	render(response, "index.html")

}

func Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	//fmt.Print(req.Method)
	if req.Method != "POST" {

		session, err := store.Get(req, "email")

		if err != nil {
			log.Print(err)
		}

		session.Values["email"] = ""
		render(res, "index.html")
		return

	}

	if req.Method == "POST" {
		email := req.FormValue("email")
		password := req.FormValue("password")

		count := strings.Count(email, "@")

		if count > 1 {

			render(res, "index.html")
			return
		}

		var user User

		r := db.Where(&User{Email: email, Password: password}).First(&user)

		if r.RowsAffected == 1 {

			session, err := store.Get(req, "email")

			if err != nil {
				log.Print(err)
			}

			session.Values["email"] = email
			session.Save(req, res)

			http.Redirect(res, req, "/profile", 302)
			return
		} else {

			render(res, "index.html")
		}

	}
}

func ViewProfile(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	/*t, err := template.ParseFiles("template/myProfile.html")
	if err != nil{
		log.Fatal(err)
	}*/
	//Get userEmail from session
	email := getUser(req)

	//if session is not set then logout
	if email == "" {
		http.Redirect(res, req, "/login", http.StatusFound)
	}
	//fmt.Print(email)

	//access values from database

	var user User

	db.Select("first_name,last_name,country").Where("email=?", email).First(&user)

	data := struct {
		UserEmail string
		FirstName string
		LastName  string
		Email     string
		Country   string
		Message   string
	}{
		UserEmail: email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     email,
		Country:   user.Country,
		Message:   "",
	}

	render(res, "myProfile.html", data)
}

func UpdateProfile(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	email := getUser(req)

	if email == "" {
		http.Redirect(res, req, "/login", http.StatusFound)
	}

	if req.Method == "POST" {
		var user User

		db.Select("id").Where("email=?", email).First(&user)
		id := user.ID

		firstName := req.FormValue("firstName")
		lastName := req.FormValue("lastName")
		mail := req.FormValue("mail")
		country := req.FormValue("country")

		r := db.Model(&user).Where("id = ?", id).Updates(User{FirstName: firstName, LastName: lastName, Email: mail, Country: country, UpdatedAt: time.Now()})
		var data interface{}

		//destroy session

		session, err := store.Get(req, "email")

		if err != nil {
			log.Print(err)
		}

		var user1 User

		db.Select("email").Where("id = ?", id).First(&user1)

		session.Values["email"] = user1.Email
		session.Save(req, res)

		if r.RowsAffected > 0 {

			data = struct {
				UserEmail string
				FirstName string
				LastName  string
				Email     string
				Country   string
				Message   string
				//Name string

			}{
				UserEmail: user1.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user1.Email,
				Country:   user.Country,
				Message:   "User Information Updated successfully!",
			}
		} else {

			data = struct {
				Message string
			}{
				Message: "Information fail to update!",
			}
		}

		render(res, "myProfile.html", data)
	}
}

func NewUser(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == "GET" {
		render(response, "registration.html")
	} else if req.Method == "POST" {

		firstName := req.FormValue("firstName")
		lastName := req.FormValue("lastName")
		email := req.FormValue("email")
		password := req.FormValue("password")
		country := req.FormValue("country")

		//fmt.Fprint(response,firstName," ",lastName," ",email," ",password," ",country)

		if firstName != "" && password != "" {

			var user User
			r := db.Where(&User{Email: email}).First(&user)

			if r.RowsAffected > 0 {
				//render(response,"error.html")
				render(response, "registration.html")
				return
			}

			db.Create(&User{FirstName: firstName, LastName: lastName, Email: email, Password: password, Country: country, CreatedAt: time.Now(), UpdatedAt: time.Now()})

			if db.NewRecord(&User{}) == true {
				const msg = `
							<!DOCTYPE html><div align="center">
							<h2>Thank You For Registration</h2><br>
								To Login Click <a href="/login">here</a>
							</div>`

				fmt.Fprint(response, msg)
				//http.Redirect(response,req,"/login",302)

			}

		} else {
			render(response, "registration.html")
		}
	}

}

func UpdatePassword(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	email := getUser(req)

	data := struct {
		UserEmail string
		Message   string
	}{
		UserEmail: email,
		Message:   "",
	}

	if req.Method == "GET" {
		t, err := template.ParseFiles("template/updatePassword.html")
		if err != nil {
			log.Fatal(err)
		}
		//Get userEmail from session

		err = t.Execute(res, data)
		if err != nil {
			log.Fatal(err)
		}
	} else if req.Method == "POST" {

		changePassword(res, req, data, email)
	}

}

func changePassword(res http.ResponseWriter, req *http.Request, data struct {
	UserEmail string
	Message   string
}, email string) {
	oldPasword := req.FormValue("oldPswd")
	newPswd := req.FormValue("newPswd")
	cnewPswd := req.FormValue("cnewPswd")

	if newPswd != cnewPswd {

		//http.Redirect(res,req,"/updatePswd",http.StatusFound)
		//return

		data.Message = "Password not match!!"
		render(res, "updatePassword.html", data)

	} else {

		var user User

		r := db.Model(&user).Where("email = ? AND password = ?", email, oldPasword).Update("password", cnewPswd)

		if r.RowsAffected > 0 {
			data.Message = "Your password is successfully updated!"
			render(res, "updatePassword.html", data)
		} else {
			data.Message = "Your Password was incorrect!"
			render(res, "updatePassword.html", data)
		}
	}

}

func getUser(req *http.Request) string {
	session, err := store.Get(req, "email")

	if err != nil {
		log.Print(err)
	}

	email := session.Values["email"].(string)

	return email

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
	router.GET("/", Index)

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

	router.ServeFiles("/assets/*filepath", http.Dir("template/assets"))

	log.Fatal(http.ListenAndServe(":8000", router))

	defer db.Close()
}
