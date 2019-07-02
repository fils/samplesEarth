package main

import (
	"log"
	"net/http"

	"../../internal/kv"
	"../../internal/sample"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

// MyServer struct for mux router
type MyServer struct {
	r *mux.Router
}

func main() {

	// move to init?
	db, err := bolt.Open("data/my.db", 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	kv.DBCon = db

	do := mux.NewRouter()
	do.HandleFunc("/id/{ID}", sample.DigitalObject)
	http.Handle("/id/", do)

	// Common assets like; css, js, images, etc...
	rcommon := mux.NewRouter()
	rcommon.PathPrefix("/common/").Handler(http.StripPrefix("/common/", http.FileServer(http.Dir("./web/static"))))
	http.Handle("/common/", &MyServer{rcommon})

	htmlRouter := mux.NewRouter()
	htmlRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	http.Handle("/", &MyServer{htmlRouter})

	log.Printf("Listening on 9990. Go to http://127.0.0.1:9990/")
	//err = http.ListenAndServeTLS(":443", "./secret/server.crt", "./secret/server.key", nil)
	err = http.ListenAndServe(":9990", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Let the Gorilla work
	s.r.ServeHTTP(rw, req)
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}
