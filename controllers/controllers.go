package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var wadu []models.Job

var db *sql.DB

func init() {
	db = database.GetDB()

	layout := "2006-01-02T15:04"
	startTime, _ := time.Parse(layout, "2021-11-24T09:07")
	endTime, _ := time.Parse(layout, "2021-11-24T09:19")
	modified, _ := time.Parse(layout, "2021-11-24T09:38")

	checkers := []models.Job{
		{
			Created:   startTime,
			EndTime:   endTime,
			ID:        1,
			JobType:   "IT",
			Modified:  modified,
			Name:      "SDE2",
			StartTime: startTime,
			State:     models.JOB_STATE_RUNNING,
		},
		{
			Created:   startTime,
			EndTime:   endTime,
			ID:        2,
			JobType:   "IT",
			Modified:  modified,
			Name:      "SDE3",
			StartTime: startTime,
			State:     models.JOB_STATE_RUNNING,
		},
		{
			Created:   startTime,
			EndTime:   endTime,
			ID:        3,
			JobType:   "IT",
			Modified:  modified,
			Name:      "SDE3",
			StartTime: startTime,
			State:     models.JOB_STATE_SUCCESS,
		},
	}

	wadu = append(wadu, checkers...)
}

func getOneJob(id string) models.Job {
	var job models.Job

	query := "SELECT id, name, created, modified,state,data FROM jobs WHERE id = $1"
	var data []byte
	row := db.QueryRow(query, id)
	err := row.Scan(&job.ID, &job.Name, &job.Created, &job.Modified, &job.State, &data)

	if err != nil {
		// Handle the error according to your needs
		fmt.Printf("Error fetching job: %s\n", err)
		return models.Job{} // Return an empty Job struct or handle the error accordingly
	}
	if data != nil {
		// Initialize job.Data if it's nil
		if job.Data == nil {
			job.Data = &models.Response{}
		}

		// Unmarshal the data into the job's Data field
		err := json.Unmarshal(data, job.Data)
		if err != nil {
			fmt.Printf("Error decoding JSON data: %s\n", err)
			// Handle the error as needed
		}
	}
	return job
	//

}

func final_ip_job(check models.Job) models.Job {
	check.Created = time.Now()
	check.JobType = models.JOB_TYPE_IP_MEASUREMENT
	check.State = models.JOB_STATE_CREATED

	wadu = append(wadu, check)
	query := `
	INSERT INTO jobs (Created, EndTime, JobType, Modified, Name, StartTime, State)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (Name) DO NOTHING
	RETURNING id
`
	var (
		newID int
		err   error
	)

	err = db.QueryRow(query, check.Created, check.EndTime, check.JobType, check.Modified, check.Name, check.StartTime, check.State).Scan(&newID)
	if err != nil {
		fmt.Println("we here")
		return models.Job{}
	}
	if newID != 0 {
		check.ID = newID
		fmt.Println(newID)
		wadu = append(wadu, check)

		_, err = db.Exec("UPDATE jobs SET State=$1 WHERE ID=$2", models.JOB_STATE_SUCCESS, check.ID)
		if err != nil {
			log.Panic(err)
			fmt.Println("we here")
		}
		fmt.Println("from the final_ip")
		fmt.Println(check)
		return check

	}
	return models.Job{}

}
func do_create_job(check models.Job) models.Job {
	check.Created = time.Now()
	check.JobType = models.JOB_TYPE_IP_MEASUREMENT

	wadu = append(wadu, check)
	query := `
	INSERT INTO jobs (Created, EndTime, JobType, Modified, Name, StartTime, State)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (Name) DO NOTHING
	RETURNING id
`
	var (
		newID int
		err   error
	)

	err = db.QueryRow(query, check.Created, check.EndTime, check.JobType, check.Modified, check.Name, check.StartTime, check.State).Scan(&newID)
	if err != nil {
		fmt.Println(err)
		return models.Job{}
	}
	if newID != 0 {
		check.ID = newID
		fmt.Println(newID)
		wadu = append(wadu, check)
		return check

	}
	return models.Job{}

}

// all api hits
func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	fmt.Println("we here")

	if wadu == nil {

		check := models.Result{
			Data:    []*models.Job{},
			Error:   "",
			Success: true,
		}
		json.NewEncoder(w).Encode(check)
		return
	}
	json.NewEncoder(w).Encode(wadu)

}

func Do_create_job(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	if r.Body == nil {
		json.NewEncoder(w).Encode("no data ")

	}
	var check models.Job

	json.NewDecoder(r.Body).Decode(&check)
	check = do_create_job(check)

	json.NewEncoder(w).Encode(check)

}
func Do_get_job(w http.ResponseWriter, r *http.Request) {
	fmt.Println("we here")
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	checkers := getOneJob(params["id"])
	if checkers.ID == 0 {

		check := models.Result{
			Data:    []*models.Job{},
			Error:   "Error=Not Found",
			Success: false,
		}
		json.NewEncoder(w).Encode(check)
		return
	}
	check := models.Result{
		Data:    []*models.Job{&checkers},
		Error:   "",
		Success: false,
	}
	json.NewEncoder(w).Encode(check)

}

func Handle_test_ip(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	checkersChan := make(chan models.Response, 1)
	newIDChan := make(chan int, 1)
	errorChan := make(chan error, 1)
	wg.Add(3)
	//2 processes create job and get the ip from the user
	if r.Body == nil {
		json.NewEncoder(w).Encode("no data ")
	}

	//one goroute to create the jobs field in the db

	//one to check the ip address provided
	var result models.Request

	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		errorChan <- err
		return

	}

	fmt.Println(result.Ipaddress)
	var checkers models.Response
	var beaters models.Job
	go func() {

		defer wg.Done()

		checkers = getDetailsofthePing(result.Ipaddress, result.Count)
		responseJSON, err := json.Marshal(checkers)
		if err != nil {
			log.Fatalf("Error marshaling response to JSON: %s", err)
		}

		filr, err := os.OpenFile("successlog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("errpr")
		}
		filr.Write(responseJSON)
		fmt.Println("Sending to checkersChan")
		checkersChan <- checkers
		fmt.Println("Sent to checkersChan")
		fmt.Println(checkersChan)
		fmt.Println("no deadlocks 1")

	}()
	go func() {

		defer wg.Done()
		var check models.Job
		check.Name = result.Ipaddress
		beaters = final_ip_job(check)
		fmt.Println(beaters.ID)

		newIDChan <- beaters.ID

		var res models.HandleTestIPResponse
		res.JobId = beaters.ID
		json.NewEncoder(w).Encode(res)
		fmt.Println("no deadlocks 2")
	}()

	go func() {
		defer wg.Done()

		newID := <-newIDChan
		fmt.Printf("we are printint newId %d", newID)
		checkers := <-checkersChan
		fmt.Println(checkers)

		jsonData, err := json.Marshal(checkers)
		if err != nil {
			errorChan <- err
			return
		}

		_, err = db.Exec("UPDATE jobs SET Data = $1 WHERE id = $2", jsonData, newID)
		if err != nil {
			errorChan <- err
			fmt.Println("Error in db.Exec:", err)
			return
		}
		fmt.Println("no deadlocks 3")
	}()

	wg.Wait()
	close(errorChan)
	if err, ok := <-errorChan; ok {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	fmt.Println("we are here really")
	// chech, _ := json.Marshal(newIDChan)
	// json.NewEncoder(w).Encode(chech)
	// var res models.HandleTestIPResponse
	// res.JobId = <-newIDChan
	// fmt.Println(res)
	// fmt.Println("the top one")
	// json.NewEncoder(w).Encode(res)
	// return

}
