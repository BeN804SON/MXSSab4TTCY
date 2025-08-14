// 代码生成时间: 2025-08-14 15:34:50
package main

import (
    "net/http"
    "strings"
    "html/template"
    "github.com/gorilla/mux"
)

// XSSFilter is a middleware that filters out potential XSS attacks.
func XSSFilter(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        for key, value := range r.URL.Query() {
            r.URL.RawQuery = strings.Map(func(r rune) rune {
                if strings.ContainsRune("<script>"<>&", r) {
                    return -1
                }
                return r
            }, r.URL.RawQuery)
        }
        next.ServeHTTP(w, r)
    })
}

// IndexHandler handles the GET requests to the root path.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    // Use html/template to escape any HTML tags in the input.
    tmpl, err := template.New("index").Parse(`<html><body>Hello, {{.}}!</body></html>`)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, template.HTMLEscapeString(r.URL.Query().Get("name")))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc(`/`, IndexHandler).Methods("GET")
    r.Use(XSSFilter)
    http.Handle(`/`, r)
    
    // Start the server on port 8080.
    http.ListenAndServe(":8080", nil)
}
