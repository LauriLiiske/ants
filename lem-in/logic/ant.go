package logic

import (
	"fmt"
	"strconv"
	"strings"
)

type Ant struct {
	AntName     string
	AntRoadsign *Room // address of the ants next room
	AntRoom     *Room // address of the current room the ant is in
}

// Adds ant to specified room
func (r *Room) AddAnt(ant Ant) {
	r.Occupants = append(r.Occupants, ant)
}

// Give ant information of what room he is in
func (ant *Ant) ChangeAntRoom(roomAddress *Room) {
	ant.AntRoom = roomAddress
}

// initialize ants structs and add them to the starting room
func AddAntsToStart(n *Nest, r *Room) {
	for i := 1; i <= n.CountOfAnts; i++ {
		ant := Ant{AntName: ("L" + strconv.Itoa(i)), AntRoom: r}
		r.AddAnt(ant)
	}
}

//Moves the ants from start to finish and gives final answer
func AntSwapper(nest *Nest) [][]string {
	var resultSlice [][]string
	var tempResultSlice []string

	var nextRoom *Room
	startAntsCopy := make([]Ant, len(nest.Start.Occupants)) 
	copy(startAntsCopy, nest.Start.Occupants) // make copy of the start ants so that we dont mess up the indexing later

	antsPointers := antsIntoPointers(startAntsCopy)
	var ant *Ant
	
	for !checkEndGoal(nest) { // loop until all the ants have reached the last room.
		for i := 0; i < nest.CountOfAnts; i++ { // step counter
			for j := 0; j < nest.CountOfAnts; j++ { // specific ant counter
				ant = antsPointers[j] 
				nextRoom = ant.AntRoadsign 
				if ant.AntRoom.IsEnd {
					//tempResultSlice = append(tempResultSlice, ant.AntName + "-" + ant.AntRoom.Name)
					continue
				}
				if len(nextRoom.Occupants) >= 1 && !nextRoom.IsEnd { // only skip the ant if the next room has an ant in it, and the next room is not an end room.
					continue
				}
				if len(ant.AntRoom.Occupants) != 0 { // if there is one ant in the antroom occupants slice
					ant.AntRoom.Occupants = removeCorrectAnt(ant, ant.AntRoom.Occupants) // remove the ant from the room's occupants list
					ant.ChangeAntRoom(nextRoom) // tell the ant what room it will be in
					ant.AntRoom.Occupants = append(ant.AntRoom.Occupants, *ant) // add the ant to the next room
					ant.AntRoadsign = ant.AntRoom.Roadsign // set ants new roadsign according to the room the ant is
					tempResultSlice = append(tempResultSlice, ant.AntName + "-" + ant.AntRoom.Name)
				}
			}
			if len(tempResultSlice) != 0 {
				resultSlice = append(resultSlice, tempResultSlice)
				tempResultSlice = nil
			}
			break
		}
	}
	return resultSlice
}

func CreateAndPrintFinalAnswer(n *Nest, RoomsSlice []string, ConnectionSlice []string, antsRoomsResult [][]string) {
	fmt.Println(n.CountOfAnts)
	for _, room := range RoomsSlice {
		strings.Split(room, "")
		if string(room[0]) == n.Start.Name {
			fmt.Println("#start")
			fmt.Println(room)
		} else if string(room[0]) == n.End.Name {
			fmt.Println("#end")
			fmt.Println(room)
		} else {
			fmt.Println(room)
		}
	}
	for _, connection := range ConnectionSlice {
		fmt.Println(connection)
	}
	fmt.Println()
	for _, antAndRoom := range antsRoomsResult {
		fmt.Println(strings.Join(antAndRoom, " "))
	}
}

func PrintFinalAnswer(resultSlice [][]string) {
	for _, slice := range resultSlice {
		fmt.Println(strings.Join(slice, " "))
	}
}

