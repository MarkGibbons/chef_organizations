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
	// "log"
	"os"
	"time"

	"github.com/go-chef/chef"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	user := os.Args[1]
	keyfile := os.Args[2]
	chefurl := os.Args[3]
	dbName := "127.0.0.1"
	dbPort := "3306"

	// TODO: Listen for update requests
	// TODO: Delete organizations that have been deleted

	// Extract and Update on a timer
	// Build a client
	client := buildClient(user, keyfile, chefurl)
	// Update cycle
	for {
		// Open database connection
		db := dbConnection(dbName, dbPort)
		// Get list of organizations
		orgList := listOrganizations(client)
		// For each organization
                for org, data := range orgList {
                        fmt.Println("ORGLIST DATA")
                        fmt.Println(data)
                        orgclient := buildClient(user, keyfile, chefurl + "/" + org + "/")
		        // Add organization if not there
                        // TODO: pass in the org data here
			org2DB(db, org)
		        // Get the list of groups, update db
                        groupsOrg2DB(orgclient, org, db)
                }
		//      Close the data base connection
                db.Close()
                time.Sleep(180 * time.Second)
	}

	// List admin Group
	groupList, err := client.Groups.Get("admins")
	if err != nil {
		fmt.Println("Issue listing admins:", err)
	}
	// Print out the list
	fmt.Println(groupList)

	// List Environments
	envInfo, err := client.Environments.List()
	if err != nil {
		fmt.Println("Issue listing environments:", err)
	}
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

func dbConnection(dbname string, dbport string) *sql.DB {
	db, err := sql.Open("mysql", "root@tcp("+dbname+":"+dbport+")/organizations")
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

// List organizations
func listOrganizations(client *chef.Client) map[string] string {
	orgList, err := client.Organizations.List()
	if err != nil {
		fmt.Println("Issue listing orgs:", err)
	}
	return orgList
}

func org2DB(db *sql.DB, org string) {
        // See if org is already there
        row := db.QueryRow("SELECT name FROM organizations  WHERE name = '" + org + "';")
        // TODO? Close row query
        fmt.Println("QUERY RESULTS")
        fmt.Println(row)
        var name string
        switch err := row.Scan(&name); err {
        case sql.ErrNoRows:
        	fmt.Println("Add this org " + org)
        case nil:
        	fmt.Println("Scanned Name" + name)
        	fmt.Println("Present org " + org)
          	return
        default:
                panic(err.Error()) // proper error handling instead of panic in your app
        }

        // Prepare statement for inserting organizations
        // TODO: Need to do organization.Get to the chef server to get the full_name. Add later
        // TODO need the org information
        stmtInsOrg, err := db.Prepare("INSERT INTO organizations (name) VALUES( ? )") // ? = placeholder
        if err != nil {
                panic(err.Error()) // proper error handling instead of panic in your app
        }
        _, err = stmtInsOrg.Exec(org)
        if err != nil {
               panic(err.Error()) // proper error handling instead of panic in your app
        }
	return
}

func groupsOrg2DB(client *chef.Client, org string, db *sql.DB) {
	//         Get the list of groups in the organization
        orgGroups(client, org)
	//           For each group`
	//             Get the members
	//             Begin
	//               Delete all rows in groups that match this group
	//               Add and or update the member entry
	//               For each member
	//                 Add a group entry with the member
	//             Commit
}

func orgGroups(client *chef.Client, org string) map[string]string {
	// List Groups
	groupInfo, err := client.Groups.List()
	if err != nil {
		fmt.Println("Issue listing groups:", err)
	}
	// Print out the list
	// fmt.Println(groupInfo)
        return groupInfo
}
