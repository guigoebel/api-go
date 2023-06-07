package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/guigoebel/api-go/configs"
	"github.com/guigoebel/api-go/internal/entity"
	"github.com/guigoebel/api-go/internal/infra/database"
	"github.com/guigoebel/api-go/internal/infra/webservice/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDb := database.NewProduct(db)
	ProductHandler := handlers.NewProductHandler(productDb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/products", ProductHandler.GetProducts)
	r.Get("/products/{id}", ProductHandler.Get)
	r.Post("/products", ProductHandler.Create)
	r.Put("/products/{id}", ProductHandler.Update)
	r.Delete("/products/{id}", ProductHandler.Delete)
	//endpoint for create product
	http.ListenAndServe(":8080", r)

}
