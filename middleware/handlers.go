package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"

	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route
	"github.com/mateuscdiniz/simplecrud/models"

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateJob create a job in the postgres db
func CreateJob(w http.ResponseWriter, r *http.Request) {
	// create an empty job of type models.Job
	var job models.Job

	// decode the json request to job
	err := json.NewDecoder(r.Body).Decode(&job)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// call insert job function and pass the job
	insertID := insertJob(job)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Job created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetJob will return a single Job by its id
func GetJob(w http.ResponseWriter, r *http.Request) {

	// get the jobid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// call the getJob function with job id to retrieve a single job
	job, err := getJob(int64(id))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// send the response
	json.NewEncoder(w).Encode(job)
}

// GetAllJob will return all the jobs
func GetAllJobs(w http.ResponseWriter, r *http.Request) {

	// get all the jobs in the db
	jobs, err := getAllJobs()

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// send all the jobs as response
	json.NewEncoder(w).Encode(jobs)
}

// UpdateJob update job's detail in the postgres db
func UpdateJob(w http.ResponseWriter, r *http.Request) {

	// get the jobid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// create an empty job of type models.Job
	var job models.Job

	// decode the json request to job
	err = json.NewDecoder(r.Body).Decode(&job)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// call update job to update the job
	updatedRows := updateJob(int64(id), job)

	// format the message string
	msg := fmt.Sprintf("Job updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteJob delete job's detail in the postgres db
func DeleteJob(w http.ResponseWriter, r *http.Request) {

	// get the jobid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// call the deleteJob, convert the int to int64
	deletedRows := deleteJob(int64(id))

	// format the message string
	msg := fmt.Sprintf("Job updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one job in the DB
func insertJob(job models.Job) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning jobid will return the id of the inserted job
	sqlStatement := `INSERT INTO jobs (name) VALUES ($1) RETURNING jobid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, job.Name).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one job from the DB by its jobid
func getJob(id int64) (models.Job, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a job of models.Job type
	var job models.Job

	// create the select sql query
	sqlStatement := `SELECT * FROM jobs WHERE jobid=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to job
	err := row.Scan(&job.ID, &job.Name)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return job, nil
	case nil:
		return job, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty job on error
	return job, err
}

// get one job from the DB by its jobid
func getAllJobs() ([]models.Job, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var jobs []models.Job

	// create the select sql query
	sqlStatement := `SELECT * FROM jobs`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var job models.Job

		// unmarshal the row object to job
		err = rows.Scan(&job.ID, &job.Name)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the job in the jobs slice
		jobs = append(jobs, job)

	}

	// return empty job on error
	return jobs, err
}

// update job in the DB
func updateJob(id int64, job models.Job) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE jobs SET name=$2, location=$3, age=$4 WHERE jobid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, job.Name)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete job in the DB
func deleteJob(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM jobs WHERE jobid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
