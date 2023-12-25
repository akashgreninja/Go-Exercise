package models

type Request struct {
	Ipaddress string `json:"ipaddress,omitempty"`
	count     int    `json:"count,omitempty"`
}

type Response struct {
	ipaddress   string   `json:"ipaddress,omitempty"`
	count       int      `json:"count,omitempty"`
	transmitted int      `json:"transmitted,omitempty"`
	received    int      `json:"received,omitempty"`
	Packetloss  int      `json:"packetloss,omitempty"`
	latency     *latency `json:"latency,omitempty"`
}

type latency struct {
	min string
	max string
}
