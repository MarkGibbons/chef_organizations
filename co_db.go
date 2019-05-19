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
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-chef/chef"
	_ "github.com/go-sql-driver/mysql"
)

var organizationsSchema = []string{

	"CREATE DATABASE IF NOT EXISTS organizations;",

	"USE organizations;",

	"CREATE TABLE IF NOT EXISTS organizations ( name TEXT, full_name  TEXT) ENGINE=INNODB;",

	"CREATE TABLE IF NOT EXISTS org_groups ( group_name TEXT, organization_name  TEXT, user_name  TEXT) ENGINE=INNODB;",

	"CREATE TABLE IF NOT EXISTS members ( user_name TEXT, email  TEXT, display_name  TEXT) ENGINE=INNODB;",
}

// main uses a loop to get information from the chef server and update a mysql data base.
// The chef server is accessed via the go-chef/chef server package and uses the chef-server api.
// The mysql data base is used as a cache for a REST interface. See co_rest.go for the REST details.
func main() {
	// Pass in the database and chef-server api credentials.
	user := os.Args[1]
	keyfile := os.Args[2]
	chefurl := os.Args[3]
	dbPWDFile  := os.Args[4]
	dbName := "127.0.0.1"
	dbPort := "3306"
	dbUser := "root"

	// TODO: Listen for update requests
	// TODO: Delete organizations that have been deleted
	// Create the database and add the schema
	dbInit(dbName, dbPort, dbUser, dbPWD(dbPWDFile))

	// Extract and Update on a timer

	// Build an api client instance.
	client := buildClient(user, keyfile, chefurl)
	// Execute the get from chef server, update mysql cycle.
	for {
		// Open database connection
		db := dbConnection(dbName, dbPort, dbUser, dbPWD(dbPWDFile))
		// Get list of organizations
		orgList := listOrganizations(client)
		// For each organization
		for org := range orgList {
			orgclient := buildClient(user, keyfile, chefurl+"/"+org+"/")
			// Add organization if not there
			org2DB(db, org)
			// Get the list of groups, update db
			groupsOrg2DB(orgclient, org, db)
		}
		memberUpdate(db, client)
		// Close the data base connection
		db.Close()
		time.Sleep(180 * time.Second)
	}
}

func buildClient(user string, keyfile string, baseurl string) *chef.Client {
	key := clientKey(keyfile)
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

func clientKey(filepath string) string {
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Couldn't read key.pem:", err)
		os.Exit(1)
	}
	return string(key)
}

func dbConnection(dbname string, dbport string, dbuser string, dbpwd string) *sql.DB {
	db, err := sql.Open("mysql", dbuser+":"+dbpwd+"@tcp("+dbname+":"+dbport+")/organizations")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}

