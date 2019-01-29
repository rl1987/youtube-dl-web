package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fs := http.FileServer(http.Dir("public_html"))
	http.Handle("/", fs)

	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query()["video_url"]) < 1 {
			http.Error(w, "video_url parameter not found", 500)
			return
		}

		url := r.URL.Query()["video_url"][0]

		// First, let's figure out filename
		cmd := exec.Command("youtube-dl", "--get-filename", url)
		outbuf, err := cmd.Output()
		if err != nil {
			spew.Dump(err)
			http.Error(w, "Cannot determine file name", 500)
			return
		}

		outputFilename := string(outbuf)

		w.Header().Set("Content-Disposition", "attachment; filename="+outputFilename)

		cmd = exec.Command("youtube-dl", "--no-part", "-o", "-", url)

		spew.Dump(cmd)

		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			spew.Dump(err)
			http.Error(w, "Internal error", 500)
			return
		}

		go cmd.Run() // XXX: disregarding errors for now

		w.WriteHeader(200)

		io.Copy(w, outPipe)
	})

	log.Println("Listening for HTTP requests...")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
