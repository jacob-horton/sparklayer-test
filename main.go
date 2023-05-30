package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

// http://api.open-notify.org/astros.json
func main() {
	url := "http://api.open-notify.org/astros.json"

	// Create the request
	request, getErr := http.NewRequest(http.MethodGet, url, nil)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Create the client and do the request
	client := http.Client{}
	response, getErr := client.Do(request)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Check the response wasn't empty
	if response.Body != nil {
		defer response.Body.Close()
	}

	// Read the body of the response
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// Parse the JSON data
	astronauts := Astronauts{}
	parseErr := json.Unmarshal(body, &astronauts)
	if parseErr != nil {
		log.Fatal(parseErr)
	}

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

	// Print as CSV
	for _, astronaut := range astronauts.People {
		fmt.Printf("%s,%s\n", astronaut.Name, astronaut.Craft)
	}
}
