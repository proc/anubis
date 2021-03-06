package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	valid "github.com/asaskevich/govalidator"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	ngr "github.com/jingweno/negroni-gorelic"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func HomePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Anubis")
}

type City struct {
	Name string
}

type Error struct {
	Message string
}

func Locate(db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var city string
		lat := r.FormValue("lat")
		lng := r.FormValue("lng")
		if !valid.IsLatitude(lat) || !valid.IsLongitude(lng) {
			e := new(Error)
			e.Message = "Missing Param"
			js, _ := json.Marshal(e)
			rw.Write(js)
			return
		}
		q := "SELECT name FROM planet_osm_polygon WHERE ST_DWithin(way, ST_TRANSFORM(ST_SETSRID(ST_MAKEPOINT($1, $2), 4326), 900913), 1) AND admin_level='8' AND name IS NOT NULL;"
		err := db.Get(&city, q, lng, lat)
		if err != nil {
			e := new(Error)
			e.Message = "No results."
			js, _ := json.Marshal(e)
			rw.Write(js)
			return
		}
		js, err := json.Marshal(&City{city})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	})
}
func NewDB(dbName string, dbUser string, dbPass string, dbPort string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user="+dbUser+" password="+dbPass+" dbname="+dbName+" sslmode=disable port="+dbPort)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func main() {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	newrelicKey := os.Getenv("NEWRELIC_KEY")
	db := NewDB(dbName, dbUser, dbPass, dbPort)
	defer db.Close()

	r := mux.NewRouter().StrictSlash(false)
	r.Handle("/locate", Locate(db))
	r.HandleFunc("/", HomePage)

	n := negroni.New()
	n.Use(ngr.New(newrelicKey, "anubis", true))
	n.UseHandler(r)
	n.Run(":3000")
}
