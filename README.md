# samples.earth

## Programs

### Sample generator

Build a simple tool to make 2 million samples. 

### Sample presenter (simple web server with landing pages)

Build a simple web site with landing pages and a sitemap.xml

Then use Gleaner to harvest it and make an index.  

Need a "sampleRegister" that reads the sitemap and pulls the 
required registering information (and reads the time stamps)

### IGSN sample schema git repo

Pull down Simons git repo for samples. 

## About

This is a simple program to generate 1 million (or any number) of samples
to simulate an environment to test use of structured data on the web to support
IGSN Sample registration.  

## Sitemap

* limit 50K URL per sitemap
* URL can be another sitemap but only 50K again, so limit is 50K^2 or 2.5 B  
* To support 1 million samples will require 20 sitemaps in an index.  

The base sitemap is:

``` xml
<url>
      <loc>http://www.example.com/</loc>
      <lastmod>2005-01-01</lastmod>
      <changefreq>monthly</changefreq>
      <priority>0.8</priority>
</url>
```

Of these we will likely use only the first two though the others can have optional meaning to the 
approach.  

## Command notes

``` Go
go run main.go
IGSN sample generator
Task (1000000/1000000) 2h23m14s [====================================================================] 100%
```

## Docker notes

docker run -p 127.0.0.1:$HOSTPORT:$CONTAINERPORT --name CONTAINER -t someimage

