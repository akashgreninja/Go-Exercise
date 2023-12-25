package controllers

func getDetailsofthePing(ipaddress string, count int) {
	//cannot execute ping directly in golang runtime using exec
	// cmd := exec.Command("ping", "-c", strconv.Itoa(count), ipaddress)

}
