package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		domains := strings.Split(scanner.Text(), " ")
		for _, v := range domains {
			checkDomain(v)
			fmt.Println()
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("error: could not read from input: %v", err)
		}
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

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
	fmt.Printf("%sDomain%s: %s\n%sHas MX%s: %v\n%sHas SPF%s: %v\n%sSPF Record%s: %s\n%sHas DMARC%s: %v\n%sDMARC Record%s: %s\n",
		Blue, Reset, domain,
		Green, Reset, hasMX,
		Yellow, Reset, hasSPF,
		Yellow, Reset, spfRecord,
		Purple, Reset, hasDMARC,
		Purple, Reset, dmarcRecord)
}
