package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	ServiceVersion string
	StartTime      time.Time
)

func init() {
	StartTime = time.Now()
	ServiceVersion = "1.0.231.fd80bea-" + time.Now().Format("20060102.150405")
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
		if !c.IsActive {
			continue
		}
		d := time.Now().Sub(c.ConnectedAt)
		m := int(d.Minutes())
		s := int((d - time.Minute*time.Duration(m)).Seconds())
		users = append(users, map[string]interface{}{
			"Name":            c.Name,
			"Addr":            c.Messager.RemoteAddr(),
			"SessionDuration": fmt.Sprintf("%dm%ds", m, s),
		})
	}
	info := map[string]interface{}{
		"CurrentUserCount":            len(users),
		"CurrentlyAuthenticatedUsers": users,
	}

	data, _ := json.Marshal(info)
	w.Write(data)
}
