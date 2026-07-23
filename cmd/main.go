package main

import (
	"fmt"
	"groupie-tracker/internal"
	"html/template"
	"net/http"
	"os"
)

func Handler(artists []internal.Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, artists)
		if err != nil {
			http.Error(w, "Could not displace page", http.StatusInternalServerError)
		}
	}
}

func main() {
	artists, err := internal.FetchArtists()
	if err != nil {
		fmt.Println("API error:", err)
		return
	}

	if len(artists) == 0 {
		fmt.Println("No artists received")
		return
	}

	//this is how i send data to the html
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler(artists))

	// the StripPrefix is used to prevent the html to repeat a path
	// and ending up somewere it doesnt exists
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running at http://localhost:" + port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Server error:", err)

	}
}
