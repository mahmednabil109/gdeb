package data

import (
	"encoding/json"
	"log"
	"os"
)

type Stakeholder struct {
	PublicKey string
	Stake     float64 //percentage
}

type Distribution struct {
	Stakeholders []Stakeholder
}

//loads json file containing stake distribution into in memory map to access stakes of users
//useful when validating transactions (check for enough credit) and validating leaders (computing threshold for leader)
func LoadStakeDist(file string, stakeDist *map[string]float64) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Println("opening config file", err.Error())
	}

	var distribution Distribution
	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&distribution); err != nil {
		log.Println("parsing config file: ", err.Error())
	}
	*stakeDist = make(map[string]float64, len(distribution.Stakeholders))
	for _, elem := range distribution.Stakeholders {
		(*stakeDist)[elem.PublicKey] = elem.Stake
	}
}
