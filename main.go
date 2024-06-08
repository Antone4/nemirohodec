package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var router = mux.NewRouter()
var store = sessions.NewCookieStore([]byte("0zxz0jRyzxk7D2CzCJ44Ix7cvOMu0P1m5VkYvf3OSyALs46KyjJguC13cE5fhXc4bjfWtY1w9IXmJot0DM37zof5mQzNy0um"))

var thisPage string = "localhost:8888"

func main() {
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		icon, err := os.ReadFile("static/Images/favicon.ico")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(icon)
	})
	//router.HandleFunc("/static/", serveStaticFiles)
	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.HandlerFunc(staticHandler)))
	http.ListenAndServe(":8080", router)
}

// func staticHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Host != "localhost:8080" {
// 		Printfln("NOT HOst")
// 		http.NotFound(w, r)
// 		return
// 	}

// 	filepath := "static/" + r.URL.Path

// 	if _, err := os.Stat(filepath); os.IsNotExist(err) {

// 		Printfln("File not exist: %v", filepath)
// 		http.NotFound(w, r)
// 		return

// 	}
// 	Printfln("ServeFile: %v", filepath)
// 	http.ServeFile(w, r, filepath)
// }

func staticHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	if strings.Contains(referer, ("http://" + thisPage + "/")) {
		// file, err := os.Open(r.URL.Path)
		// defer file.Close()
		// if err != nil {
		// 	Printfln("Error file not found: %v", err.Error())
		// 	http.Error(w, "File not found", http.StatusNotFound)
		// 	return
		// }

		// fileInfo, err := file.Stat()
		// if err != nil {
		// 	http.Error(w, "Unable to read file information", http.StatusInternalServerError)
		// 	return
		// }
		http.ServeFile(w, r, "static/"+r.URL.Path)
		//http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	} else {
		http.Redirect(w, r, "/PageNotFound", http.StatusMovedPermanently)
	}
}
