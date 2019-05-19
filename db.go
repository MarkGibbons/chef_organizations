// Chef Organization
// Data base common functions
package co

import (
	"database/sql"
        "io/ioutil"
        "strings"
	_ "github.com/go-sql-driver/mysql"
)

type DbConnectionRequest struct {
  Server string
  User string
  PwdFile string
  Pwd string
  Port string
  Database string
}

func DbConnection(dbc DbConnectionRequest) *sql.DB {
        dbc.Pwd = DbPWD(dbc.PwdFile)
        db, err := sql.Open("mysql", dbc.User+":"+dbc.Pwd+"@tcp("+dbc.Server+":"+dbc.Port+")/"+dbc.Database)
        if err != nil {
                panic(err.Error()) // proper error handling instead of panic in your app
        }
        return db
}

// DbPWD reads and returns a password from a specified file path.
func DbPWD(dbpwdFile string) string {
	pwd, err := ioutil.ReadFile(dbpwdFile)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return strings.TrimSpace(string(pwd))
}
