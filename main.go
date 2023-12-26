package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"log/syslog"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/gorilla/mux"
)

var db *sql.DB
var logger *syslog.Writer

func init() {
	var err error
	// Connect to syslog
	logger, err = syslog.New(syslog.LOG_INFO, "hbapp")
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Connected to syslog")

	// Connect to your database
	db, err = sql.Open("mysql", "kamailio:kamailio@tcp(localhost:3306)/kamailio")
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Database Connected")
	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Received Pong from Database")
}

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain string `json:"domain"`
}

func main() {
	r := mux.NewRouter()

	// Define your API endpoint
	r.HandleFunc("/check/{input}", CheckExistence).Methods("GET")
	r.HandleFunc("/isuseronline/{input}", IsUserOnline).Methods("GET")
	r.HandleFunc("/addSubs", AddSubscriber).Methods("POST")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":9090", r))
	logger.Info("Server Started listening on port 9090")
}

// IsUserOnline is the handler for the API endpoint
func IsUserOnline(w http.ResponseWriter, r *http.Request) {
	// Get the input parameter from the URL
	vars := mux.Vars(r)
	input := vars["input"]

	// Call a function to check if the input exists in the database
	exists, err := checkDatabase(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result as JSON
	if exists {
		fmt.Fprint(w, `{"isUserOnline": true}`)
	} else {
		fmt.Fprint(w, `{"isUserOnline": false}`)
	}

}

// Add Subscriber
func AddSubscriber(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	logger.Info("Request Received to add Subscriber")
	fmt.Printf("%v\n", r.Body)
	var requestBody RequestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		logger.Err("Invalid JSON")
		return
	}

	fmt.Printf("%v\n", requestBody)

	// Call a function to insert data into the database
	err = insertIntoDatabase(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	fmt.Fprint(w, "Data inserted successfully")
}

// insertIntoDatabase inserts data into the MySQL database
func insertIntoDatabase(data RequestBody) error {
	// Perform a SQL INSERT query
	// Replace "your_table" and the column names with your actual table and column names
	query := "INSERT INTO subscriber (username, domain, password) VALUES (?, ?, ?)"
	_, err := db.Exec(query, data.Username, data.Domain, data.Password)
	if err != nil {
		return err
	}

	return nil
}


// CheckExistence is the handler for the API endpoint
func CheckExistence(w http.ResponseWriter, r *http.Request) {
	// Get the input parameter from the URL
	vars := mux.Vars(r)
	input := vars["input"]

	// Call a function to check if the input exists in the database
	exists, err := checkDatabase(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result as JSON
	if exists {
		fmt.Fprint(w, `{"exists": true}`)
	} else {
		fmt.Fprint(w, `{"exists": false}`)
	}
}

// checkDatabase is a function to check if the input exists in the database
func checkDatabase(input string) (bool, error) {
	// Perform a query to check existence in your database
	// Replace "your_table" and "your_column" with your actual table and column names
	query := "SELECT COUNT(*) FROM location WHERE username = ?"
	var count int
	err := db.QueryRow(query, input).Scan(&count)
	if err != nil {
		return false, err
	}

	// If count is greater than 0, the input exists in the database
	return count > 0, nil
}

