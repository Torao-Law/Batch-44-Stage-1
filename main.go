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
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": false,
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Blog struct {
	Id          int
	Title       string
	Author      string
	Content     string
	Image       string
	Post_date   time.Time
	Format_date string
}

func main() {
	// declarartion new router
	router := mux.NewRouter()

	connection.DatabaseConnect()
	// create static folder.
	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// create handling URl
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/blog", blog).Methods("GET")
	router.HandleFunc("/add-blog", addBlog).Methods("POST")
	router.HandleFunc("/contact-me", getContact).Methods("GET")
	router.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	router.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	router.HandleFunc("/register", formRegister).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", formLogin).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")

	// running local server
	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", router)
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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["Username"] = session.Values["Name"].(string)
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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["Username"] = session.Values["Name"].(string)
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

		each.Author = "Dandi Gans"
		each.Format_date = each.Post_date.Format("Jan 21, 2000")

		// if session.Values["IsLogin"] != true {
		// 	each.IsLogin = false
		// } else {
		// 	each.IsLogin = session.Values["IsLogin"].(bool)
		// }

		result = append(result, each)
	}

	resp := map[string]interface{}{
		"Data":  Data,
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
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// parsing template html
	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	BlogDetail := Blog{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, image, content, post_date FROM tb_blog WHERE id=$1", id).Scan(&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Image, &BlogDetail.Content, &BlogDetail.Post_date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	BlogDetail.Author = "Dandi Gans"
	BlogDetail.Format_date = BlogDetail.Post_date.Format("19 August 2000")

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

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	var image string
	image = "imageByVariable.png"

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_blog(title, content, image) VALUES ($1, $2, $3)", title, content, image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	fmt.Println(title)
	fmt.Println(content)
	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

// function delete blog
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_blog WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

// function handling register
func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	temp, err := template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	temp.Execute(w, nil)
}

// hashing password register
func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_user(name, email, password) VALUES ($1, $2, $3);", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// function handling register
func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	temp, err := template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	temp.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT Id, email, name, password FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Options.MaxAge = 10800

	session.AddFlash("Login succes", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
