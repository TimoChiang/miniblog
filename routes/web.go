package routes

import (
	"github.com/gorilla/mux"
	c "miniblog/http/controller"
	"net/http"
)


func SetRouters() *mux.Router{
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets"))))
	router.HandleFunc("/", c.GetArticles).Methods("GET")
	router.HandleFunc("/article/{id:[0-9]+}", c.GetArticle).Methods("GET")
	router.HandleFunc("/articles", c.CreateArticle).Methods("POST")
	router.HandleFunc("/articles/new", c.NewArticle).Methods("GET")
	router.HandleFunc("/login", c.SignIn).Methods("GET")
	router.HandleFunc("/login", c.PostSignIn).Methods("POST")
	router.HandleFunc("/logout", c.SignOut).Methods("GET")
	return router
}

func Serve(port string, router *mux.Router) {
	http.ListenAndServe(":" + port, router)
}