package main

import (
	"html/template"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Name          string
	Size          float64
	Updated       time.Time
	Width, Height int
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func main() {
	port := ":3000"
	log.Println("Starting Web Server 127.0.0.1" + port)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("tmp"))))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fileFormUrlGet := r.URL.Query().Get("file")

	// Delete a file
	if fileFormUrlGet != "" {
		err := os.Remove("tmp/" + fileFormUrlGet)

		if err != nil {
			log.Println(err)
			setFlash(w, "warning", "The file '"+fileFormUrlGet+"' can't be remove")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		setFlash(w, "success", "The file '"+fileFormUrlGet+"' has been remove successul")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles("views/form_upload.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entries, err := ioutil.ReadDir("tmp/")

	if err == nil {
		files := []FileInfo{}

		for _, entry := range entries {
			// Open the file
			dimensions, err := os.Open("tmp/" + entry.Name())

			if err != nil {
				log.Println(err)
			}

			// To get dimensions in pixels
			imgConf, _, err := image.DecodeConfig(dimensions)

			if err != nil {
				log.Println(err)
			}

			// File informations
			f := FileInfo{
				Name:    entry.Name(),
				Size:    math.Ceil(float64(entry.Size()) / 1024),
				Updated: entry.ModTime(),
				Width:   imgConf.Width,
				Height:  imgConf.Height,
			}

			dimensions.Close() // Very important for "os.Remove"

			files = append(files, f)
		}

		// Cookies messages
		cookieSuccess, a := r.Cookie("success")
		cookieWarning, b := r.Cookie("warning")

		success := ""
		warning := ""

		if a == nil {
			success = cookieSuccess.Value
		}

		if b == nil {
			warning = cookieWarning.Value
		}

		data := map[string]interface{}{
			"Files":   files,
			"Success": success,
			"Warning": warning,
		}

		tmpl.Execute(w, data)

	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, fileheader, err := r.FormFile("file")

	// No file in the input
	if err != nil {
		log.Println(err)
		setFlash(w, "warning", "No file to upload")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	defer file.Close()

	// Size limitation
	var maxSize int64 = 2000000 // 2 Mb = 2 * 10^6 bytes

	if r.ContentLength > maxSize {
		setFlash(w, "warning", "The file '"+fileheader.Filename+"' can't be uploaded: too heavy size")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	extension := filepath.Ext(fileheader.Filename)

	// Files extensions allowed
	if extension == ".gif" || extension == ".jpg" || extension == ".png" {
		out, err := os.Create("tmp/" + fileheader.Filename)

		if err != nil {
			log.Println(err)
		}

		defer out.Close()

		_, err = io.Copy(out, file)

		if err != nil {
			log.Println(err)
		}

		defer out.Close()

		setFlash(w, "success", "The file '"+fileheader.Filename+"' has been upload successul")
	} else {
		setFlash(w, "warning", "Extension "+extension+" not accepted")
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Cookie message
func setFlash(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: 1}
	http.SetCookie(w, cookie)
}
