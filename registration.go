package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	_ "github.com/go-sql-driver/mysql"
)

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
			dt := time.Now().Format("2006-01-02 15:04:05")

			db.Create(&User{FirstName: firstName, LastName: lastName, Email: email, Password: password, Country: country, CreatedAt: dt, UpdatedAt: dt})

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