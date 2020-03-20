package routes

import (
	"github.com/gorilla/mux"
	"log"
	c "miniblog/http/controller"
	"net/http"
)

// initial router and set default path
func InitialRouters() *mux.Router{
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets"))))
	router.HandleFunc("/health_check", c.CheckHealth).Methods("GET")
	return router
}

func Serve(port string, router *mux.Router) {
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("serve failed :%v", err)
	}
}

func SetUserRouters(router *mux.Router, handler *c.UserHandler){
	router.HandleFunc("/login", handler.SignIn).Methods("GET")
	router.HandleFunc("/login", handler.PostSignIn).Methods("POST")
	router.HandleFunc("/logout", handler.SignOut).Methods("GET")
}

func SetArticleRouters(router *mux.Router, handler *c.ArticleHandler){
	router.HandleFunc("/", handler.GetArticles).Methods("GET")
	router.HandleFunc("/article/{id:[0-9]+}", handler.GetArticle).Methods("GET")
	router.HandleFunc("/articles", handler.CreateArticle).Methods("POST")
	router.HandleFunc("/articles/new", handler.NewArticle).Methods("GET")
}
