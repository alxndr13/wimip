package main

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type indexResp struct {
	Message string `json:"message"`
}

type ipResp struct {
	IP string `json:"ip"`
}

func main() {
	log.Println("Starting wimip..")
	r := mux.NewRouter()
	r.HandleFunc("/", indexhandler).Methods("GET")
	r.HandleFunc("/ip", wimiphandler).Methods("GET")
	definedRoutes := []string{}
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		log.Println("Available Route:", t)
		definedRoutes = append(definedRoutes, t)
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(http.ListenAndServe("localhost:3000", r))
}

// wandlers
func indexhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indexResp{Message: "Welcome to wimip"})
}

func wimiphandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	x := strings.Split(r.RemoteAddr, ":")
	json.NewEncoder(w).Encode(ipResp{IP: string(x[0])})
}
