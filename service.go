package main

import (
	"net/http"
	"os"
	"fmt"
	"regexp"
	"log"
	"mime"
	"path/filepath"
	"strconv"
)

const (
	rootPath = "/"
	subscribePath = "/subscribe"
	emailField = "email"
	emailRegexp = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

func root(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if (path == "") {
		path = "index.html"
	}
	data, contentType, err := getResource("static/" + path)
	if (err != nil) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func getResource(path string) ([]byte, string, error) {
	extension := filepath.Ext(path)
	contentType := mime.TypeByExtension(extension)
	if (contentType == "") {
		contentType = "text/plain"
	}
	data, err := Asset(path)
	if (err != nil) {
		return nil, "", err
	}
	return data, contentType, nil
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get(emailField)
	if !isValidEmail(email) {
		log.Printf("Skipping invalid email [%s]\n", email)
		http.Error(w, "no valid email provided", http.StatusBadRequest)
		return
	}
	err := saveEmail(email)
	if (err != nil) {
		log.Printf("Server error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func isValidEmail(email string) bool {
	if (len(email) == 0) {
		return false
	}
	regex := regexp.MustCompile(emailRegexp)
	return regex.MatchString(email)
}

func saveEmail(email string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n", email)); err != nil {
		return err
	}
	return nil
}

func addCommonHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Almighty 1.0 with BlackJack Patch 1.1, Hookers Patch 2.1")
		fn(w, r)
	}
}

func mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(rootPath, addCommonHeaders(root))
	mux.HandleFunc(subscribePath, addCommonHeaders(subscribe))
	return mux
}
