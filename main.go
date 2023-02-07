package main

import (
	"Personal-Web/connection"
	"context"
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

// type Home interface{
// 	helloworld(int, model) (string, error)
// }

type Blog struct {
	Id          int
	Title       string
	Author      string
	Content     string
	Image       string
	Post_date   time.Time
	Format_date string
}

// var Blogs = []Blog{
// 	{
// 		Title:     "Pasar Coding di Indonesia Dinilai Masih Menjanjikan",
// 		Post_date: time.Now().String(),
// 		Author:    "Dandi Saputra",
// 		Content:   "Halo ini testing dan ini merupakan dummy data",
// 	},
// }

func main() {
	// declarartion new router
	router := mux.NewRouter()

	connection.DatabaseConnect()
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

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, title, image, content, post_date FROM public.tb_blog;")

	var result []Blog
	for rows.Next() {
		var each = Blog{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Image, &each.Content, &each.Post_date)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "Dandi Saputra"
		each.Format_date = each.Post_date.Format("Jan 21, 2000")

		result = append(result, each)
	}

	resp := map[string]interface{}{
		"Title": Data,
		"Blogs": result,
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

	// id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// parsing template html
	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	BlogDetail := Blog{}

	// for i, data := range Blogs {
	// 	if i == id {
	// 		BlogDetail = Blog{
	// 			Title:     data.Title,
	// 			Post_date: data.Post_date,
	// 			Author:    data.Author,
	// 			Content:   data.Content,
	// 		}
	// 	}
	// }

	resp := map[string]interface{}{
		"Data": Data,
		"Blog": BlogDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")

	// // var newBlog = Blog{
	// // 	Title:     title,
	// // 	Post_date: time.Now().String(),
	// // 	Author:    "Dandi Saputra",
	// // 	Content:   content,
	// // }

	// Blogs = append(Blogs, newBlog)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

// function delete blog
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	// Blogs = append(Blogs[:id], Blogs[id+1:]...)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}
