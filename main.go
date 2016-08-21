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
	r.HandleFunc("/image/{id}", getHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}


func postHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Create flags for open images
		pathToImage := "./logo.png"

		// Parse image and find: Height, Width, Name
		file, error := os.Open(pathToImage)
		if error != nil {
			fmt.Println("Invalid path to image")
		}

		img, _, error := image.DecodeConfig(file)
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
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {

}