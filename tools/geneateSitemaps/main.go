package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/minio/minio-go"
	"github.com/tidwall/gjson"
)

type SiteMapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

type Sitemap struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
}

type URLSet struct {
	XMLName xml.Name  `xml:"urlset"`
	URLs    []URLNode `xml:"url"`
}

type URLNode struct {
	Loc        string  `xml:"loc"`
	Lastmod    string  `xml:"lastmod"`
	Changefreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

func main() {
	fmt.Println("samples.Earth sitemap builder")

	// TODO make the s3Maker and then time both of these
	start := time.Now()
	s3Maker()
	elapsed := time.Since(start)
	log.Printf("s3Maker took %s", elapsed)

	// start = time.Now()
	// kvMaker()
	// elapsed = time.Since(start)
	// log.Printf("kvMaker took %s", elapsed)

}

func s3Maker() {
	// mc := minioConnection("clear.local", "9000", "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	mc := minioConnection("localhost", "9111", "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")

	// Create a done channel.
	doneCh := make(chan struct{})
	defer close(doneCh)
	// Recursively list all objects in 'mytestbucket'
	recursive := true
	for message := range mc.ListObjectsV2("samplesearth", "", recursive, doneCh) {
		fmt.Println(message)
	}
}

func kvMaker() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("../../datalocal/my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket"))

		una := []URLNode{}
		count := 0 // sitemap count

		b.ForEach(func(k, v []byte) error {
			loc := fmt.Sprintf("http://samples.earth/id/%s", k)
			lastmod := getDate(v)
			changefreq := "yearly"
			priority := 0.8

			un := URLNode{loc, lastmod, changefreq, priority}
			una = append(una, un)

			if len(una) == 40000 {
				fmt.Println(len(una))
				err := writeSitemap(una, count)
				count = count + 1
				if err != nil {
					log.Println(err)
				} else {
					una = nil
				}
			}
			// fmt.Printf("key=%s, value=%d\n", k, len(v))
			return nil
		})
		fmt.Println(len(una))
		err := writeSitemap(una, count)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Writing Sitemap Index for %d files\n", count+1)
		err = sitemapindex(count)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
}

func getDate(v []byte) string {
	d := gjson.Get(string(v), "dateCreated.@value")
	log.Println(d.String())
	return d.String()
}

func sitemapindex(c int) error {
	sa := []Sitemap{}

	for i := 0; i <= c; i++ {
		now := time.Now()
		s := Sitemap{fmt.Sprintf("http://samples.earth/sitemap%d.xml", i), fmt.Sprint(now)}
		sa = append(sa, s)
	}

	smi := SiteMapIndex{Sitemaps: sa}

	filename := "./output/sitemap.xml"
	file, _ := os.Create(filename)
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")

	if err := enc.Encode(smi); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return nil
}

func writeSitemap(una []URLNode, count int) error {
	us := URLSet{URLs: una}

	filename := fmt.Sprintf("./output/sitemap%d.xml", count)
	file, _ := os.Create(filename)
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(us); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return nil
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
