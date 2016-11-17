package main

import (
	"html/template"
	"log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)



//noinspection ALL
func Index(w http.ResponseWriter,r *http.Request,_ httprouter.Params)  {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.ParseFiles("template/test.html")
	check(err)

	data := struct {
		Title string
		Name string

	}{
		Title : "Test Title",
		Name: "Alam",

	}

	err = t.Execute(w, data)
	check(err)


}


func main() {



	router := httprouter.New()
	router.GET("/",Index)

	log.Fatal(http.ListenAndServe(":8000",router))



}