// dbInit creates the organization data base and tables.
func dbInit(dbname string, dbport string, dbuser string, dbpwd string) {
	fmt.Println(dbuser+":"+dbpwd+"@tcp("+dbname+":"+dbport+")")
	db, err := sql.Open("mysql", dbuser+":"+dbpwd+"@tcp("+dbname+":"+dbport+")/")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for _, stmt := range organizationsSchema {
		fmt.Println(stmt)
		_, err := db.Exec(stmt)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	db.Close()
	return
}

// listOrganizations uses the chef server api to list all organizations
func listOrganizations(client *chef.Client) map[string]string {
	orgList, err := client.Organizations.List()
	if err != nil {
		fmt.Println("Issue listing orgs:", err)
	}
	return orgList
}

func org2DB(db *sql.DB, org string) {
	// See if org is already there
	checkOrg := db.QueryRow("SELECT name FROM organizations  WHERE name = '" + org + "';")
	var name string
	switch err := checkOrg.Scan(&name); err {
	case sql.ErrNoRows:
	case nil:
		return
	default:
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Prepare statement for inserting organizations
	stmtInsOrg, err := db.Prepare("INSERT INTO organizations (name) VALUES( ? )")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	_, err = stmtInsOrg.Exec(org)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
        stmtInsOrg.Close()
	return
}

func groupsOrg2DB(client *chef.Client, org string, db *sql.DB) {
	//         Get the list of groups in the organization
	groupList := orgGroups(client, org)
	stmtDelGrp, err := db.Prepare("DELETE FROM org_groups WHERE organization_name = ? AND group_name = ?;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for group := range groupList {
		// skip chefs internal groups
		if isUsag(group) || group == "clients" {
			continue
		}

		groupInfo := getGroup(client, group)

		// Update the group entries in a transaction
		tx, err := db.Begin()

		// Delete all rows in groups that match this group
		_, err = stmtDelGrp.Exec(org, group)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// Consolidate the member and actor arrays for the group
		groupMembers := getGroupMembers(client, groupInfo)
		// Add and/or update the member entry unless it exists
		groupMembers2DB(groupMembers, org, group, db)

		tx.Commit()
	}
	stmtDelGrp.Close()
}

func orgGroups(client *chef.Client, org string) map[string]string {
	groupList, err := client.Groups.List()
	if err != nil {
		fmt.Println("Issue listing groups:", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return groupList
}

func getGroup(client *chef.Client, group string) chef.Group {
	groupInfo, err := client.Groups.Get(group)
	if err != nil {
		fmt.Println("Issue getting: "+group, err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return groupInfo
}

func getMember(client *chef.Client, member string) chef.User {
	memberInfo, err := client.Users.Get(member)
	if err != nil {
		fmt.Println("Issue getting: "+member, err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return memberInfo
}

func isUsag(group string) bool {
	match, err := regexp.MatchString("^[0-9a-f]+$", group)
	if err != nil {
		fmt.Println("Issue with regex", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return len(group) == 32 && match
}

// getGroupMember gets all the members in a group.
// Chef groups have three lists of members.  There is a list of 
// actors, a list of users and a list of nested groups
func getGroupMembers(client *chef.Client, groupInfo chef.Group) []string {
        // TODO: Verify we are extracting the correct set of users
	members := usersFromGroups(client, groupInfo.Groups)
	members = append(members, groupInfo.Actors...)
	members = append(members, groupInfo.Users...)
	members = unique(members)
	return members
}

func usersFromGroups(client *chef.Client, groups []string) []string {
	var members []string
	for _, group := range groups {
		groupInfo, err := client.Groups.Get(group)
		if err != nil {
			fmt.Println("Issue with regex", err)
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		members = getGroupMembers(client, groupInfo)
	}
	return members
}

func groupMembers2DB(groupMembers []string, org string, group string, db *sql.DB) {
	// Add and/or update the member entry unless it exists
	// Add a org_groups row for each member
	stmtInsOrgGroup, err := db.Prepare("INSERT INTO org_groups (group_name, organization_name, user_name) VALUES( ?, ?, ? )")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for _, member := range groupMembers {
		_, err = stmtInsOrgGroup.Exec(group, org, member)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println("Add group org member ", group, org, member)
	}
        stmtInsOrgGroup.Close()
}

func memberUpdate(db *sql.DB, client *chef.Client) {
	// Get a unique list of all the users
	users, err := db.Query("SELECT user_name FROM org_groups;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var members []string
	for users.Next() {
		var name string
		err = users.Scan(&name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		members = append(members, name)
	}
	members = unique(members)
	users.Close()
	stmtInsMember, err := db.Prepare("INSERT INTO members (user_name, email, display_name) VALUES( ?, ?, ? )")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
        // transaction - begin, delete existing users, add replacement data
	tx, err := db.Begin()
	_, err = db.Query("DELETE FROM members;")
	for _, member := range members {
		// Extract information for each user
		memberInfo := getMember(client, member)
		// Update the data base with a new set of user records
		_, err = stmtInsMember.Exec(member, memberInfo.Email, memberInfo.DisplayName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
        stmtInsMember.Close()
	tx.Commit()
}

// unique takes and array and return the unique elements of the array
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

// dbPWD read a password from a specified file path
func dbPWD(dbpwdFile string) string {
	pwd, err := ioutil.ReadFile(dbpwdFile)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return strings.TrimSpace(string(pwd))
}
