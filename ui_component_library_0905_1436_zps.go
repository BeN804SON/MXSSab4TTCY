// 代码生成时间: 2025-09-05 14:36:37
package main

import (
    "fmt"
    "net/http"
    "html/template"
    "github.com/gorilla/mux"
)

// Component defines the structure for UI components
type Component struct {
    Name    string
    Version string
    Template *template.Template
}

// NewComponent creates a new component with the given name and version
func NewComponent(name, version string, templates ...string) *Component {
    var t *template.Template
    // Load all templates
    for _, tmpl := range templates {
        t = template.Must(t.New(name).ParseFiles(tmpl))
    }
    return &Component{
        Name:    name,
        Version: version,
        Template: t,
    }
}

// Render renders the template with the given data
func (c *Component) Render(w http.ResponseWriter, data interface{}) error {
    // Error handling for template execution
    if err := c.Template.ExecuteTemplate(w, "base", data); err != nil {
        return err
    }
    return nil
}

// StartServer starts the HTTP server with the given component
func StartServer(component *Component) {
    r := mux.NewRouter()
    // Define routes for the component
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        component.Render(w, nil) // Render the component without any data
    })

    // Start the server
    if err := http.ListenAndServe(":8080", r); err != nil {
        fmt.Println("Error starting server: ", err)
    }
}

func main() {
    // Create a new component
    comp := NewComponent("Button", "1.0.0", "./templates/base.html")
    // Start the server with the component
    StartServer(comp)
}

// templates/base.html
// {{ define "base" }}
// <!DOCTYPE html>
// <html>
// <head>
//     <title>{{ .Name }}</title>
// </head>
// <body>
//     <h1>{{ .Name }} - {{ .Version }}</h1>
// </body>
// </html>
// {{ end }}