package urlshort

import (
	"net/http"
	"encoding/json"
	yaml "gopkg.in/yaml.v2"
	"github.com/boltdb/bolt" // Import the BoltDB package
)


var db *bolt.DB // Declare a global variable to hold the BoltDB instance





// Initialize the BoltDB and create a bucket for URL mappings.
// This function should be called from the main package to initialize the database.
func InitDB(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("urlMappings"))
		return err
	})
	if err != nil {
		return err
	}

	return nil
}





// DBHandler reads the URL mappings from the BoltDB and returns
// an http.HandlerFunc that attempts to map paths to their corresponding URLs.
func DBHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var dest string
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("urlMappings"))
			val := b.Get([]byte(path))
			if val != nil {
				dest = string(val)
			}
			return nil
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if dest != "" {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}





// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}
//here we are turning it into a string and mapping it 
//we are extracting the path which is the key and the url which is the value 

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}


//we are passing pathurl which is of yaml to bytes and sending it to yaml handler 
//when the  data is parsed to the pathurl array the path and url will turn to string
func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}




// JSONHandler will parse the provided JSON data and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON data, then the
// fallback http.Handler will be called instead.
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJSON(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := json.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}



type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
