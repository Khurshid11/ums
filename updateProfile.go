package main

import (
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
	_ "github.com/go-sql-driver/mysql"


)

func UpdateProfile(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

/*
	type Data struct {
		UserEmail string
		FirstName string
		LastName  string
		Email     string
		Country   string
		Message   string
		//Name string
	}*/

	//var data Data

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
		dt := time.Now().Format("2006-01-02 15:04:05")

		r := db.Select("id").Where("email = ? AND id != ?", mail,id).First(&user)

		if r.RowsAffected > 0 {
			//render(response,"error.html")
			data.Message =  "Email Already Exists!"

		} else {

			r = db.Model(&user).Where("id = ?", id).Updates(User{FirstName: firstName, LastName: lastName, Email: mail, Country: country, UpdatedAt: dt})


			//destroy session

			session, err := store.Get(req, "store_email")
			session.Values["email"] = ""
			if err != nil {
				log.Print(err)
			}


			var user1 User

			db.Select("email").Where("id = ?", id).First(&user1)

			session.Values["email"] = user1.Email
			session.Save(req, res)

			if r.RowsAffected > 0 {

				data.FirstName=user.FirstName
					data.LastName=user.LastName
					data.Email=user1.Email
					data.Country=user.Country
					data.Message="User Information Updated successfully!"

			} else {

				data.FirstName=user.FirstName
				data.LastName=user.LastName
				data.Email=user1.Email
				data.Country=user.Country
				data.Message="Fail to update Information!"
			}
		}
		getHeader(res,req)
		render(res, "myProfile.html", data)
	}
}