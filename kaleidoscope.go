package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type img struct {
	Fname   string `json:"filename"`
	DataURL string `json:"data,omitempty"`
}

type orientation struct {
	Landscape bool
	Times     []byte
}

var host = flag.String("h", "localhost", "host")
var port = flag.String("p", "8080", "port")
var dir = flag.String("d", "D:\\Pictures", "dir to scan")
var landscape = flag.Bool("o", true, "landscape(true) portrait(false)")
var imgFileRe = regexp.MustCompile("jpe?g$")
var images = make([]*img, 0, 65536)

var landscapeOrientation = orientation{true, make([]byte, 16)}
var portraitOrientation = orientation{false, make([]byte, 12)}

//go:embed index.tmpl assets
var content embed.FS

func buildIndex(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to stat %s: %s", path, err)
		} else if info.Mode().IsRegular() && imgFileRe.MatchString(path) {
			images = append(images, &img{Fname: path})
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	flag.Parse()

	var index *template.Template
	if data, err := content.ReadFile("index.tmpl"); err != nil {
		log.Fatal(err)
	} else {
		index = template.Must(template.New("").Parse(string(data)))
	}

	buildIndex(*dir)
	rand.Seed(time.Now().Unix())

	fsHandler := http.FileServer(http.FS(content))

	http.HandleFunc("/random", serveRandomImage)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			or := portraitOrientation
			if *landscape {
				or = landscapeOrientation
			}

			if err := index.Execute(w, or); err != nil {
				http.Error(w, "failed: "+err.Error(), http.StatusInternalServerError)
			}
		} else {
			fsHandler.ServeHTTP(w, req)
		}
	})
	log.Fatal(http.ListenAndServe(net.JoinHostPort(*host, *port), nil))
}

func serveRandomImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json")

	img := images[rand.Intn(len(images))]

	if img.DataURL == "" {
		data, err := os.ReadFile(img.Fname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		img.DataURL = `data:image/jpeg;base64,` + base64.StdEncoding.EncodeToString(data)
	}

	if err := json.NewEncoder(w).Encode(img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
