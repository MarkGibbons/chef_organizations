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
        "chef_organizations/co"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	// "github.com/go-chef/chef"
	chef "github.com/MarkGibbons/chefapi"
	_ "github.com/go-sql-driver/mysql"
)

// organizationsSchema defines the organization database and its tables.
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
        var dbc co.DbConnectionRequest
	dbc.PwdFile  = os.Args[4]
	dbc.Server = "127.0.0.1"
	dbc.Port = "3306"
	dbc.User = "root"
	dbc.Database = "organizations"

	// TODO: Listen for update requests
	// TODO: Delete organizations that have been deleted
	// Create the database and add the schema
	dbInit(dbc)

	// Extract and Update on a timer

	// Build an api client instance.
	client := buildClient(user, keyfile, chefurl)
	// Execute the get from chef server, update mysql cycle.
	for {
		// Open database connection
		db := co.DbConnection(dbc)
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

// buildClient creates a connection to a chef server using the chef api.
func buildClient(user string, keyfile string, baseurl string) *chef.Client {
	key := clientKey(keyfile)
	client, err := chef.NewClient(&chef.Config{
		Name:    user,
		Key:     string(key),
		BaseURL: baseurl,
		// goiardi is on port 4545 by default, chef-zero is 8889, chef-server is on 443
	})
	if err != nil {
		fmt.Println("Issue setting up client:", err)
		os.Exit(1)
	}
	return client
}

// clientKey reads the pem file containing the credentials needed to use the chef client.
func clientKey(filepath string) string {
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Couldn't read key.pem:", err)
		os.Exit(1)
	}
	return string(key)
}

// dbInit creates the organization data base and tables.
func dbInit(dbc co.DbConnectionRequest) {
	db, err := sql.Open("mysql", dbc.User+":"+dbc.Pwd+"@tcp("+dbc.Server+":"+dbc.Port+")/")
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

// org2DB adds organizations to the organizations database in the organizations table.
func org2DB(db *sql.DB, org string) {
	// See if org is already there
	stmtOrgName := db.Prepare("SELECT name FROM organizations  WHERE name = ?)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	checkOrg := stmtOrgName.Exec(org)
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

// groupsOrg2DB adds the groups and members to the database.
// The Org/group/member connections are updated by replacing all existing rows that contain
// this org/group being processeda. A transaction is used so that requesters see the 
// group members update as an atomic action.
func groupsOrg2DB(client *chef.Client, org string, db *sql.DB) {
	//         Get the list of groups in the organization
	groupList := orgGroups(client, org)
	stmtDelGrp, err := db.Prepare("DELETE FROM org_groups WHERE organization_name = ? AND group_name = ?;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for group := range groupList {
		// skip chefs internal groups
		if co.IsUSAG(group) || group == "clients" {
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

// orgGroups gets a list of groups, from the chef server, belonging to an organization.
func orgGroups(client *chef.Client, org string) map[string]string {
	groupList, err := client.Groups.List()
	if err != nil {
		fmt.Println("Issue listing groups:", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return groupList
}

// getGroup gets group information from the chef server. The
// members of the group and nested groups are retrieved.
func getGroup(client *chef.Client, group string) chef.Group {
	groupInfo, err := client.Groups.Get(group)
	if err != nil {
		fmt.Println("Issue getting: "+group, err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return groupInfo
}

// getMember gets the information associated with a particular user account.
func getMember(client *chef.Client, member string) chef.User {
	memberInfo, err := client.Users.Get(member)
	if err != nil {
		fmt.Println("Issue getting: "+member, err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return memberInfo
}

// getGroupMember gets all the members in a group form the chef server.
// Chef groups have three lists of members.  There is a list of 
// actors, a list of users and a list of nested groups
func getGroupMembers(client *chef.Client, groupInfo chef.Group) []string {
        // TODO: Verify we are extracting the correct set of users
	members := usersFromGroups(client, groupInfo.Groups)
	members = append(members, groupInfo.Actors...)
	members = append(members, groupInfo.Users...)
	members = co.Unique(members)
	return members
}

// usersFromGroups gets the nested groups. getGroupMembers and userFromGroups
// call each other in a recursive fashion to expand the nested groups
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

// groupMembers2DB updates the members table.
// It adds and/or updates a member entry.
// It adds a org_groups row for each member
func groupMembers2DB(groupMembers []string, org string, group string, db *sql.DB) {
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

//memberUpdate updates the member information in the database.
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
	members = co.Unique(members)
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
