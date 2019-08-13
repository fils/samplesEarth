package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
	"github.com/gosuri/uiprogress"
	"github.com/icrowley/fake"
	"github.com/minio/minio-go"
	"github.com/rs/xid"
	// gominio "github.com/minio/minio-go"
)

type ranData struct {
	PID       string
	Lat       int
	Long      int
	Paragraph string
	Date      string
}

func main() {
	fmt.Println("IGSN sample generator")
	rand.Seed(time.Now().UnixNano())

	count := 3000000
	bar := uiprogress.AddBar(count).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Task (%d/%d)", b.Current(), count)
	})
	uiprogress.Start() // start rendering

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	for n := 0; n <= count; n++ {
		guid := xid.New()
		long := fake.LongitudeDegrees()
		lat := fake.LatitudeDegrees()
		// para := babbler.Babble()
		//para := fake.Paragraph()
		para := randomdata.Paragraph()
		date := randate().Format(time.RFC3339)

		rd := ranData{guid.String(), lat, long, para, date}
		jld := newRandomSample(rd) //  send in some XID value for ID..   send in some random lat longs
		// fmt.Printf("\n ------- \n %s  \n ----------\n", jld)

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			err := b.Put([]byte(guid.String()), []byte(jld))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return err
		})
		bar.Incr()
	}

	time.Sleep(time.Second)
	uiprogress.Stop()
}

func newRandomSample(rd ranData) string {
	var buf bytes.Buffer

	t := template.Must(template.New("template").Parse(s2))
	err := t.Execute(&buf, rd)
	if err != nil {
		fmt.Println(err)
	}

	return buf.String()
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)

	}
	return res
}

func randate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2010, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func minioConnection(minioVal, portVal, accessVal, secretVal string) *minio.Client {
	endpoint := fmt.Sprintf("%s:%s", minioVal, portVal)
	accessKeyID := accessVal
	secretAccessKey := secretVal
	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Println(err)
	}
	return minioClient
}

const s2 = `
{
	"@context": {
		"@vocab": "http://schema.org/",
		"datacite": "http://purl.org/spar/datacite/",
		"earthcollab": "https://library.ucar.edu/earthcollab/schema#",
		"geolink": "http://schema.geolink.org/1.0/base/main#",
		"vivo": "http://vivoweb.org/ontology/core#",
		"dbpedia": "http://dbpedia.org/resource/",
		"geo-upper": "http://www.geoscienceontology.org/geo-upper#"

	},
	"@id": "http://sample.igsn.org/soilarchive/{{.PID}}",
	"@type": [
		"http://www.w3.org/2002/07/owl#Thing",
		"http://www.w3.org/ns/sosa/Sample"

	],
	"spatialCoverage": {
		"@type": "Place",
		"geo": {
			"@type": "GeoCoordinates",
			"latitude":  {{.Lat}},
			"longitude":  {{.Long}}

		}

	},
	"http://purl.org/dc/terms/title": "ANZ soil sample {{.PID}}",
	"description": "{{.Paragraph}}",

	"additionalType": [{
			"@id": "http://pid.geoscience.gov.au/def/voc/igsn-codelists/PhysicalSample"

		},
		{
			"@id": "http://pid.geoscience.gov.au/def/voc/igsn-codelists/soil"

		}
	],
	"creator": {
		"@id": "http://sample.igsn.org/soilarchive/CDS-NSW"

	},
	"dateCreated": {
		"@type": "http://www.w3.org/2001/XMLSchema#date",
		"@value": "{{.Date}}"

	},
	"title": "ANZ soil sample",
	"url": {
		"@id": "http://samples.earth/id/{{.PID}}"

	},
	"http://www.w3.org/2000/01/rdf-schema#label": "ANZ soil sample",
	"http://www.w3.org/ns/dcat#landingPage": {
		"@id": "http://samples.earth/doc/{{.PID}}"

	},
	"http://www.w3.org/ns/sosa/isResultOf": {
		"@id": "_:b0"

	},
	"http://www.w3.org/ns/sosa/isSampleOf": [{
			"@id": "http://www.anzsoil.org/data/csiro-natsoil/anzsoilml201/soil/soil_199.CAN.C410"

		},
		{
			"@id": "http://www.anzsoil.org/data/csiro-natsoil/anzsoilml201/soilhorizon/soil_horizon_199.CAN.C410.1.2"

		}
	]

}
`

