package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
)

func UpdatePassword(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	getHeader(res,req)

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