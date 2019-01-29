package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

		tmpDir, err := ioutil.TempDir("/tmp", "youtube_dl_web")
		if err != nil {
			log.Fatal(err)
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
	})

	log.Println("Listening for HTTP requests...")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
