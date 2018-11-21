package main // import "github.com/ypeckstadt/couchdb-cleaner"

import (
	"encoding/json"
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/pkg/errors"
	"github.com/ypeckstadt/couchdb-cleaner/environment"
	"github.com/ypeckstadt/couchdb-cleaner/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	// Load env variables
	var environment environment.Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	checkPanic(err)


	log.Printf("CouchDB cleaner started for %s", environment.URLs)


	// get the URLs
	couchDBURLs := strings.Split(environment.URLs, ",")

	if len(couchDBURLs[0]) == 0 {
		panic(errors.New("no couchDB urls where provided"))
	}

	ticker := time.NewTicker(time.Duration(environment.CleanInterval) * time.Millisecond)
	for {
		select {
			case <- ticker.C:
				go func() {
					log.Println("Starting cleanup and compaction ...")
					for _, url := range couchDBURLs {

						// Get all databases
						databases := getAllDatabases(url)

						// loop through the databases
						if len(databases[0]) > 0 {
							for _, database := range databases {
								compactDatabase(url, database)
								databaseViews := getViewsForDatabase(url, database)

								for _, view := range databaseViews.Rows {
									viewID := strings.Replace(view.ID, "_design/", "",1)
									compactView(url, database, viewID)
								}
								cleanupViews(url, database)
							}
						}
					}
					log.Println("Cleanup and compaction have finished")
				}()
		}
	}
}

func cleanupViews(url string, database string) {
	response, err := http.Post(url + "/" + database + "/_view_cleanup", "application/json", nil)
	check(err)
	if response.Status == "202 Accepted" {
		log.Printf("Views have been cleaned up for database %s at %s", database, url)
	} else {
		log.Printf("Views have not been cleaned up for database %s at %s has failed", database, url)
	}
}

func compactView(url string, database string, viewID string) {
	response, err := http.Post(url + "/" + database + "/_compact/" + viewID, "application/json", nil)
	check(err)
	if response.Status == "202 Accepted" {
		log.Printf("View %s has been compacted for database %s at %s", viewID, database, url)
	} else {
		log.Printf("View %s has not been compacted for database %s at %s has failed", viewID, database, url)
	}
}

func getViewsForDatabase(url string, database string) model.Designs {
	var views model.Designs

	response, err := http.Get(url + "/" + database +"/_design_docs")
	check(err)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(contents, &views)
	check(err)

	return views
}

func compactDatabase(url string, database string) {
	response, err := http.Post(url + "/" + database + "/_compact", "application/json", nil)
	check(err)
	if response.Status == "202 Accepted" {
		log.Printf("Database %s has been compacted at %s", database, url)
	} else {
		log.Printf("Database %s compaction at %s has failed", database, url)
	}
}

func getAllDatabases(url string) []string {
	response, err := http.Get(url + "/_all_dbs")
	check(err)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	var databases []string
	err = json.Unmarshal(contents, &databases)
	check(err)

	fmt.Println(databases)

	return databases
}

func checkPanic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}