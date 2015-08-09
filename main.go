package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type authHandler struct {
	next http.Handler
}

func init() {
	// set up gomniauth
	gomniauth.SetSecurityKey("some long key")
	gomniauth.WithProviders(
		facebook.New("key", "secret",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:8080/auth/callback/github"),
		google.New("636296155193-3u8lt2kmcr42mt49qmcvoq726dv9q9kj.apps.googleusercontent.com", "pq9s2KpbPyt6g-0kDM_ef-7F",
			"http://localhost:8080/auth/callback/google"),
	)
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

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// error
		panic(err.Error())
	} else {
		// success. call the next handler
		h.next.ServeHTTP(w, r)
	}
}
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func main() {
	http.Handle("/", MustAuth(&templateHandler{filename: "index.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/rooms", handleRoom)
	http.HandleFunc("/ws", handleWebsocket)

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