const s = `{
	"@graph": [
	  {
		"@id": "_:b0",
		"@type": "http://www.w3.org/ns/sosa/Sampling",
		"http://www.w3.org/ns/sosa/usedProcedure": {
		  "@id": "http://www.anzsoil.org/def/au/soil/observation-method/soil-pit"
		}
	  },
	  {
		"@id": "_:b1",
		"@type": "http://purl.org/dc/terms/Location",
		"http://www.w3.org/ns/dcat#centroid": {
		  "@type": "http://www.opengis.net/ont/geosparql#asWKT",
		  "@value": "POINT({{.Long}} {{.Lat}})"
		}
	  },
	  {
		"@id": "http://sample.igsn.org/soilarchive/{{.PID}}",
		"@type": [
		  "http://www.w3.org/2002/07/owl#Thing",
		  "http://www.w3.org/ns/sosa/Sample"
		],
		"http://purl.org/dc/terms/created": {
		  "@type": "http://www.w3.org/2001/XMLSchema#date",
		  "@value": "1959-10-08"
		},
		"http://purl.org/dc/terms/creator": {
		  "@id": "http://sample.igsn.org/soilarchive/CDS-NSW"
		},
		"http://purl.org/dc/terms/issued": {
		  "@type": "http://www.w3.org/2001/XMLSchema#date",
		  "@value": "2017-01-03"
		},
		"http://purl.org/dc/terms/spatial": {
		  "@id": "_:b1"
		},
		"http://purl.org/dc/terms/title": "ANZ soil sample {{.PID}}",
		"http://schema.org/description": "{{.Paragraph}}",
		"http://purl.org/dc/terms/type": [
		  {
			"@id": "http://pid.geoscience.gov.au/def/voc/igsn-codelists/PhysicalSample"
		  },
		  {
			"@id": "http://pid.geoscience.gov.au/def/voc/igsn-codelists/soil"
		  }
		],
		"http://schema.org/additionalType": [
		  {
			"@id": "http://pid.geoscience.gov.au/def/voc/igsn-codelists/PhysicalSample"
		  },
		  {
			"@id": "http://pid.geoscience.gov.au/def/voc/igsn-codelists/soil"
		  }
		],
		"http://schema.org/creator": {
		  "@id": "http://sample.igsn.org/soilarchive/CDS-NSW"
		},
		"http://schema.org/dateCreated": {
		  "@type": "http://www.w3.org/2001/XMLSchema#date",
		  "@value": "{{.Date}}"
		},
		"http://schema.org/identifier": "soil_specimen_{{.}}",
		"http://schema.org/title": "ANZ soil sample",
		"http://schema.org/url": {
		  "@id": "http://samples.earth/id/{{.PID}}"
		},
		"http://www.w3.org/2000/01/rdf-schema#label": "ANZ soil sample",
		"http://www.w3.org/ns/dcat#landingPage": {
		  "@id": "http://samples.earth/doc/{{.PID}}"
		},
		"http://www.w3.org/ns/sosa/isResultOf": {
		  "@id": "_:b0"
		},
		"http://www.w3.org/ns/sosa/isSampleOf": [
		  {
			"@id": "http://www.anzsoil.org/data/csiro-natsoil/anzsoilml201/soil/soil_199.CAN.C410"
		  },
		  {
			"@id": "http://www.anzsoil.org/data/csiro-natsoil/anzsoilml201/soilhorizon/soil_horizon_199.CAN.C410.1.2"
		  }
		]
	  }
	]
  }
`
