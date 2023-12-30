package models

import "time"

// ip models
type Request struct {
	Ipaddress string `json:"ipaddress,omitempty"`
	Count     int    `json:"count,omitempty"`
}

type Response struct {
	Ipaddress   string   `json:"ipaddress,omitempty"`
	Count       int      `json:"count,omitempty"`
	Transmitted int      `json:"transmitted,omitempty"`
	Received    int      `json:"received,omitempty"`
	Packetloss  string   `json:"packetloss,omitempty"`
	Latency     *Latency `json:"latency,omitempty"`
}

type Latency struct {
	Min  string
	Max  string
	Avg  string
	Mdev string
}
type Errorlog struct {
	Ipaddress string
	Count     int
}

//JOBS MODELS

type JobState string

const (
	JOB_STATE_CREATED   JobState = "ACTIVE"
	JOB_STATE_SCHEDULED JobState = "JOB_STATE_SCHEDULED"
	JOB_STATE_RUNNING   JobState = "JOB_STATE_RUNNING"
	JOB_STATE_SUCCESS   JobState = "JOB_STATE_SUCCESS"
	JOB_STATE_FAILED    JobState = "JOB_STATE_FAILED"
)

type JobType string

const (
	JOB_TYPE_IP_MEASUREMENT JobType = "IT"
)

type Job struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	State     JobState  `json:"state"`
	JobType   JobType   `json:"job_type"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Created   time.Time
	Modified  time.Time
}

type Result struct {
	Data    []*Job
	Error   string
	Success bool
}
