package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": true,
}

type Blog struct {
	Title     string
	Post_date string
	Author    string
	Content   string
}

var Blogs = []Blog{
	{
		Title:     "Pasar Coding di Indonesia Dinilai Masih Menjanjikan 0",
		Post_date: "12 Jul 2021 22:30 WIB",
		Author:    "Dandi Saputra",
		Content:   "Halo ini testing dan ini merupakan dummy data",
	},
	{
		Title:     "Pasar Coding di Indonesia Dinilai Masih Menjanjikan 1",
		Post_date: "12 Jul 2021 22:30 WIB",
		Author:    "Dandi Saputra",
		Content:   "Halo ini testing dan ini merupakan dummy data",
	},
	{
		Title:     "Pasar Coding di Indonesia Dinilai Masih Menjanjikan 2",
		Post_date: "12 Jul 2021 22:30 WIB",
		Author:    "Dandi Saputra",
		Content:   "Halo ini testing dan ini merupakan dummy data",
	},
	{
		Title:     "Pasar Coding di Indonesia Dinilai Masih Menjanjikan 3",
		Post_date: "12 Jul 2021 22:30 WIB",
		Author:    "Dandi Saputra",
		Content:   "Halo ini testing dan ini merupakan dummy data",
	},
}

func main() {
	// declarartion new router
	router := mux.NewRouter()

	// create static folder.
	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// create handling URl
	router.HandleFunc("/hello", helloWorld).Methods("GET")
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/blog", blog).Methods("GET")
	router.HandleFunc("/add-blog", addBlog).Methods("POST")
	router.HandleFunc("/contact-me", getContact).Methods("GET")
	router.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	router.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")

	// running local server
	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", router)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

// function handling index.html
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("views/index.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// function handling blog.html
func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("views/blog.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Title": Data,
		"Blogs": Blogs,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// function handling contact me
func getContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("views/contactme.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// function handling blog-detail.html with query string
func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// parsing template html
	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data": Data,
		"Id":   id,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	var newBlog = Blog{
		Title:     title,
		Post_date: time.Now().String(),
		Author:    "Dandi Saputra",
		Content:   content,
	}

	Blogs = append(Blogs, newBlog)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

// function delete blog
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Blogs = append(Blogs[:id], Blogs[id+1:]...)

	fmt.Println(id)
	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}
