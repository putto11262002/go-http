package main

import bookapi "github.com/putto112620002/go-http/book_api"

func main(){
	api := bookapi.NewBookApi()
	api.Run()
}