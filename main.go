package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

type Person struct {
	Craft string `json:"craft"`
	Name  string `json:"name"`
}

type Astronauts struct {
	Number int      `json:"number"`
	People []Person `json:"people"`
}

// Write astronauts to CSV file
func writeToFile(filename string, astronauts Astronauts) {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Close file when function ends
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Write to file
	file.WriteString("Name,Craft\n")
	for _, astronaut := range astronauts.People {
		file.WriteString(fmt.Sprintf("%s,%s\n", astronaut.Name, astronaut.Craft))
	}
}

// Get the astronaut data from the specified URL
func getData(url string) []byte {
	// Create the request
	request, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		panic(reqErr)
	}

	// Create the client and do the request
	client := http.Client{}
	response, getErr := client.Do(request)
	if getErr != nil {
		panic(getErr)
	}

	// Close body when function ends
	if response.Body != nil {
		defer response.Body.Close()
	}

	// Read the body of the response
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		panic(readErr)
	}

	return body
}

// Parse the JSON data
func parseAstronauts(body []byte) Astronauts {
	astronauts := Astronauts{}
	parseErr := json.Unmarshal(body, &astronauts)
	if parseErr != nil {
		panic(parseErr)
	}

	return astronauts
}

// Sort astronauts by craft, reverse alphabetically, then by name alphabetically
func sortAstronauts(astronauts *Astronauts) {
	// Sort by craft, reverse alphabetically
	sort.SliceStable(astronauts.People, func(i, j int) bool {
		leftCraft := astronauts.People[i].Craft
		rightCraft := astronauts.People[j].Craft

		// If crafts are equal, sort by name alphabetically
		if leftCraft == rightCraft {
			return astronauts.People[i].Name < astronauts.People[j].Name
		} else {
			return leftCraft > rightCraft
		}
	})
}

// http://api.open-notify.org/astros.json
func main() {
	url := "http://api.open-notify.org/astros.json"
	body := getData(url)
	astronauts := parseAstronauts(body)
	sortAstronauts(&astronauts)

	// Print as CSV
	fmt.Println("Name,Craft")
	for _, astronaut := range astronauts.People {
		fmt.Printf("%s,%s\n", astronaut.Name, astronaut.Craft)
	}

	// Write to CSV file
	writeToFile("output.csv", astronauts)
}
