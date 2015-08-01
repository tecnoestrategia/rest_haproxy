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

type Backends map[string]interface{}

type Backend struct {
	Ip string
}

func getBackend(line string) (backend string, err error) {

	// Regex for Backend
	match_bkend, err := regexp.Compile(`^\s*backend`)
	if err != nil {
		return "Failed", err
	}

	if match_bkend.MatchString(line) {
		log.Println("MATCHED BACKEND: ", line)
		larry := strings.Fields(line)
		name := larry[1]
		backend = name
		log.Println(backend)
		return backend, nil
	}

	return "null", nil
}

func getIp(line string) (ip string, err error) {
	// Regex for Server
	match_srv, err := regexp.Compile(`^\s*server`)
	if err != nil {
		return "Failed.", err // there was a problem with the regular expression.
	}

	if match_srv.MatchString(line) {
		log.Println("MATCHED SERVER:\n", line)
		larry := strings.Fields(line)
		ip := larry[2]
		log.Println(ip)
		return ip, nil
	}
	return "null", nil
}

func parsefile(filename string) (Backends, error) {

	var backends Backends
	backend := make(map[Backend]string)

	// Handle the file
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		backend_name, _ := getBackend(line)
		log.Println("LINE ", line)
		if backend_name != "null" {
			for scanner.Scan() {

				server_ip, _ := getIp(line)

				if server_ip != "null" {

					log.Println("Server: ", server_ip)
					backends[backend_name] = backend[Backend{Ip: server_ip}]
					log.Println("Backends Hash:\n")
					log.Println(backends)
					continue

				}

				// Only pass current backend then THROW
				check_backend, _ := getBackend(line)
				if backends[backend_name] != check_backend {
					log.Println("BREAKING INNER LOOP")
					break
				}
			}
		}
	}

	log.Println("Final Hash:\n")
	log.Println(backends)
	return backends, nil
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
