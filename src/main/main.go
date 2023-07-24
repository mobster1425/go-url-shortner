package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	// "github.com/gophercises/urlshort"
	"feyin/go-url-shortner/urlshort"
)

func main() {
	mux := defaultMux()

	// Read the YAML file from the command line flag
	yamlFilename := flag.String("yaml", "data/paths.yaml", "a YAML file containing the URL mappings")
	jsonFilename := flag.String("json", "data/paths.json", "a JSON file containing the URL mappings")
	dbPath := flag.String("db", "data/urlshort.db", "path to the BoltDB database file")
	flag.Parse()

	// Initialize the BoltDB and create a bucket for URL mappings
	err := urlshort.InitDB(*dbPath)
	if err != nil {
		panic(err)
	}

	// Build the MapHandler using the mux as the fallback
	//this on eis a different dataits not using yaml or json
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	// Build the YAMLHandler using the mapHandler as the fallback
	yamlFile, err := os.Open(*yamlFilename)
	if err != nil {
		panic(err)
	}
	defer yamlFile.Close()

	// Read the contents of the file into a byte slice
	yamlBytes, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yamlBytes, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open(*jsonFilename)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.JSONHandler(jsonData, yamlHandler)
	if err != nil {
		panic(err)
	}

	// Build the DBHandler using the jsonHandler as the fallback
	dbHandler := urlshort.DBHandler(jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

// ... (defaultMux and hello functions remain the same)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

/*
for each of these callbacks there is a layer, the mapHandler uses mux as the callback
while yaml handler uses mapHandler as the callback to handle incorrect path which redirets to mux to handle the path not found, so mapHandler directs to mux
same thing with jsonHandler it redirects the callback to yamlHandler so that would also redirects the callback to mapHandler
which redirects to mux to handle the paths not found
//same thing with the DbHandler





*/
