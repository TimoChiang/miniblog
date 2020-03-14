package main

import (
	"miniblog/routes"
)


func main() {
	router := routes.SetRouters()
	routes.Serve("8888", router)
}
