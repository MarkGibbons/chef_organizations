//
// Chef group rest api
//
package main

import (
	"fmt"
	"log"
	"net/http"

        "github.com/gorilla/mux"
)

func main() {
    // Authenication ?  
    
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", OrgIndex)
    router.HandleFunc("/organizations", OrgIndex)
    router.HandleFunc("/organizations/{org}", OrgShow)
    router.HandleFunc("/organizations/{org}/groups", OrgGroups)
    router.HandleFunc("/organizations/{org}/groups/{group}", OrgGroupShow)
    log.Fatal(http.ListenAndServe(":8080", router))
}

func OrgIndex( w http.ResponseWriter, r *http.Request)  {
  fmt.Fprintln(w, "Organizations Index")
  // DB query to get all the orgs
}

func OrgShow( w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  org := vars["org"]
  fmt.Fprintln(w, "Organization show " + org)
  // Pass organization name
  // canonical form of the org name
  // matches?
}

func OrgGroups( w http.ResponseWriter, r *http.Request)  {
  fmt.Fprintln(w, "Organization Groups ")
}

func OrgGroupShow( w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  org := vars["org"]
  group := vars["group"]
  fmt.Fprintln(w, "Organization " + org + " Group " + group)
}
