package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Services struct {
	Available []map[string]string
}

func get_backend(line string) (name string) {
	// Define our regex to parse
	match_bkend, err := regexp.Compile(`^\s*backend`)
	if match_bkend.MatchString(line) {
		log.Println("MATCHED BACKEND: ", line)
		larry := strings.Fields(line)
		log.Println("BACKEND: ", larry[1])
		name := larry[1]
		return name
		//backend := make(map[string]string)
		//backend[name] = ""
	}
	return name
}

func parsefile(filename string) (*Services, error) {
	// Backend Array Defined
	var backends []map[string]map[string]string

	// Regex for Server
	match_srv, err := regexp.Compile(`^\s*server`)
	if err != nil {
		return nil, err // there was a problem with the regular expression.
	}

	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		current_backend := get_backend(line)
		backend := make(map[string]map[string]string)

		if match_srv.MatchString(line) {
			log.Println("MATCHED SERVER:\n", line)
			larry := strings.Fields(line)
			log.Println("LENGTH: ", len(larry))
			dest := strings.Split(larry[2], ":")
			ip := dest[0]
			port := dest[1]
			if len(larry) == 6 {
				mgmt := larry[5]
				log.Println("MGMT: ", mgmt)
			} else {
				mgmt := dest[1]
				log.Println("MGMT: ", mgmt)
			}
			log.Println("IP: ", ip)
			log.Println("MGMT: ", mgmt)
			log.Println("PORT: ", port)

			backend[current_backend][ip] = ip
			backend[current_backend][port] = port
			backend[current_backend][mgmt] = mgmt
			backends = append(backends, backend)
		}
	}

	return &Services{
			Available: backends,
		},
		nil
}

func response(rw http.ResponseWriter, request *http.Request) {
	services, err := parsefile("/etc/haproxy/haproxy.cfg")
	if err != nil {
		log.Println("ERROR: ", err)
	}
	json, err := json.Marshal(services)
	rw.Write([]byte(json))
}

func main() {
	http.HandleFunc("/services", response)
	http.ListenAndServe(":3000", nil)
}
