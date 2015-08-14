package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	ServiceVersion string
	StartTime      time.Time
)

func init() {
	StartTime = time.Now()
}

func info(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	info := map[string]interface{}{
		"Version":         ServiceVersion,
		"StartTimeSecs":   StartTime.UTC().Unix(),
		"CurrentTimeSecs": now.UTC().Unix(),
		"Uptime":          now.Sub(StartTime),
	}
	data, _ := json.Marshal(info)
	w.Write(data)
}

func connections(w http.ResponseWriter, req *http.Request) {
	users := []interface{}{}
	for _, c := range clients {
		users = append(users, map[string]interface{}{
			"Name": c.Name,
		})
	}
	info := map[string]interface{}{
		"CurrentUserCount":            len(clients),
		"CurrentlyAuthenticatedUsers": users,
	}

	data, _ := json.Marshal(info)
	w.Write(data)
}
