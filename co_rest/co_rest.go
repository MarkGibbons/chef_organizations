//
// Chef organization rest api
//

// TODO: REST documentation
package main

import (
        "co"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

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

var dbc co.DbConnectionRequest

// co_rest main routes REST requests
func main() {
	// TODO: Authenication
	dbc.PwdFile = os.Args[1]
	dbc.Server = "127.0.0.1"
	dbc.Port = "3306"
	dbc.User = "root"
        dbc.Database = "organizations"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", userIndex)
	router.HandleFunc("/users/{member}", userShow)
	router.HandleFunc("/organizations/{org}/groups/{group}", orgGroupShow)
	router.HandleFunc("/organizations/{org}/groups", orgGroups)
	router.HandleFunc("/organizations/{org}", orgShow)
	router.HandleFunc("/organizations", orgIndex)
	router.HandleFunc("/", orgIndex)
	l, err := net.Listen("tcp4", ":8111")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Fatal(http.ServeTLS(l, router, "cert.pem" , "key.pem"))
}

// Return a list of organizations
func orgIndex(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	db := co.DbConnection(dbc)
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

// orgShow executes a DB query to get a specific organization
func orgShow(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	db := co.DbConnection(dbc)
        stmtQryOrg, err := db.Prepare("SELECT name FROM organizations where name = ? ;)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	results, err := stmtQryOrg.Exec(org)
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
        stmtQryOrg.Close()
	db.Close()
	orgs = co.Unique(orgs)
	jsonPrint(w, orgs)
        return
}

// orgGroups executes a DB query to get all the groups in an org
func orgGroups(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	db := co.DbConnection(dbc)
	stmtQryOrgGrp, err := db.Prepare("SELECT group_name FROM org_groups where organization_name = ? ;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	results, err := stmtQryOrgGrp.Exec(org)
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
        stmtQryOrgGrp()
	db.Close()
	groups = co.Unique(groups)
	jsonPrint(w, groups)
        return
}

// orgGroupShow executes a DB query to get the users in an organization and group.
func orgGroupShow(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	org := cleanInput(vars["org"])
	group := cleanInput(vars["group"])
	db := co.DbConnection(dbc)
	stmtQryOrgGrp, err := db.Prepare("SELECT user_name FROM org_groups where organization_name = ? AND group_name = ? ;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	results, err := stmtQryOrgGrp.Exec(org, group)
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
        stmtQryOrgGrp.Close()
	db.Close()
	members = co.Unique(members)
	jsonPrint(w, members)
        return
}

// userIndex gets a list of all the users in the database members table.
func userIndex(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	db := co.DbConnection(dbc)
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

// userShow gets the details of a specific member.
func userShow(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	member := cleanInput(vars["member"])
	db := co.DbConnection(dbc)
	stmtQryNm, err := db.Prepare("SELECT user_name, email, display_name FROM members where user_name = ? ;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	results, err := stmtQryNm.Exec(member)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
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
        stmtQryNm.Close()
	db.Close()
	jsonPrintUser(w, users)
        return
}

// cleanInput restrict the values that can be specified via the rest interface
// Allow only word characters for the org, group and member names
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

// jsonPrint marshall json data from an array to return to the rest call
func jsonPrint(w http.ResponseWriter, out []string) {
	jsondat := &myJSON{Array: out}
	encjson, _ := json.Marshal(jsondat)
	fmt.Fprintf(w, "%q", string(encjson))
}

// jsonPrintUser marshal user data from an array to return to the rest call
func jsonPrintUser(w http.ResponseWriter, out []userInfo) {
	jsondat := &myJSONUser{Array: out}
	encjson, _ := json.Marshal(jsondat)
	fmt.Fprintf(w, "%q", string(encjson))
}
