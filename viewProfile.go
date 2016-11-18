package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	FirstName string
	LastName  string
	Email     string
	Country   string
	Message   string
}

var data =new (Data)

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

	data.FirstName=user.FirstName
	data.LastName=user.LastName
	data.Email=email
		data.Country  =user.Country
		data.Message=   ""

	getHeader(res,req)
	render(res, "myProfile.html", data)
}