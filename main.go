package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	// fmt.Println("started server ")
	// r := router.Router()

	// log.Fatal(http.ListenAndServe(":4000", r))
	// fmt.Println("Listening at port 4000")
	//cannot execute ping directly in golang runtime using exec
	ipaddress := "8.8.8.8"
	count := 5
	result := make(map[string]interface{})

	cmd := exec.Command("ping", "-c", strconv.Itoa(count), ipaddress)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	done := cmd.Run()
	if done != nil {
		log.Fatalf("Failed to execute ping command: %s", done)
	}
	output := stdout.String()
	lines := strings.Split(output, "\n")
	// fmt.Println(output)

	// fmt.Println("done here")
	// fmt.Println(lines)
	for _, check := range lines {
		if strings.Contains(check, "packets transmitted") {
			parts := strings.Fields(check)
			fmt.Println(parts)
			transmitted, _ := strconv.Atoi(parts[0])
			fmt.Println(transmitted)
			received, _ := strconv.Atoi(parts[3])
			packetLoss := parts[6]

			result["transmitted"] = transmitted
			result["received"] = received
			result["Packet-loss"] = packetLoss
		}

	}

	fmt.Println(result)
}
