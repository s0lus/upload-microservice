package main

import (
	_"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"encoding/json"
	_"image/jpeg"
	_"image/png"
	"io/ioutil"
	"net/http"
	"fmt"
	"os"
	"image"
	"path"
	"io"
)

func main() {
	var config Config

	jsonStream, error := ioutil.ReadFile("./config.json");
	if error != nil {
		fmt.Println("Error reading 'config.json' file")
	}

	json.Unmarshal(jsonStream, &config)

	db, error := gorm.Open(config.DataBase.Dialect, config.DataBase.ConnectionData)
	if error != nil {
		fmt.Println("Error connection to database")
	}

	defer db.Close()

	db.DB().SetMaxIdleConns(config.DataBase.IdleConnections)
	db.DB().SetMaxOpenConns(config.DataBase.MaxOpenConnections)

	db.LogMode(true);
	//db.AutoMigrate(&Image{})

	r := mux.NewRouter()

	r.HandleFunc("/upload/image/", postHandler(db)).Methods("POST")
	r.HandleFunc("/image/{id}", getHandler(db)).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}

func postHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Create flags for open images
		pathToImage := "./logo.png"

		// Parse image and find: Height, Width, Name
		fileToParse, error := os.Open(pathToImage)
		if error != nil {
			fmt.Println("Invalid path to image")
		}

		img, _, error := image.DecodeConfig(fileToParse)
		if error != nil {
			fmt.Println("Error decode image")
		}

		var image Image

		baseName := path.Base(pathToImage)

		image.Name = baseName
		image.Height = img.Height
		image.Width = img.Width

		// Add new record to database (record: info about image)
		newImage := Image{Name:image.Name, Width:image.Width, Height:image.Height}

		db.Create(&newImage)

		error = r.ParseForm()
		if error != nil {
			fmt.Println("Error parsing form")
		}

		r.ParseMultipartForm(32 << 20)

		// "Image" - name in form (key)
		file, _, error := r.FormFile("Image")
		if error != nil {
			fmt.Println("Error call FormFile()")
		}
		defer file.Close()

		f, error := os.OpenFile(pathToImage, os.O_WRONLY | os.O_CREATE, 0666)
		if error != nil {
			fmt.Println("Error reading 'logo.png' file")
		}
		defer f.Close()

		io.Copy(f, file)
	}
}

func getHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		imageID := params["id"]

		var image Image

		db.Where("id = ?", imageID).Find(&image)

		if image.ID == 0 {
			fmt.Println("Image doesn`t exist")
		}
	}

}