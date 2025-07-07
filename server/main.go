package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Item はデータモデル
type Item struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func main() {
	log.Println("Starting server...")
	db, err := gorm.Open(sqlite.Open("/tmp/db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&Item{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// host index.html
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var list []Item
			if err := db.Find(&list).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(list)

		case "POST":
			var it Item
			if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := db.Create(&it).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(it)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
