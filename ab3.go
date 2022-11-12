package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func HomeEndpoint(w http.ResponseWriter, r *http.Request) {
	//Notice the metadata url is not the standard EC2 169.254.169.254 address due to the Fargate instances
	resp, err := http.Get("http://169.254.170.2/v2/metadata")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		jsonTxt := scanner.Text()
		fmt.Fprintln(w, jsonTxt)
		//TODO [Used to demonstrate a bad deployment - will result in no code returned, comment out the previous line]
		//jsonTxtUpper := strings.ToUpper(jsonTxt)
		//fmt.Fprintln(w, jsonTxtUpper)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", HomeEndpoint)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
