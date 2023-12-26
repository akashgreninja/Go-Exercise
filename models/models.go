package models

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
