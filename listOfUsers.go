package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

func ViewUsers(res http.ResponseWriter,req *http.Request,_ httprouter.Params)  {


	getHeader(res,req)

	type UserList struct {
		SlNo	int
		Id 	int
		FirstName string
		LastName  string
		UserEmail     string
		Country   string
		RegisteredDate string
		//Name string

	}

	var user []User
	r := db.Select("id,first_name,last_name,email,country,created_at").Find(&user)

	nusers := r.RowsAffected



	users := make([]UserList,nusers)




	cnt :=1

	for i, list := range user{
		users[i].Id = list.ID
		users[i].SlNo = cnt
		users[i].FirstName=list.FirstName
		users[i].LastName=list.LastName
		users[i].UserEmail=list.Email
		users[i].Country=list.Country
		users[i].RegisteredDate=list.CreatedAt
		cnt++

		//fmt.Fprint(res,i,cnt)
	}

	//fmt.Fprint(res,users)
	render(res,"userList.html",users)


}