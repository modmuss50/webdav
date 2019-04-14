package webdav

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func (c *Config) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		c.receiveFile(w, r)
	} else {
		file := c.getPath(r.URL.Path)
		http.ServeFile(w, r, file)
	}
}

func (c *Config) receiveFile(w http.ResponseWriter, r *http.Request) {
	file := c.getPath(r.URL.Path)
	os.MkdirAll(filepath.Dir(file), os.ModePerm)
	defer r.Body.Close()
	out, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(out, r.Body)
}

func (c *Config) getPath(file string) string {
	return path.Join(c.Scope, file)
}
