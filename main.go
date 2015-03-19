package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
		"Host": r.Host,
	}

	t.templ.Execute(w, data)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
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

const socketBufferSize = 1024

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func handleWebsocket(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	log.Println("connect socket", socket)

	for {
		_, m, _ := socket.ReadMessage()
		log.Println(string(m))
	}
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", &templateHandler{filename: "index.html"})
	r.HandleFunc("/rooms", createRoom).Methods("POST")
	http.HandleFunc("/ws", handleWebsocket)

	http.Handle("/", r)

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
