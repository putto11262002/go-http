package main

import logserver "github.com/putto112620002/go-http/log_server"

func main(){
	logServer := logserver.NewLogServer("127.0.0.1:8081")
	logServer.Run()
}