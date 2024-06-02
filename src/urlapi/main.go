package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"url/pkg/mongodb"

	"url/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Link string
}

func main() {
	http.HandleFunc("/", urlRequest)
	http.ListenAndServe(":"+config.Config.ServerPort, nil)
}

func urlRequest(w http.ResponseWriter, r *http.Request) {
	shortUrl := strings.TrimPrefix(r.URL.Path, "/")
	switch {
	case shortUrl == "":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, here's URLAPI"))
		return

	case strings.HasPrefix(shortUrl, "!"):
		shortUrl := strings.TrimPrefix(shortUrl, "!")
		url, err := getUrlHandler(w, r, shortUrl)
		if err != nil && err != mongo.ErrNoDocuments {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
			return
		}
		message := Message{
			Link: url,
		}
		exePath, err := os.Executable()
		exeDir := filepath.Dir(exePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		renderTemplate(w, exeDir+"/template.html", message)
	default:
		url, err := getUrlHandler(w, r, shortUrl)
		if err != nil || err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	}
}

func getUrlHandler(w http.ResponseWriter, r *http.Request, shortUrl string) (string, error) {
	mongoClient, err := mongodb.MongoConnect(config.Config.DbHost, config.Config.DbPort, config.Config.DbUser, config.Config.DbPass)
	if err != nil {
		return "", err
	}
	defer mongoClient.Disconnect(context.TODO())

	longUrl, err := mongodb.GetKey(mongoClient, config.Config.CollectionName, config.Config.DbName, shortUrl)
	if err != nil {
		return "", err
	}

	return longUrl, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data Message) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
