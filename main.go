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
	"io"
	"image"
	"flag"
)

func main() {
	configFlag := flag.String("config", "./config.json", "Set path to config file for database")

	flag.Parse();

	var config Config

	jsonStream, err := ioutil.ReadFile(*configFlag);
	if err != nil {
		fmt.Println("Error reading 'config.json' file")
	}

	json.Unmarshal(jsonStream, &config)

	db, err := gorm.Open(config.DataBase.Dialect, config.DataBase.ConnectionData)
	if err != nil {
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
		r.ParseMultipartForm(32 << 20)

		files := r.MultipartForm.File["Image"]

		for i := range files {
			fileIn, err := files[i].Open()
			if err != nil {
				fmt.Println("Error opening files")
			}

			fileOut, err := os.Create("./images/" + files[i].Filename)
			if err != nil {
				fmt.Println("Error 'images' folder doesn`t exists")
			}

			_, err = io.Copy(fileOut, fileIn)
			if err != nil {
				fmt.Println("Error copying files")
			}
		}

		fmt.Fprintln(w, "Files upload successfully");

		var imageInfo Image

		for i := 0; i < len(files); i++ {
			fileToParse, err := files[i].Open()
			if err != nil {
				fmt.Println("Error opening files")
			}

			img, _, err := image.DecodeConfig(fileToParse)
			if err != nil {
				fmt.Println("Error parsing files")
			}

			files, err := ioutil.ReadDir("./images/")
			if err != nil {
				fmt.Println("Error reading ./images/ directory")
			}

			imageInfo.Name = files[i].Name()
			imageInfo.Height = img.Height
			imageInfo.Width = img.Width

			newImage := Image{Name:imageInfo.Name, Width:imageInfo.Width, Height:imageInfo.Height}
			db.Create(&newImage)
		}

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