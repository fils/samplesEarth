package sample

import (
	"log"
	"net/http"

	"../kv"

	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

type Data struct {
	Name   string
	Date   string
	Desc   string
	JSONLD string
}

func DigitalObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oid := vars["ID"]

	db := kv.DBCon

	var v []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v = b.Get([]byte(oid))
		return nil
	})

	results := Data{Name: JSONLeaf(v, "@id"), Date: JSONLeaf(v, "dateCreated.@value"), Desc: JSONLeaf(v, "description"), JSONLD: string(v)}

	ht, err := template.New("some template").ParseFiles("web/templates/object.html") // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

func JSONLeaf(jld []byte, path string) string {
	value := gjson.Get(string(jld), path)

	return value.String()
}
