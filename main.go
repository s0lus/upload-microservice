package main

import (
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"github.com/jinzhu/gorm"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"net/http"
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

	db.DB().SetMaxIdleConns(config.DataBase.IdleConnections)
	db.DB().SetMaxOpenConns(config.DataBase.MaxOpenConnections)

	db.LogMode(true);
	db.AutoMigrate(&Image{})

	r := mux.NewRouter()

	r.HandleFunc("/upload/image", postHandler).Methods("POST")
	r.HandleFunc("/image/{id}", getHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {

}

func getHandler(w http.ResponseWriter, r *http.Request) {

}
