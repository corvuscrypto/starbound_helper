package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var lastString string

var userMap map[string]bool
var serverOnline = false

var filePaths struct {
	StarboundLog     string
	OutputDir        string
	StarboundAddress string
	StarboundPort    string
}

func pingTest() {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", net.JoinHostPort(filePaths.StarboundAddress, filePaths.StarboundPort))
	var nserverOnline bool
	for {
		check, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			nserverOnline = false
		} else {
			nserverOnline = true
			check.Close()
		}
		if nserverOnline != serverOnline {
			serverOnline = nserverOnline
			if !serverOnline {
				userMap = make(map[string]bool)
			}
			updateFile()
		}
		<-time.Tick(60 * time.Second)
	}
}

func handleInfo(data []string) {

	for _, v := range data {
		if strings.Contains(v, "Info:") {

			if strings.Contains(v, "Logged in") {
				v = strings.Split(v, "as player '")[1]
				username := strings.Split(v, "'")[0]

				userMap[username] = true

				if !serverOnline {
					serverOnline = true
				}

			} else if strings.Contains(v, "disconnected") {
				v = strings.Split(v, "Client '")[1]
				username := strings.Split(v, "'")[0]
				delete(userMap, username)
			} else if strings.Contains(v, "Server shutdown gracefully") {
				userMap = make(map[string]bool)
				serverOnline = false
			}
		}

	}
	updateFile()
}

func updateFile() {
	var userlist []string
	for k := range userMap {
		userlist = append(userlist, k)
	}

	var serverInfo struct {
		Status string
		Users  []string
	}
	if serverOnline {
		serverInfo.Users = userlist
	}
	if serverOnline {
		serverInfo.Status = "Online"
	} else {
		serverInfo.Status = "Offline"
	}

	data, _ := json.Marshal(serverInfo)

	ioutil.WriteFile(filePaths.OutputDir+"/starbound_data.json", data, os.ModeExclusive)

}

func grabInfo() {
	file, err := os.Open(filePaths.StarboundLog)
	if err != nil {
		log.Fatal(err)
	}
	var newlines []string
	readReady := false
	reader := bufio.NewScanner(file)

	for reader.Scan() {
		if !readReady && lastString != "" {
			if reader.Text() == lastString {
				readReady = true
			}
		} else {
			readReady = true
			newlines = append(newlines, reader.Text())
			if reader.Text() != "" {
				lastString = reader.Text()
			}
		}
	}
	//we scanned whole file and found nothing of our last read line
	//We must read the entire file into the program again, 100 lines at a time
	//to keep memory low still. Also reset the userMap
	if !readReady {
		userMap = make(map[string]bool)
		file.Seek(0, 0)
		reader = bufio.NewScanner(file)
		counter := 0
		for reader.Scan() {
			counter++
			newlines = append(newlines, reader.Text())
			if reader.Text() != "" {
				lastString = reader.Text()
			}
			if counter == 100 {
				handleInfo(newlines)
				newlines = []string{}
				counter = 0
			}
		}
	}

	if len(newlines) > 0 {
		handleInfo(newlines)
	}
}

func poll() {

	for {
		timer := time.Tick(250 * time.Millisecond)
		<-timer
		grabInfo()
	}

}

func loadConfig() {
	dir, _ := os.Getwd()
	file, err := os.Open(dir + "/paths.cfg")
	if err != nil {
		log.Fatal(err)

	}
	data, _ := ioutil.ReadAll(file)
	sData := string(data)
	for _, v := range strings.Split(sData, "\n") {

		if len(v) < 2 {
			continue
		}
		if string(strings.TrimSpace(v)[0]) == "#" {
			continue
		}

		if strings.Contains(v, "=") {
			kvPair := strings.Split(v, "=")
			kvPair[0] = strings.TrimSpace(kvPair[0])
			kvPair[1] = strings.Replace(strings.TrimSpace(kvPair[1]), "\"", "", -1)

			switch kvPair[0] {
			case "starbound_log":
				filePaths.StarboundLog = kvPair[1]
			case "output_directory":
				filePaths.OutputDir = kvPair[1]
			case "starbound_address":
				filePaths.StarboundAddress = kvPair[1]
			case "starbound_port":
				filePaths.StarboundPort = kvPair[1]
			}
		} else {
			continue
		}

	}
}

func main() {
	loadConfig()
	go pingTest()
	userMap = make(map[string]bool)
	go poll()
	done := make(chan bool)
	<-done
}
