package main

import ("net/http"; "io")
import "fmt"
func Hello(res http.ResponseWriter, req *http.Request){
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		`<h1>Hello World</h1>`,
	)
}

func main(){
	http.HandleFunc("/metrics", Hello)
	fmt.Printf("Server is running at 0.0.0.0:9000")
	http.ListenAndServe("127.0.0.1:9000", nil)
}