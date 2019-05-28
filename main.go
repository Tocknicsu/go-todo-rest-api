package main

import (
  "net/http"
  "encoding/json"
  "io/ioutil"

  "github.com/gorilla/mux"
  "github.com/rs/cors"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type Todo struct {
  Title string
  Done bool
  gorm.Model
}

func writeJson (w http.ResponseWriter, data interface{}) {
  res, err := json.Marshal(data)
  if err != nil {
    return
  }
  w.Write(res)
}

func GetTodosHandler (w http.ResponseWriter, r *http.Request) {
  var todos []Todo
  db.Find(&todos)
  writeJson(w, todos)
}

func PostTodosHandler (w http.ResponseWriter, r *http.Request) {
  body, _ := ioutil.ReadAll(r.Body)
  data := map[string]interface{} {}
  json.Unmarshal(body, &data)

  todo := Todo{
    Title: data["Title"].(string),
  }

  db.Create(&todo)
  writeJson(w, todo)
}

func GetTodoHandler (w http.ResponseWriter, r *http.Request) {
  var todo Todo
  vars := mux.Vars(r)

  db.First(&todo, "id = ?", vars["id"])
  writeJson(w, todo)
}

func PutTodoHandler (w http.ResponseWriter, r *http.Request) {
  body, _ := ioutil.ReadAll(r.Body)
  var todo Todo
  vars := mux.Vars(r)
  data := map[string]interface{} {}
  json.Unmarshal(body, &data)

  db.First(&todo, "id = ?", vars["id"])
  db.Model(&todo).Updates(data)

  writeJson(w, todo)
}

func DeleteTodoHandler (w http.ResponseWriter, r *http.Request) {
  var todo Todo
  vars := mux.Vars(r)

  db.First(&todo, "id = ?", vars["id"])
  db.Delete(&todo)
}

func main() {
  db, _ = gorm.Open("sqlite3", "test.db")

  db.AutoMigrate(&Todo{})

  r := mux.NewRouter()
  r.HandleFunc("/todos/", GetTodosHandler).Methods("GET")
  r.HandleFunc("/todos/", PostTodosHandler).Methods("POST")
  r.HandleFunc("/todos/{id}/", GetTodoHandler).Methods("GET")
  r.HandleFunc("/todos/{id}/", PutTodoHandler).Methods("PUT")
  r.HandleFunc("/todos/{id}/", DeleteTodoHandler).Methods("DELETE")

  c := cors.New(cors.Options{
    AllowedOrigins: []string{"http://localhost:3000"},
    AllowCredentials: true,
    Debug: true,
  })

  handler := c.Handler(r)

  http.ListenAndServe(":3030", handler)
}

