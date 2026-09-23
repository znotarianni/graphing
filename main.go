package main

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed index.html
var page []byte

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func main() {
	dir := env("DATA_DIR", "/data")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(page)
	})

	// http.FileServer handles Range requests, so large files can be read in slices.
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir(dir))))

	mux.HandleFunc("/api/files", func(w http.ResponseWriter, r *http.Request) {
		type entry struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		}
		out := []entry{}
		filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			switch strings.ToLower(filepath.Ext(p)) {
			case ".csv", ".tsv", ".txt", ".json", ".h5", ".hdf5":
				if info, e := d.Info(); e == nil {
					rel, _ := filepath.Rel(dir, p)
					out = append(out, entry{filepath.ToSlash(rel), info.Size()})
				}
			}
			return nil
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	})

	addr := ":" + env("PORT", "8080")
	log.Printf("plotter listening on %s, data dir %s", addr, dir)
	log.Fatal(http.ListenAndServe(addr, mux))
}
