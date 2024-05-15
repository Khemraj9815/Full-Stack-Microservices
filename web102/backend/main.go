package main


import (
    "fmt"
    "log"
    "net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    fmt.Fprintf(w, "POST request successful")

	http.ServeFile(w, r, "templates/login.html")
   
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/hello" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }


    fmt.Fprintf(w, "Hello!")
}

// ? Handler for the home route
func authicationHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "authenticated")
}



func main() {
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    http.HandleFunc("/form", formHandler)
    http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", loginpageHandler)
    http.HandleFunc("/auth", authicationHandler)

    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func loginpageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./login.html")
}