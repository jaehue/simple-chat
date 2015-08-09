package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		var views []string

		files, _ := ioutil.ReadDir(filepath.Join("templates", "partials"))
		for _, f := range files {
			views = append(views, filepath.Join("templates", "partials", f.Name()))
		}
		views = append([]string{filepath.Join("templates", t.filename)}, views...)

		t.templ = template.Must(template.ParseFiles(views...))
	})

	data := map[string]interface{}{
		"Host":  r.Host,
		"Rooms": rooms,
	}

	t.templ.Execute(w, data)
}

func handleRoom(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		b, err := json.Marshal(rooms)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
		return
	}

	if req.Method != "POST" {
		return
	}

	decoder := json.NewDecoder(req.Body)
	var room Room
	decoder.Decode(&room)
	if len(room.Name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	newRoom := NewRoom(room.Name)

	b, err := json.Marshal(newRoom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func main() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/rooms", handleRoom)
	http.HandleFunc("/ws", handleWebsocket)

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
