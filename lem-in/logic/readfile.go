package logic

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadFile() []string {
	if len(os.Args) < 2 {
		log.Fatal("ERROR: No filename specified.")
	}
	file, err := os.Open("./examples/" + os.Args[1])
	if err != nil {
		fmt.Println(err)
		log.Fatal("ERROR: There was a problem with your input file.")
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		if scanner.Text() != "" {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}

func GetData(lines []string, n *Nest) ([]string, []string, string, string) {
	// slice of rooms data (name, x, y)
	var RoomsSlice []string
	// slice of rooms connections (room-room)
	var ConnectionsSlice []string
	// start room (name, x, y)
	var StartRoom string
	// end room (name, x, y)
	var EndRoom string

	for i, line := range lines {
		if i == 0 {
			n.CountOfAnts = validateAndGetAntsCount(line)
			if n.CountOfAnts <= 0 {
				log.Fatal("ERROR: The amount of ants is incorrect.")
			}
		} else if line == "##start" {
			if i+1 <= len(lines) { // check if next line exists
				validateRoom(lines[i+1])
				StartRoom = lines[i+1]
			} else { // if there is no next line then there is no starting room
				log.Fatal("ERROR: There are no start room coords.")
			}
		} else if line == "##end" {
			if i+1 <= len(lines) {
				validateRoom(lines[i+1])
				EndRoom = lines[i+1]
			} else {
				log.Fatal("ERROR: There are no end room coords.")
			}
		} else if line[0] == '#' {
			continue
		} else if validateConnection(line) { // we have correct connection slice
			ConnectionsSlice = append(ConnectionsSlice, line) // add connection to connections
		} else if validateRoom(line) { //we have correct room coords
			RoomsSlice = append(RoomsSlice, line) // add room to roomsSlice
		}
	}
	return RoomsSlice, ConnectionsSlice, StartRoom, EndRoom
}

// check if ants count is an integer
func validateAndGetAntsCount(line string) int {
	countOfAnts, err := strconv.Atoi(line) // check if ants count is an integer
	if err != nil {
		log.Fatal("ERROR: The amount of ants is not a number.")
	}
	return countOfAnts
}

// validate the room, true if rooms coords are correct, false if something is wrong
func validateRoom(line string) bool {
	data := strings.Fields(line)

	_, err := strconv.Atoi(data[1]) // convert string to int (we make sure that the coord is a number and not letters)
	if err != nil {
		log.Fatal("ERROR: Rooms coordinate X can not be a string.")
		return false
	}

	_, err2 := strconv.Atoi(data[2]) // convert string to int
	if err2 != nil {
		log.Fatal("ERROR: Rooms coordinate X can not be a string.")
		return false
	}
	return len(data) == 3 // check that the room coords have 3 different values
}

//returns true if the connection is valid
func validateConnection(line string) bool {
	if !strings.Contains(line, "-") { // if the connection does not have - in it
		return false
	}
	data := strings.Split(line, "-")
	return len(data) == 2 //check that there are only 2 parts to the connection slice
}
