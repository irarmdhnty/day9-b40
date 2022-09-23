package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"my-project/connection"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	connection.ConnectDatabase()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/project", project).Methods("GET")
	r.HandleFunc("/add-project", addProject).Methods("POST")
	r.HandleFunc("/detail/{index}", detail).Methods("GET")
	r.HandleFunc("/delete/{index}", delete).Methods("GET")

	fmt.Println("server on in port 8000")
	http.ListenAndServe("localhost:8000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description FROM blog")
	// fmt.Println(data)

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Name, &each.Desc)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	fmt.Println(result)

	card := map[string]interface{}{
		"Add": result,
	}

	tmpl.Execute(w, card)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/addProject.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

type Project struct {
	ID           int
	Name         string
	Start_date   string
	End_date     string
	Duration     string
	Desc         string
}

var data = []Project{}


func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var start_date = r.PostForm.Get("startDate")
	var end_date = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("desc")

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, start_date)
	dateEnd, _ := time.Parse(layout, end_date)

	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := daysInHours / 30
	yearInMonth := monthInDay / 12 

	var duration string
	var month, _ float64 = math.Modf(monthInDay)
	var year, _ float64 = math.Modf(yearInMonth)

	if year > 0 {
		duration = strconv.FormatFloat(year, 'f', 0, 64) + " Years"
		// fmt.Println(year, " Years")
	} else if month > 0 {
		duration = strconv.FormatFloat(month, 'f', 0, 64) + " Months"
		// fmt.Println(month, " Months")
	} else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Days"
		// fmt.Println(daysInHours, " Days")
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
		// fmt.Println(hours, " Hours")
	} else {
		duration = "0 Days"
		// fmt.Println("0 Days")
	}


	var newData = Project{
		Name:         name,
		Start_date:   start_date,
		End_date:     end_date,
		Duration:     duration,
		Desc:         desc,
	
	}

	data = append(data, newData)
	// fmt.Println(data)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Detail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)

	for i, data := range data {
		if index == i {
			Detail = Project{
				Name:       data.Name,
				Start_date: data.Start_date,
				End_date:   data.End_date,
				Desc:       data.Desc,
			}
		}
	}

	data := map[string]interface{}{
		"Details": Detail,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
}

func delete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	data = append(data[:index], data[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}


