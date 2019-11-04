package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	. "./config"
	. "./dao"
	. "./models"
	"github.com/fasthttp/router"
	"github.com/gorilla/mux"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"gopkg.in/mgo.v2/bson"
)

var config = Config{}
var dao = MoviesDAO{}

// GET list of movies
func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
	movies, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movies)
}

// GET a movie by its ID
func FindCardsByExternalCodeEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindByExternalCode(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func FindCardsByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindByIdCard(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

// POST a new movie
func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Card

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, movie)
}

func CreateMovieEndPointStandalone(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Card
	movie.IndicePk = "card_id"
	movie.CoverImage = "1234"
	movie.Description = "100"
	movie.ID = bson.NewObjectId()

	if err := dao.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var movie2 Card
	movie2.IndicePk = "external_code"
	movie2.CoverImage = "1xxx"
	movie2.Description = "99"
	movie2.ID = bson.NewObjectId()

	if err := dao.Insert(movie2); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, movie2)
}

// PUT update an existing movie
func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Card
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dao.Update(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing movie

func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Card

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dao.Delete(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	http.DefaultClient.Timeout = time.Minute * 10
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes

func mainx() {
	r := mux.NewRouter()
	r.HandleFunc("/externalCode/{id}", FindCardsByExternalCodeEndPoint).Methods("GET")
	r.HandleFunc("/card/{id}", FindCardsByIdEndPoint).Methods("GET")
	r.HandleFunc("/add", CreateMovieEndPointStandalone).Methods("GET")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func FindCardsByIdEndPoint2(ctx *fasthttp.RequestCtx) {
	movie, err := dao.FindCartao()
	if err != nil {
		return
	}
	response, _ := json.Marshal(movie)
	ctx.Write(response)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func FindCardsByIdEndPoint3(ctx *fasthttp.RequestCtx) {
	movie, err := dao.FindExternalCode()
	if err != nil {
		return
	}
	response, _ := json.Marshal(movie)
	ctx.Write(response)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func main() {

	ln, err := reuseport.Listen("tcp4", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	r := router.New()
	r.GET("/card/10", FindCardsByIdEndPoint2)
	r.GET("/externalCode/1xxx", FindCardsByIdEndPoint3)

	if err = fasthttp.Serve(ln, r.Handler); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}

}
