package dao

import (
	"log"

	. "../models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MoviesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "cards_detail2"
)

// Establish a connection to database

func (m *MoviesDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of movies

func (m *MoviesDAO) FindAll() ([]Card, error) {
	var movies []Card
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

// Find a movie by its id

func (m *MoviesDAO) FindById(id string) (Card, error) {
	var movie Card
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
	return movie, err
}

func (m *MoviesDAO) FindByIdCard(id string) (Card, error) {

	query := bson.M{"indicePk": "card_id", "indiceSk": id}

	var movie Card
	err := db.C(COLLECTION).Find(query).One(&movie)
	return movie, err
}

func (m *MoviesDAO) FindCartao() (Card, error) {
	query := bson.M{"indicePk": "card_id", "indiceSk": "10"}

	var movie Card
	err := db.C(COLLECTION).Find(query).One(&movie)
	return movie, err
}

func (m *MoviesDAO) FindExternalCode() (Card, error) {

	query := bson.M{"indicePk": "external_code", "indiceSk": "1xxx"}

	var movie Card
	err := db.C(COLLECTION).Find(query).One(&movie)
	return movie, err
}

func (m *MoviesDAO) FindByExternalCode(id string) (Card, error) {

	query := bson.M{"indicePk": "external_code", "indiceSk": id}

	var movie Card
	err := db.C(COLLECTION).Find(query).One(&movie)
	return movie, err
}

// Insert a movie into database

func (m *MoviesDAO) Insert(movie Card) error {
	err := db.C(COLLECTION).Insert(&movie)
	return err
}

// Delete an existing movie

func (m *MoviesDAO) Delete(movie Card) error {
	err := db.C(COLLECTION).Remove(&movie)
	return err
}

// Update an existing movie
func (m *MoviesDAO) Update(movie Card) error {
	err := db.C(COLLECTION).UpdateId(movie.ID, &movie)
	return err
}
