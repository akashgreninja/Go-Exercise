package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getDetailsofthePing(ipaddress string, count int) models.Response {

	result := make(map[string]interface{})

	// var check models.Response
	// if net.ParseIP(ipaddress) == nil {
	// 	fmt.Println("we here")
	// 	return models.Response{Ipaddress: "IP address not valid, please check"}
	// }

	cmd := exec.Command("ping", "-c", strconv.Itoa(count), ipaddress)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	done := cmd.Run()
	if done != nil {

		check := models.Errorlog{
			Ipaddress: ipaddress,
			Count:     count,
		}
		responseJSON, err := json.Marshal(check)
		if err != nil {
			log.Fatalf("Error marshaling response to JSON: %s", err)
		}

		filr, err := os.OpenFile("faillog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("errpr")
		}
		filr.Write(responseJSON)
		return models.Response{Ipaddress: "IP address not valid, please check"}
	}
	output := stdout.String()
	lines := strings.Split(output, "\n")
	// fmt.Println(output)

	// fmt.Println("done here")

	// fmt.Println(output)

	for _, check := range lines {
		if strings.Contains(check, "packets transmitted") {
			parts := strings.Fields(check)
			fmt.Println(parts)
			transmitted, _ := strconv.Atoi(parts[0])
			// fmt.Println(transmitted)
			received, _ := strconv.Atoi(parts[3])
			packetLoss := parts[5]
			fmt.Println(packetLoss)
			// fmt.Println("above thiss ")

			result["transmitted"] = transmitted
			result["received"] = received
			result["packet-loss"] = packetLoss
		}
		if strings.Contains(check, "rtt min/avg/max/mdev") {
			fmt.Println("Trueeeeeeeeee")
			//splitting the equal and non equal sign first
			parts := strings.Split(check, "=")
			// fmt.Println(parts)
			//just a check to see if the data returned else will goto error
			if len(parts) >= 2 {

				latencyInfo := strings.TrimSpace(parts[1])
				//splitting by /
				latencyParts := strings.Split(latencyInfo, "/")
				// fmt.Println(latencyParts)

				if len(latencyParts) >= 4 {
					result["latency"] = map[string]string{
						"min":  strings.TrimSpace(latencyParts[0]),
						"avg":  strings.TrimSpace(latencyParts[1]),
						"max":  strings.TrimSpace(latencyParts[2]),
						"mdev": strings.TrimSpace(latencyParts[3]),
					}
				}
			}
		}

	}
	fmt.Println(result["packet-loss"])
	// fmt.Println(result)

	packetLossValue, packetLossOK := result["packet-loss"].(string)

	// Perform type assertions for transmitted and received
	transmittedValue, transmittedOK := result["transmitted"].(int)
	receivedValue, receivedOK := result["received"].(int)

	// Perform type assertions for latency map fields
	latencyMap, latencyMapOK := result["latency"].(map[string]string)
	minValue, minOK := latencyMap["min"]
	avgValue, avgOK := latencyMap["avg"]
	maxValue, maxOK := latencyMap["max"]
	mdevValue, mdevOK := latencyMap["mdev"]

	if packetLossOK && transmittedOK && receivedOK && latencyMapOK && minOK && avgOK && maxOK && mdevOK {
		response := models.Response{
			Ipaddress:   ipaddress,
			Count:       count,
			Transmitted: transmittedValue,
			Received:    receivedValue,
			Packetloss:  packetLossValue,
			Latency: &models.Latency{
				Min:  minValue,
				Avg:  avgValue,
				Max:  maxValue,
				Mdev: mdevValue,
			},
		}

		return response

	} else {
		log.Println("Type assertion failed or values are nil")
		if packetLossValue == "" || !transmittedOK || !receivedOK || !latencyMapOK || minValue == "" || avgValue == "" || maxValue == "" || mdevValue == "" {
			return models.Response{Ipaddress: "IP address not valid, please check"}
		}
	}
	return models.Response{}

}

func Statistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	if r.Body == nil {
		json.NewEncoder(w).Encode("num-failed")
	}
	var result models.Request

	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		json.NewEncoder(w).Encode("num-failed")

	}

	fmt.Println(result.Ipaddress)
	checkers := getDetailsofthePing(result.Ipaddress, result.Count)
	//diff between io.Write string anf file.Write check
	responseJSON, err := json.Marshal(checkers)
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %s", err)
	}

	filr, err := os.OpenFile("successlog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("errpr")
	}
	filr.Write(responseJSON)

	json.NewEncoder(w).Encode(checkers)

}
