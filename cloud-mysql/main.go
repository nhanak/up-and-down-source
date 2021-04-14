package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// InsertRequest parsed json struct for making an insert into cloud mysql db
type InsertRequest struct {
	WinnerName          string
	WinnerCredits       string
	WinnerHealth        string
	LoserName           string
	LoserCredits        string
	LoserHealth         string
	GameLengthInSeconds string
}

func parseInsertRequest(r *http.Request) (*InsertRequest, bool) {
	var ir InsertRequest
	err := json.NewDecoder(r.Body).Decode(&ir)
	if err != nil {
		log.Printf("Error decoding body of insert request: %v", err)
		return nil, false
	}
	return &ir, true
}

func main() {
	http.HandleFunc("/insert", insertHandler)
	port := "666"
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("InsertHandler Pinged!")
	publishGameOverToMySQL(r)
}

//https://github.com/GoogleCloudPlatform/cloudsql-proxy
func publishGameOverToMySQL(r *http.Request) {
	ir, ok := parseInsertRequest(r)
	if !ok {
		log.Printf("Something went wrong with parsing insert request")
		return
	}

	log.Printf("DIALING SQL")
	db, err := sql.Open("mysql", "neil@tcp(127.0.0.1:3306)/up_and_down")
	if err != nil {
		log.Printf("Something went wrong opening db connection: %v", err)
		return
	}
	defer db.Close()

	q := `INSERT INTO unranked (winnerName, winnerCredits, winnerHealth, loserName, loserCredits, loserHealth, gameLengthInSeconds) VALUES ( '` + ir.WinnerName + `', ` + ir.WinnerCredits + `, ` + ir.WinnerHealth + `, '` + ir.LoserName + `', ` + ir.LoserCredits + `, ` + ir.LoserHealth + `, ` + ir.GameLengthInSeconds + ` );`
	log.Printf("Query is: %v", q)
	insert, err := db.Query("INSERT INTO unranked (winnerName, winnerCredits, winnerHealth, loserName, loserCredits, loserHealth, gameLengthInSeconds) VALUES (?,?,?,?,?,?,?);", ir.WinnerName, ir.WinnerCredits, ir.WinnerHealth, ir.LoserName, ir.LoserCredits, ir.LoserHealth, ir.GameLengthInSeconds)
	if err != nil {
		log.Printf("Something went wrong opening db connection: %v", err)
		return
	}
	defer insert.Close()

}
