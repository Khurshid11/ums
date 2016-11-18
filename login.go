package main

import (
	"strings"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	_ "github.com/go-sql-driver/mysql"
)

func Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	//fmt.Print(req.Method)
	//----------If request is get do the following-----------

	if req.Method != "POST" {

		session, err := store.Get(req, "store_email")

		if err != nil {
			log.Print(err)
		}

		session.Values["email"] = ""
		render(res, "index.html")
		return

	}

	//------------end for get Request-------------


	//------------If Request is Post--------Execute----------


	// struct for templating which is define in index.html
	type Message struct {
		EmailError string
		PasswordError string
		EmailPasswordError string
	}

	loginData :=&Message{"","",""}



	if req.Method == "POST" {
		email := req.FormValue("email")
		password := req.FormValue("password")

		count := strings.Count(email, "@")
		if count > 1{
			loginData.EmailError="Invalid EmailID"

		}

		if password=="" {
			loginData.PasswordError="You need to provide password"
		}

		if count > 1 && password==""{
			loginData.PasswordError=""
			loginData.EmailError=""
			loginData.EmailPasswordError="Email or Password Incorrect!"
		}

		if (loginData.EmailPasswordError=="") || (loginData.EmailError=="") || (loginData.PasswordError=="") {

			var user User

			r := db.Where(&User{Email: email, Password: password}).First(&user)

			if r.RowsAffected == 1 {

				session, err := store.Get(req, "store_email")

				if err != nil {
					log.Print(err)
				}

				session.Values["email"] = email
				session.Save(req, res)

				http.Redirect(res, req, "/profile", 302)
				return
			} else {

				loginData.PasswordError=""
				loginData.EmailError=""
				loginData.EmailPasswordError="Email or Password Incorrect!"
				render(res, "index.html", loginData)
			}
		}else {
			render(res,"index.html",loginData)
		}

	}
}