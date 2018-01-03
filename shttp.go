/*
   Copyright © 2017 Jose Gonzalez Krause (josef@hackercat.ninja)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "127.0.0.1:8080", "Sock to listen to")
var servePath = flag.String("path", "./", "Path to serve")
var allowUpload = flag.Bool("upload", false, "Allow file upload")

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir(*servePath))
	http.Handle("/", fs)

	log.Printf("Listening on http://%s", *addr)
	log.Printf("Serving files from: %s", *servePath)
	if *allowUpload {
		log.Println("Allowing file upload")
		http.HandleFunc("/upload", upload)
	}

	http.ListenAndServe(*addr, nil)
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.New("upload")
		t.Parse(fmt.Sprintf(uploadHTML, *addr))
		t.Execute(w, nil)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		fmt.Fprintf(w, "Uploaded: %s", handler.Filename)

		f, err := os.OpenFile(*servePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
	}
}

var uploadHTML = `
<html>
<head>
       <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="http://%s/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="submit" value="upload" />
</form>
</body>
</html>
`
