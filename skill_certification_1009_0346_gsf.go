// 代码生成时间: 2025-10-09 03:46:26
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

// Skill represents the structure of a skill.
type Skill struct {
    ID      int    "json:\"id\""
    Name    string "json:\"name\""
    Level   int    "json:\"level\""
    Certified bool   "json:\"certified\""
}

// NewSkill creates a new Skill.
func NewSkill(id int, name string, level int, certified bool) Skill {
    return Skill{
        ID:      id,
        Name:    name,
        Level:   level,
        Certified: certified,
    }
}

// SkillHandler handles HTTP requests related to skills.
type SkillHandler struct {
    skills map[int]Skill
}

// NewSkillHandler creates a new SkillHandler with an empty skills map.
func NewSkillHandler() *SkillHandler {
    return &SkillHandler{
        skills: make(map[int]Skill),
    }
}

// GetAllSkills handles GET requests to fetch all skills.
func (h *SkillHandler) GetAllSkills(w http.ResponseWriter, r *http.Request) {
    // Convert the map to a slice for JSON serialization.
    skills := make([]Skill, 0, len(h.skills))
    for _, skill := range h.skills {
        skills = append(skills, skill)
    }
    
    // Serialize skills to JSON and write the response.
    if err := json.NewEncoder(w).Encode(skills); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// AddSkill handles POST requests to add a new skill.
func (h *SkillHandler) AddSkill(w http.ResponseWriter, r *http.Request) {
    var skill Skill
    if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    
    // Check if the skill already exists.
    if _, exists := h.skills[skill.ID]; exists {
        http.Error(w, "There's already a skill with that ID.", http.StatusConflict)
        return
    }
    
    // Add the new skill to the map.
    h.skills[skill.ID] = skill
    
    // Write the successful response.
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(skill); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// main function to run the HTTP server.
func main() {
    r := mux.NewRouter()
    
    // Create a new skill handler.
    skillHandler := NewSkillHandler()
    
    // Define routes with their methods and corresponding handlers.
    r.HandleFunc("/skills", skillHandler.GetAllSkills).Methods("GET")
    r.HandleFunc("/skills", skillHandler.AddSkill).Methods("POST")
    
    // Start the server.
    fmt.Println("Starting the server on port 8080...
")
    log.Fatal(http.ListenAndServe(":8080", r))
}
