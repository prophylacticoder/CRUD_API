package main

import (
  "net/http"
  "log"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "strconv"
)

type coffee struct {
  Type string `JSON:type`
}

var coffeePool = make(map[int]coffee)
var counter int

func create(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    var t coffee
    err := json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
      http.Error(w, "Invalid JSON.", http.StatusBadRequest)
      return
    }
    coffeePool[counter] = t
    fmt.Fprintf(w, "CID: %v/Type: %v", counter, t.Type)
    counter++
  }
}

func read(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  i, _ := strconv.Atoi(vars["id"])

  if _, ok := coffeePool[i]; !ok {
    http.Error(w, "Invalid ID.", http.StatusBadRequest)
    return
  }

  var s = struct{
    Id int `JSON:id`
    coffee
  }{
    i,
    coffeePool[i],
  }
  b, _ := json.Marshal(s)

  w.Write(b)

}

func readAll(w http.ResponseWriter, r *http.Request){
  b, err := json.Marshal(coffeePool)
  if err != nil {
    http.Error(w, "", http.StatusInternalServerError)
    return
  }
  w.Write(b)
}

func update(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPatch {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])

    if _, ok := coffeePool[id]; !ok {
      http.Error(w, "Invalid ID.", http.StatusBadRequest)
      return
    }

    coffeePool[id] = coffee{vars["type"]}
    fmt.Fprintf(w, "CID: %v | Type: %v", id, coffeePool[id].Type)
  }
}

func del(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodDelete {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])
    if _, ok := coffeePool[id]; !ok {
      http.Error(w, "Invalid ID.", http.StatusBadRequest)
    }
    delete(coffeePool, id)
  }
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/create", create)
  r.HandleFunc("/read/{id:[0-9]+}", read)
  r.HandleFunc("/readall", readAll)
  r.HandleFunc("/update/{id:[0-9]+}/{type}", update)
  r.HandleFunc("/delete/{id:[0-9]+}", del)
  log.Fatal(http.ListenAndServe(":3000", r))
}
