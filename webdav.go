package main

import (
	"log"
	"net/http"
	"strings"
)

// Config is the configuration of a WebDAV instance.
type Config struct {
	*User
	Users map[string]*User
}

// ServeHTTP determines if the request is for this plugin, and if all prerequisites are met.
func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := c.User

	// Gets the correct user for this request.
	username, _, ok := r.BasicAuth()
	if ok {
		if user, ok := c.Users[username]; ok {
			u = user
		}
	}

	// Checks for user permissions relatively to this PATH.
	if !u.Allowed(r.URL.Path) {
		http.Error(w, "Not authorized", http.StatusForbidden)
		log.Println("User not authorized for path", username)
		return
	}

	if r.Method == "HEAD" {
		w = newResponseWriterNoBody(w)
	}

	// If this request modified the files and the user doesn't have permission
	// to do so, return forbidden.
	if (r.Method == "PUT" || r.Method == "POST" || r.Method == "MKCOL" ||
		r.Method == "DELETE" || r.Method == "COPY" || r.Method == "MOVE") &&
		!u.Modify {
		http.Error(w, "Not authorized", http.StatusForbidden)
		return
	}

	c.HandleRequest(w, r)
}

// Rule is a dissalow/allow rule.
type Rule struct {
	Path string
}

// User contains the settings of each user.
type User struct {
	Scope  string
	Modify bool
	Rules  []*Rule
}

// Allowed checks if the user has permission to access a directory/file
func (u User) Allowed(url string) bool {
	var rule *Rule
	i := len(u.Rules) - 1

	for i >= 0 {
		rule = u.Rules[i]

		if strings.HasPrefix(url, rule.Path) {
			return true
		}

		i--
	}

	return false
}

// responseWriterNoBody is a wrapper used to suprress the body of the response
// to a request. Mainly used for HEAD requests.
type responseWriterNoBody struct {
	http.ResponseWriter
}

// newResponseWriterNoBody creates a new responseWriterNoBody.
func newResponseWriterNoBody(w http.ResponseWriter) *responseWriterNoBody {
	return &responseWriterNoBody{w}
}

// Header executes the Header method from the http.ResponseWriter.
func (w responseWriterNoBody) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write suprresses the body.
func (w responseWriterNoBody) Write(data []byte) (int, error) {
	return 0, nil
}

// WriteHeader writes the header to the http.ResponseWriter.
func (w responseWriterNoBody) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