//remove a specific ant from the given Ants slice
func removeCorrectAnt(ant *Ant, antsSlice []Ant) []Ant {
	for i := 0; i < len(antsSlice); i++ {
		if antsSlice[i].AntName == ant.AntName {
			if len(antsSlice) == 1 {
				antsSlice = []Ant{}
				break
			} else {
				antsSlice = append(antsSlice[:i], antsSlice[i+1:]...)
				break
			}
		}
	}
	return antsSlice
}

//turn the Ants slice into Ants pointers Slice
func antsIntoPointers(antsSlice []Ant) []*Ant {
	var antsPointers []*Ant
	for i := range antsSlice {
		antPointer := &antsSlice[i]
		antsPointers = append(antsPointers, antPointer) // append the Ants structs address to the slice of addresses
	}
	return antsPointers
}

//checks if all the ants have reached the last room
func checkEndGoal(n *Nest) bool {
	return n.CountOfAnts == len(n.End.Occupants)
}

// give ants in the starting room their most optimal paths
func GiveAntsRoadsigns(n *Nest) {
	currentPath := 0
	//we always start from the shortest path so safe to give first ant the first path
	for ant := 0; ant < n.CountOfAnts; ant++ { // loopime läbi sipelgad
		if currentPath + 1 < len(n.ConfirmedPaths) {	// kontrollime, et me ei oleks viimasel pathil, sest muidu läheb scoreTwo lappama
			pathOneScore := n.ConfirmedPaths[currentPath].Length + n.ConfirmedPaths[currentPath].AntsOnPath
			pathTwoScore := n.ConfirmedPaths[currentPath+1].Length + n.ConfirmedPaths[currentPath+1].AntsOnPath
			firstPathScore := n.ConfirmedPaths[0].Length + n.ConfirmedPaths[0].AntsOnPath // the FIRST pathscore
			if pathOneScore < pathTwoScore && pathOneScore >= firstPathScore && currentPath != 0 {	// kui praeguse pathi score on väiksem kui järgmise pathi score ja praegune on võrdne või suurem esimesest
				// siis pane sipelgas ikkagi esimesele pathile ja alusta otsast peale
				currentPath = 0	// see on see erijuht, kui järgmine path on palju suurem kui esimene path ja seega sinna pole mõtet pannagi
				(n.Start.Occupants[ant]).AntRoadsign = n.ConfirmedPaths[currentPath].Roadsign
				(&n.ConfirmedPaths[currentPath]).AntsOnPath += 1
			} else if pathOneScore <= pathTwoScore {	// kui praegune path on väiksema või võrdse skooriga võrreldes järgmisega pane sipelgas praegusele pathile
				(n.Start.Occupants[ant]).AntRoadsign = n.ConfirmedPaths[currentPath].Roadsign
				(&n.ConfirmedPaths[currentPath]).AntsOnPath += 1
			} else if pathOneScore > pathTwoScore {	// kui praegune path score on suurem kui järgmise pathi score, panbe sipelgas järgmisele pathile
				currentPath += 1
				(n.Start.Occupants[ant]).AntRoadsign = n.ConfirmedPaths[currentPath].Roadsign
				(&n.ConfirmedPaths[currentPath]).AntsOnPath += 1
			}
		} else {
			pathOneScore := n.ConfirmedPaths[currentPath].Length + n.ConfirmedPaths[currentPath].AntsOnPath
			pathTwoScore := n.ConfirmedPaths[0].Length + n.ConfirmedPaths[0].AntsOnPath
			if pathOneScore <= pathTwoScore {
				(n.Start.Occupants[ant]).AntRoadsign = n.ConfirmedPaths[currentPath].Roadsign
				(&n.ConfirmedPaths[currentPath]).AntsOnPath += 1
			} else {
				currentPath = 0
				(n.Start.Occupants[ant]).AntRoadsign = n.ConfirmedPaths[currentPath].Roadsign
				(&n.ConfirmedPaths[currentPath]).AntsOnPath += 1
			}
		}
	}
}
