package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/guigoebel/api-go/configs"
	"github.com/guigoebel/api-go/internal/entity"
	"github.com/guigoebel/api-go/internal/infra/database"
	"github.com/guigoebel/api-go/internal/infra/webservice/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Guilherme Goebel
// @contact.url    https://www.linkedin.com/in/guilherme-goebel/
// @contact.email  guigoebel

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	configs, err := configs.LoadConfig(".")
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

	userDb := database.NewUser(db)
	UserHandler := handlers.NewUserHandler(userDb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))
	// r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", ProductHandler.GetProducts)
		r.Get("/{id}", ProductHandler.Get)
		r.Post("/", ProductHandler.Create)
		r.Put("/{id}", ProductHandler.Update)
		r.Delete("/{id}", ProductHandler.Delete)
	})

	r.Post("/users", UserHandler.Create)
	r.Post("/users/generate_token", UserHandler.GetJWT)

	//endpoint for create product
	http.ListenAndServe(":8080", r)

}

// func LogRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("Request: %s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }
