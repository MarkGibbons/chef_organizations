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
        _ "github.com/go-sql-driver/mysql"
)

func main() {
	// user := os.Args[1]
	// password := os.Args[2]
        // get the user and password more securely
        db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/organizations")
        if err != nil {
                 fmt.Println(err.Error())
        }
        defer db.Close()
        // Prepare statement for inserting organizations
        // name, full_name
	stmtInsOrg, err := db.Prepare("INSERT INTO organizations VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtInsOrg.Close() // Close the statement when we leave main() / the program terminates
        _, err = stmtInsOrg.Exec("uis", "Unix and Linux")
        if err != nil {
 	       panic(err.Error()) // proper error handling instead of panic in your app
	}


}
