package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/alecthomas/template"
	"github.com/gosuri/uiprogress"
	"github.com/minio/minio-go"
	"github.com/rs/xid"
	// gominio "github.com/minio/minio-go"
)

func main() {
	fmt.Println("IGSN sample generator")
	rand.Seed(time.Now().UnixNano())

	// get a minio connection
	// build sample object
	// loadToMinio
	count := 1000000
	bar := uiprogress.AddBar(count).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Task (%d/%d)", b.Current(), count)
	})
	uiprogress.Start() // start rendering

	// mc := minioConnection("clear.local", "9000", "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	mc := minioConnection("localhost", "9111", "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")

	// Set up the the semaphore and conccurancey
	semaphoreChan := make(chan struct{}, 100) // a blocking channel to keep concurrency under control
	defer close(semaphoreChan)
	wg := sync.WaitGroup{} // a wait group enables the main process a wait for goroutines to finish

	for n := 0; n <= count; n++ {
		wg.Add(1)

		go func(n int) {
			semaphoreChan <- struct{}{}

			guid := xid.New()
			//fmt.Println(randFloats(-90, 90, 6))
			//fmt.Println(randFloats(-180, 180, 6))
			// fmt.Println(randomdata.Paragraph())

			jld := newRandomSample(guid.String()) //  send in some XID value for ID..   send in some random lat longs
			b := bytes.NewBuffer([]byte(jld))
			// load into minio

			contentType := "application/ld+json" // really Nq right now
			//n, err := mc.PutObject("doa-meta", objectName, b, int64(b.Len()), minio.PutObjectOptions{ContentType: contentType, UserMetadata: usermeta})
			_, err := mc.PutObject("samplesearth", guid.String(), b, int64(b.Len()), minio.PutObjectOptions{ContentType: contentType})
			// log.Printf("Loading metadata object: %d\n", n)
			if err != nil {
				log.Printf("Error loading metadata object to minio bucket %d, %s : %s\n", n, "igsntest", err)
			}

			wg.Done()
			bar.Incr()
			<-semaphoreChan
		}(n)
	}
	wg.Wait()

	time.Sleep(time.Second)
	uiprogress.Stop()

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

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)

	}
	return res
}

func newRandomSample(guid string) string {
	var buf bytes.Buffer

	t := template.Must(template.New("template").Parse(s))
	err := t.Execute(&buf, guid)
	if err != nil {
		fmt.Println(err)
	}

	return buf.String()
}

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
		  "@value": "POINT(146.067917 -34.79847)"
		}
	  },
	  {
		"@id": "http://sample.igsn.org/soilarchive/{{.}}",
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
		"http://purl.org/dc/terms/title": "ANZ soil sample {{.}}",
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
		  "@value": "1959-10-08"
		},
		"http://schema.org/identifier": "soil_specimen_{{.}}",
		"http://schema.org/title": "ANZ soil sample",
		"http://schema.org/url": {
		  "@id": "http://samples.earth/id/{{.}}"
		},
		"http://www.w3.org/2000/01/rdf-schema#label": "ANZ soil sample",
		"http://www.w3.org/ns/dcat#landingPage": {
		  "@id": "http://samples.earth/doc/{{.}}"
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
