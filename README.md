#  Overview
The chef organization application extracts information from the chef server that is not available to unauthorized users.
The application consists of 4 processes. A mysql database, a process written in go to extract information from the chef
server, a go process the supports a REST interface to access the database and an nginx server provides a UI.  The extract process
uses the go-chef/chef chef server api with additions for organization and user support.  The UI is written in Java Script.


# Links to help with development

## Extract info
https://docs.chef.io/api_chef_server.html
https://www.thepolyglotdeveloper.com/2017/07/consume-restful-api-endpoints-golang-application/
https://tutorialedge.net/golang/consuming-restful-api-with-go/
https://tutorialedge.net/golang/golang-mysql-tutorial/  good example of using the mysql api
https://stackoverflow.com/questions/16029441/how-to-delete-multiple-rows-in-sql-where-id-x-to-y
http://www.mysqltutorial.org/mysql-insert-statement.aspx
Will need the pivotal user credentials

## DB tables
http://go-database-sql.org/references.html

## Java script
// https://www.taniarascia.com/how-to-connect-to-an-api-with-javascript/  Seems like a good tutorial
// https://github.com/taniarascia/sandbox/tree/master/ghibli source code for the js sample application, index, css and js
Generally good blog stuff https://www.taniarascia.com/blog/


# Start and stop commands for testing

## Run co_db to gather group info from the uis org
go run -tags debug co_db.go xmjg ~/.chef/xmjg.pem https://chefp01.nordstrom.net/organizations/uis/

## Run co_db to update the db
sudo env 'GOPATH=/home/xmjg/go' go run co_db.go pivotal /u01/restaurant/tools/pivotal.pem https://chefp01.nordstrom.net/organizations /u01/restaurant/tools/menu/pivotal.pem

## Run co_rest to start the rest interface
sudo env 'GOPATH=/home/xmjg/go' go run co_rest.go 
