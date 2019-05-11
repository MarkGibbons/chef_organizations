//
// Get chef organization information from the chef server via the api
// Update the organization data base with the extracted information
//
// https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html
// http://www.mysqltutorial.org/mysql-create-table/
//
// https://gowalker.org/github.com/go-chef/chef  and https://github.com/go-chef/chef  api
//
package main

import (
	"fmt"
	"io/ioutil"
	// "log"
	"os"

	// "github.com/gorilla/mux"
	"github.com/go-chef/chef"
)

func main() {
	user := os.Args[1]
	keyfile := os.Args[2]
	chefurl := os.Args[3]

	// Listen for update requests
	// Extract and Update on a timer

	// Get client key
	key := clientKey(keyfile)

	// Build a client
	client := buildClient(user, key, chefurl)
        fmt.Println(client)

	// List organizations
	orgList, err := client.Organizations.List()
	if err != nil {
		fmt.Println("Issue listing orgs:", err)
	}
	// Print out the list
	fmt.Println(orgList)

	// List uis organization
	uisShow, err := client.Organizations.Get("uis")
	if err != nil {
		fmt.Println("Issue getting uis org:", err)
	}
	// Print out the list
	fmt.Println(uisShow)

	// List admin Group
	groupList, err := client.Groups.Get("admins")
	if err != nil {
		fmt.Println("Issue listing admins:", err)
	}
	// Print out the list
	fmt.Println(groupList)

	// List Groups
	groupInfo, err := client.Groups.List()
	if err != nil {
		fmt.Println("Issue listing groups:", err)
	}
	// Print out the list
	fmt.Println(groupInfo)

	// List Environments
	envInfo, err := client.Environments.List()
	if err != nil {
		fmt.Println("Issue listing environments:", err)
	}
	// Print out the list
	fmt.Println(envInfo)

	// List Cookbooks - works
//	cookList, err := client.Cookbooks.List()
	//if err != nil {
		//fmt.Println("Issue listing cookbooks:", err)
	//}
	// Print out the list
	//fmt.Println(cookList)

	// Extract the organizations
	// Extract the groups from each organization
        // organization := "uis"
        // groups := listGroups(&client, organization)
        // fmt.Println(groups)

	// Extract the group members

	// router := mux.NewRouter().StrictSlash(true)
	// log.Fatal(http.ListenAndServe(":8080", router))
}

func clientKey(filepath string) string {
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Couldn't read key.pem:", err)
		os.Exit(1)
	}
	return string(key)
}

// func listGroups(client string, key string, baseurl string) *chef.Client {
// 	client, err := client.Groups.List() {
// 	})
// 	if err != nil {
// 		fmt.Println("Issue setting up client:", err)
//		os.Exit(1)
//	}
//	return client
//}

func buildClient(user string, key string, baseurl string) *chef.Client {
	client, err := chef.NewClient(&chef.Config{
		Name:    user,
		Key:     string(key),
		BaseURL: baseurl,
		// goiardi is on port 4545 by default. chef-zero is 8889. chef-server is on 443
	})
	if err != nil {
		fmt.Println("Issue setting up client:", err)
		os.Exit(1)
	}
	return client
}
