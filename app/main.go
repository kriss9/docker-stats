package main

import (
	"net/http"
	"path"
	"io"
	"log"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// try to filter files to download, example only
	r.HandleFunc("/download/{name:[a-z]+[0-9]+}/{ext}", DownloadHandler)

	// set content type to html example
	r.HandleFunc("/get/{html}", GetHandler).Methods("GET")

	//-
	http.Handle("/", r)

	// wait for clients
	http.ListenAndServe(":8080", nil)
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)

	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()

	// r.HandleFunc("/get/{html}", GetHandler).Methods("GET")

	vars := mux.Vars(req)
	html := vars["html"]

	// use path.Base() and append '.html', might prevent directory traversal attack :)
	fpath := "./static/html/" + path.Base(html+".html")

	res.Header().Set("Content-Type", "text/html")
	if err = servefile(res, fpath); nil != err {
		status = http.StatusInternalServerError
		return
	}
}

func servefile(res http.ResponseWriter, fpath string) (err error) {
	outfile, err := os.OpenFile(fpath, os.O_RDONLY, 0x0444)
	if nil != err {
		return
	}

	// 32k buffer copy
	written, err := io.Copy(res, outfile)
	if nil != err {
		return
	}

	log.Println("served file:", outfile.Name(), ";length:", written)
	return
}

func DownloadHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)

	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()

	// r.HandleFunc("/download/{name:[a-z]+[0-9]+}/{ext}", DownloadHandler)

	//vars := mux.Vars(req)
	//name := vars["name"]
	//ext := vars["ext"]

	// use path.Base(), might prevent directory traversal attack
	fpath := "main.go"
	//fpath := "./static/fileserver/" + path.Base(name+"."+ext)

	if err = servefile(res, fpath); nil != err {
		status = http.StatusInternalServerError
		return
	}

}

