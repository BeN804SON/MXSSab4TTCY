// 代码生成时间: 2025-08-31 12:10:49
package main

import (
    "net/http"
    "strings"
    "github.com/gorilla/mux"
)

// ThemeSwitcherHandler is a handler function to switch themes.
func ThemeSwitcherHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    theme := vars["theme"]

    // Check if the theme is valid. This is a simple example and only two themes are supported.
    if theme != "light" && theme != "dark" {
        http.Error(w, "Invalid theme", http.StatusBadRequest)
        return
    }

    // Store the theme in a cookie.
    http.SetCookie(w, &http.Cookie{Name: "theme", Value: theme, Path: "/"})
    fmt.Fprintf(w, "Theme switched to %s", theme)
}

// main function to start the server.
func main() {
    router := mux.NewRouter()

    // Define routes.
    router.HandleFunc("/switch-theme/{theme}", ThemeSwitcherHandler).Methods("GET")

    // Start the server.
    http.ListenAndServe(":8080", router)
}
