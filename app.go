package main

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"
	"strconv"
	"encoding/json"

	"github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

// The Initialize method will take in the details required to connect to the database.
// It will create a database connection and wire up the routes to respond according to the requirements.
// export APP_DB_USERNAME=postgres
// export APP_DB_PASSWORD=
// export APP_DB_NAME=postgres
func (a *App) Initialize(user, password, dbname string) {
    connectionString := fmt.Sprintf("user=postgres password=root dbname=postgres sslmode=disable")

    var err error
    a.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()
		a.initializeRoutes()
}

// The Run method will simply start the application.
func (a *App) Run(addr string) { 
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

// This handler retrieves the id of the product to be fetched from the requested URL, and uses the getProduct method, created in the previous section, to fetch the details of that product.
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid product ID")
			return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
					respondWithError(w, http.StatusNotFound, "Product not found")
			default:
					respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// A handler to fetch a list of products
// This handler uses the count and start parameters from the querystring to fetch count number of products, starting at position start in the database. By default, start is set to 0 and count is set to 10. If these parameters aren’t provided, this handler will respond with the first 10 products.
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
			count = 10
	}
	if start < 0 {
			start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// This handler assumes that the request body is a JSON object containing the details of the product to be created. It extracts that object into a product and uses the createProduct method to create a product with these details.
func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid product ID")
			return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
			return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
			return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}