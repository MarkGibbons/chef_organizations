//
// Chef group rest api
//
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type myJSON struct {
	Array []string
}

func main() {
	// TODO: Authenication
	dbName := "127.0.0.1"
	dbPort := "3306"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", OrgIndex)
	router.HandleFunc("/organizations", OrgIndex)
	router.HandleFunc("/organizations/{org}", OrgShow)
	router.HandleFunc("/organizations/{org}/groups", OrgGroups)
	router.HandleFunc("/organizations/{org}/groups/{group}", OrgGroupShow)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func OrgIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Organizations Index")
	db := dbConnection(dbName, dbPort)
	// DB query to get all the orgs
	results, err := db.Query("SELECT name FROM organizations")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var organizations []string
	for results.Next() {
		var name string
		err = results.Scan(&name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		organizations = append(organizations, name)
	}
	results.Close
	db.Close

	// turn it into json and return it
	jsondat := &myJSON{Array: organizations}
	encjson, _ := json.Marshal(jsondat)
	fmt.Println(string(encjson))

}

func OrgShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	fmt.Fprintln(w, "Organization show "+org)
	results, err := db.Query("SELECT name FROM organizations  WHERE name = '" + org + "';")
	// Pass organization name
	// canonical form of the org name
	// matches?
}

func OrgGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Organization Groups ")
}

func OrgGroupShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	group := vars["group"]
	fmt.Fprintln(w, "Organization "+org+" Group "+group)
}

func dbConnection(dbname string, dbport string) *sql.DB {
        db, err := sql.Open("mysql", "root@tcp("+dbname+":"+dbport+")/organizations")
        if err != nil {
                fmt.Println(err.Error())
        }
        return db
}
