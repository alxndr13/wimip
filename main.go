package main

import (
	"encoding/json"
	"fmt"
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

var routesJSON = map[string]interface{}{}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.Println("Starting wimip..")
	r := mux.NewRouter()
	r.HandleFunc("/", indexhandler).Methods("GET")
	r.HandleFunc("/ip", wimiphandler).Methods("GET")
	r.Use(loggingMiddleware)
	err := generateListOfRoutes(r)
	if err != nil {
		log.Warnln("Could not generate list of routes")
	}

	log.Println("Starting on localhost:3000")
	log.Fatalln(http.ListenAndServe("localhost:3000", r))
}

func generateListOfRoutes(r *mux.Router) error {
	definedRoutes := []string{}
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		definedRoutes = append(definedRoutes, t)
		return nil
	})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	routesJSON["routes"] = definedRoutes
	return nil
}

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s - %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// handlers
func indexhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// return pretty JSON, if possible
	pretty, err := json.MarshalIndent(routesJSON, "", "   ")
	if err != nil {
		log.Warnln("Could not prettify json. Will return normal json")
		json.NewEncoder(w).Encode(routesJSON)
	}
	fmt.Fprintf(w, string(pretty))
}

func wimiphandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	x := strings.Split(r.RemoteAddr, ":")
	json.NewEncoder(w).Encode(ipResp{IP: string(x[0])})
}
