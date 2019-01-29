package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

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

		tmpDir, err := ioutil.TempDir("/tmp", "youtube_dl_web")
		if err != nil {
			spew.Dump(err)
			http.Error(w, "Cannot create tmp dir", 500)
			return
		}
		defer os.Remove(tmpDir)

		outputTempl := tmpDir + "/%(title)s.%(ext)s"

		cmd := exec.Command("youtube-dl", "-o", outputTempl, url)

		spew.Dump(cmd)

		err = cmd.Run()
		if err != nil {
			spew.Dump(err)
			http.Error(w, "downloading failed", 500)
			return
		}

		files, err := ioutil.ReadDir(tmpDir)
		if err != nil || len(files) != 1 {
			spew.Dump(err)
			http.Error(w, "downloading failed - cannot find downloaded file", 500)
			return
		}

		outputFile := files[0]

		w.Header().Set("Content-Disposition", "attachment; filename="+outputFile.Name())
		w.Header().Set("Content-Length", strconv.FormatInt(outputFile.Size(), 10))

		filename := filepath.Join(tmpDir, outputFile.Name())

		openfile, err := os.Open(filename)
		if err != nil {
			spew.Dump(err)
			http.Error(w, "cannot open file", 500)
			return
		}

		defer openfile.Close()

		io.Copy(w, openfile)
	})

	log.Println("Listening for HTTP requests...")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
