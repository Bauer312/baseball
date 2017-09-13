/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package util

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

/*
GetDBConnection provides a common database connection source for all of the tools
*/
func GetDBConnection() (*sql.DB, error) {
	dbuser, ok := os.LookupEnv("BASEBALL_DB_USER")
	if ok == false {
		return nil, errors.New("BASEBALL_DB_USER environment variable is not set")
	}
	dbpass, ok := os.LookupEnv("BASEBALL_DB_PASS")
	if ok == false {
		return nil, errors.New("BASEBALL_DB_PASS environment variable is not set")
	}
	dbname, ok := os.LookupEnv("BASEBALL_DB_NAME")
	if ok == false {
		return nil, errors.New("BASEBALL_DB_NAME environment variable is not set")
	}
	dbhost, ok := os.LookupEnv("BASEBALL_DB_HOST")
	if ok == false {
		dbhost = "localhost"
	}

	connectionString := fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' sslmode='%s'",
		dbuser, dbpass, dbname, dbhost, "disable")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to ping the database")
		return nil, err
	}
	return db, nil
}
