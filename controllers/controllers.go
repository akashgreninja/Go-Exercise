package controllers

import (
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var wadu []models.Job

func init() {

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
	b, _ := strconv.Atoi(id)
	// var emptyJob models.Job

	for _, i := range wadu {
		if i.ID == b {
			return i
		}
	}
	// return emptyJob
	return models.Job{}

}

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

	err := json.NewDecoder(r.Body).Decode(&check)
	check.Created = time.Now()
	check.JobType = models.JOB_TYPE_IP_MEASUREMENT

	check.ID = 1

	if err != nil {
		json.NewEncoder(w).Encode("num-failed")

	}
	wadu = append(wadu, check)
	json.NewEncoder(w).Encode(wadu)

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
		Error:   "Error=Not Found",
		Success: false,
	}
	json.NewEncoder(w).Encode(check)

}
