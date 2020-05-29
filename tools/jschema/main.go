package main

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func main() {

	//schemaLoader := gojsonschema.NewReferenceLoader("file:///home/fils/src/Projects/samplesEarth/tools/jschema/core.schema.json")
	//schemaLoader := gojsonschema.NewReferenceLoader("https://raw.githubusercontent.com/IGSN/igsn-json/master/schema.igsn.org/json/registration/v0.1/core.schema.json")
	schemaLoader := gojsonschema.NewReferenceLoader("https://raw.githubusercontent.com/IGSN/igsn-json/feature/update-links/schema.igsn.org/json/registration/v0.1/core.schema.json")
	//documentLoader := gojsonschema.NewReferenceLoader("file:///home/fils/src/Projects/samplesEarth/tools/jschema/minimal_registration.json")
	documentLoader := gojsonschema.NewReferenceLoader("file:///home/fils/src/Projects/samplesEarth/tools/jschema/full_registration.json")

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}
