package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	message              string = "static-server is starting on the port: %v"
	listenport           string = ":%v"
	baseTemplatePath     string = "template/base.html"
	showDirTemplatePath  string = "template/dir.html"
	notFoundTemplatePath string = "template/404.html"
)

var showDirTemplate *template.Template
var notFoundTemplate *template.Template

func handler(w http.ResponseWriter, r *http.Request) {

	pathname := strings.Trim(r.URL.Path, "/")
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	if pathname == "" {
		ShowDir(w, wd)
		return
	}

	paths := strings.Split(pathname, "/")

	for i, v := range paths {

		isLast := len(paths) == i+1

		fileInfos, err := ioutil.ReadDir(wd)

		if err != nil {
			log.Fatal(err)
		}
		for _, file := range fileInfos {

			filename := file.Name()
			basename := path.Base(filename)

			if basename != v {
				continue
			}

			if file.IsDir() {
				wd = fmt.Sprintf("%s/%s", wd, v)
				if isLast {
					ShowDir(w, wd)
					return
				}
				break
			}

			fullpath := path.Join(wd, filename)
			data, err := ioutil.ReadFile(fullpath)

			if err != nil {
				log.Fatal(err)
				return
			}

			w.Write(data)
			return
		}
	}
	NotFound(w)
}

func ShowDir(w http.ResponseWriter, wd string) {
	fileInfos, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}
	list := make([]string, len(fileInfos))
	for i, file := range fileInfos {
		filename := file.Name()
		basename := path.Base(filename)
		if file.IsDir() {
			basename += "/"
		}
		list[i] = basename
	}
	data := struct {
		Title string
		Items []string
	}{
		Title: fmt.Sprintf("Directory listing for %s/", wd),
		Items: list,
	}
	showDirTemplate.ExecuteTemplate(w, "base", data)
}

func NotFound(w http.ResponseWriter) {
	notFoundTemplate.ExecuteTemplate(w, "base", nil)
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)

	showDirTemplate, _ = template.ParseFiles(baseTemplatePath, showDirTemplatePath)
	notFoundTemplate, _ = template.ParseFiles(baseTemplatePath, notFoundTemplatePath)
}

func main() {

	port := flag.Int("port", 8000, "Port Number")
	flag.Parse()

	log.Info(fmt.Sprintf(message, *port))

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(listenport, *port), nil)
}
