//
// Chef group rest api
//
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type myJSON struct {
	Array []string
}

type myJSONUser struct {
	Array []userInfo
}

type userInfo struct {
	Name string
	Email string
	Display string
}

var dbName string
var dbPort string
var dbUser string
var dbPwdFile string

func main() {
	// TODO: Authenication
	dbName = "127.0.0.1"
	dbPort = "3306"
	dbPwdFile = os.Args[1]
	dbUser = "root"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", UserIndex)
	router.HandleFunc("/users/{member}", UserShow)
	router.HandleFunc("/organizations/{org}/groups/{group}", OrgGroupShow)
	router.HandleFunc("/organizations/{org}/groups", OrgGroups)
	router.HandleFunc("/organizations/{org}", OrgShow)
	router.HandleFunc("/organizations", OrgIndex)
	router.HandleFunc("/", OrgIndex)
	l, err := net.Listen("tcp4", ":8111")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Fatal(http.ServeTLS(l, router, "cert.pem" , "key.pem"))
}

func OrgIndex(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	db := dbConnection(dbName, dbPort, dbUser, dbPwd(dbPwdFile))
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
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	db := dbConnection(dbName, dbPort, dbUser, dbPwd(dbPwdFile))
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
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	db := dbConnection(dbName, dbPort, dbUser, dbPwd(dbPwdFile))
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
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	group := cleanInput(vars["group"])
	db := dbConnection(dbName, dbPort, dbUser, dbPwd(dbPwdFile))
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

func UserIndex(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	db := dbConnection(dbName, dbPort, dbUser, dbPwd(dbPwdFile))
	// DB query to get all the users
	results, err := db.Query("SELECT user_name FROM members")
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

	jsonPrint(w, members)
        return
}

func UserShow(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	member := cleanInput(vars["member"])
	db := dbConnection(dbName, dbPort, dbUser, dbPwd(dbPwdFile))
	// DB query to get a specific member
	results, err := db.Query("SELECT user_name, email, display_name FROM members where user_name = '" + member + "';")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// TODO package the query results if any into a map
        var users []userInfo
	for results.Next() {
		var user userInfo
		err = results.Scan(&user.Name, &user.Email, &user.Display)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
            	users = append(users, user)
	}
	results.Close()
	db.Close()
	jsonPrintUser(w, users)
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

func dbConnection(dbname string, dbport string, dbuser string, dbpwd string) *sql.DB {
        db, err := sql.Open("mysql", dbuser+":"+dbpwd+"@tcp("+dbname+":"+dbport+")/organizations")
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

func jsonPrintUser(w http.ResponseWriter, out []userInfo) {
	// turn it into json and return it
	jsondat := &myJSONUser{Array: out}
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

func dbPwd(dbpwd_file string) string {
        pwd, err := ioutil.ReadFile(dbpwd_file)
        if err != nil {
                panic(err.Error()) // proper error handling instead of panic in your app
        }
        return strings.TrimSpace(string(pwd))
}
