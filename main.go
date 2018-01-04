package main

import (
	"fmt"
	"github.com/joshbetz/config"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var user string
	var pass string
	// config
	c := config.New("config.json")
	c.Get("user", &user)
	c.Get("pass", &pass)

	// form file takes the post input id file
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Println(w, err)
		return
	}

	defer file.Close()
	if user != r.FormValue("user") || pass != r.FormValue("pass") {
		fmt.Fprint(w, "Auth failed.")
		return
	}
	out, err := os.OpenFile("/var/www/html/img/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprint(w, "Could not write file")
		return
	}

	defer out.Close()

	// write content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "https://enju.us/img/")
	fmt.Fprintf(w, header.Filename)
}

func main() {
	fmt.Print("running")
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
