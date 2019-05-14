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
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type myJSON struct {
	Array []string
}
var dbName string
var dbPort string

func main() {
	// TODO: Authenication
	dbName = "127.0.0.1"
	dbPort = "3306"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/organizations/{org}/groups/{group}", OrgGroupShow)
	router.HandleFunc("/organizations/{org}/groups", OrgGroups)
	router.HandleFunc("/organizations/{org}", OrgShow)
	router.HandleFunc("/organizations", OrgIndex)
	router.HandleFunc("/", OrgIndex)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func OrgIndex(w http.ResponseWriter, r *http.Request) {
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
	results.Close()
	db.Close()

	jsonPrint(w, organizations)
        return
}

func OrgShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	db := dbConnection(dbName, dbPort)
	// DB query to get all the groups in an org
	results, err := db.Query("SELECT name FROM organizations where name = '" + org + "';")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var orgs []string
	for results.Next() {
		var name string
		err = results.Scan(&name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		orgs = append(orgs, name)
	}
	results.Close()
	db.Close()
	orgs = unique(orgs)
	jsonPrint(w, orgs)
        return
}

func OrgGroups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	db := dbConnection(dbName, dbPort)
	// DB query to get all the groups in an org
	results, err := db.Query("SELECT group_name FROM org_groups where organization_name = '" + org + "';")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var groups []string
	for results.Next() {
		var name string
		err = results.Scan(&name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		groups = append(groups, name)
	}
	results.Close()
	db.Close()
	groups = unique(groups)
	jsonPrint(w, groups)
        return
}

func OrgGroupShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	group := cleanInput(vars["group"])
	db := dbConnection(dbName, dbPort)
	// DB query to get all the members in a group in an org
	results, err := db.Query("SELECT user_name FROM org_groups where organization_name = '" + org + "' AND group_name = '" + group + "';")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var members []string
	for results.Next() {
		var name string
		err = results.Scan(&name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		members = append(members, name)
	}
	results.Close()
	db.Close()
	members = unique(members)
	jsonPrint(w, members)
        return
}

func cleanInput(in string) string {
	match, err := regexp.MatchString("^[[:word:]]+$", in)
        if err != nil {
                fmt.Println(err.Error())
        }
	var out string
        if match {
		out = in
		return out
	}
	return "Invalid-Request"
}

func dbConnection(dbname string, dbport string) *sql.DB {
        db, err := sql.Open("mysql", "root@tcp("+dbname+":"+dbport+")/organizations")
        if err != nil {
                fmt.Println(err.Error())
        }
        return db
}

func jsonPrint(w http.ResponseWriter, out []string) {
	// turn it into json and return it
	jsondat := &myJSON{Array: out}
	encjson, _ := json.Marshal(jsondat)
	fmt.Fprintf(w, "%q", string(encjson))
}

func unique(in []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range in {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
