package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

func main() {
	listenaddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenaddr = ":" + val
	}
	http.HandleFunc("/api/check_email", checkDomain)

	log.Printf("%sListening On:%s http://127.0.0.1%s", Red, Reset, listenaddr)
	log.Fatal(http.ListenAndServe(listenaddr, nil))
}

func checkDomain(w http.ResponseWriter, r *http.Request) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	message := fmt.Sprintf("Domain: %s\nHas MX: %v\nHas SPF: %v\nSPF Record: %s\nHas DMARC: %v\nDMARC Record: %s\n",
		domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
	fmt.Fprint(w, message)
}
