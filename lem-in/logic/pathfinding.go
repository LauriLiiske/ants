package logic

import (
	"log"
	"sort"
)

type Path struct {
	// Each Room object in the path has its roadsign pointing to the next Room
	Length     int   // Length of the path in steps between Start - End
	Roadsign   *Room // The address of the first room from Start
	AntsOnPath int   // amount of ants currently on the path
}

// Struct for the graph
type Nest struct {
	CountOfAnts    int
	Vertices       []Room
	RoomReference  map[string]*Room
	Start, End     *Room
	ConfirmedPaths []Path
}

// Searches the given Nest object for possible paths
// Appends the results in the ConfirmedPath field
func Pathfinding(Graph *Nest) {
	var successfulPath Path
	var allPaths [][]*Room
	Graph.Start.SetVisited()
	AllPossiblePaths(Graph, Graph.Start, []*Room{}, &allPaths)

	// Sorts found solutions by length
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	finalPaths := BestCombOfPaths(Non_IntersectingPaths(allPaths), Graph.CountOfAnts)

	for _, path := range finalPaths {
		successfulPath = buildPath(path)
		Graph.ConfirmedPaths = append(Graph.ConfirmedPaths, successfulPath)
	}
}

// Searches the graph for a path. Returns the end room with the Roadsign pointing to the parent Room.
// Returns nil when no end was found
func BreathFirstSearch(Graph *Nest) *Room {
	var Q []*Room // Room queue for the algorithm to visit
	var v *Room   // Room currently being visited

	Graph.Start.SetVisited()
	Q = append(Q, Graph.Start)
	for len(Q) > 0 {
		v, Q = Q[0], Q[1:] // Takes out the first element in the queue and assigns it to v
		if v.IsEnd {       // End is found and returned
			resetVisits(Graph)
			return v
		}
		// Visits every connection of the room, sets the roadsign and adds them to the queue
		for _, connection := range v.Connections {
			// Edgecase when start and end are directly connected.
			// Skips this iteration, otherwise gets stuck in a loop
			if v.IsStart && connection.IsEnd && len(Graph.ConfirmedPaths) > 0 {
				continue
			}

			if !connection.Visited {
				connection.SetVisited()
				connection.Roadsign = v // Connected rooms roadsign set to point to the current room v
				Q = append(Q, connection)
			}
		}
	}
	return nil
}

// Runs though the BreathFirstSearch's return and turns the roadsigns to point towards the end
// Parameter recieved by the function should be an endroom
func PathTraceback(activeRoom *Room) Path {
	if !activeRoom.IsEnd {
		log.Fatalln("INCORRECT PARAMATER: Recieved *Room is not End room")
	}

	// The direction of movement for prevRoom and nextRoom is
	// only within the scope of this function.
	// This is the opposite direction of the resulting Path
	// that will lead from the Start to End
	var prevRoom, nextRoom *Room
	var PathLength int = 1 // Starts at one to count the step in the next section

	// This section sets the End room to be previous room
	// and sets up for the loop to take over
	prevRoom = activeRoom
	activeRoom = prevRoom.Roadsign

	for !activeRoom.IsStart {
		activeRoom.MarkOnPath() // Marks the Room as on path so the path search would skip them next time
		nextRoom = activeRoom.Roadsign
		activeRoom.Roadsign = prevRoom
		activeRoom, prevRoom = nextRoom, activeRoom
		PathLength++
	}
	return Path{Length: PathLength, Roadsign: prevRoom}
}

func buildPath(path []*Room) Path {
	newPath := Path{
		Length:     len(path) - 1,
		Roadsign:   path[1],
		AntsOnPath: 0,
	}
	for i := 1; !path[i].IsEnd; i++ {
		path[i].Roadsign = path[i+1]
	}
	return newPath
}

// Turns a Path object into a human-readable slice of room names in order
// Mostly for testing
func ReadablePath(pathOfInterest Path, Graph Nest) []string {
	var ResultingPath []string
	ResultingPath = append(ResultingPath, Graph.Start.Name)
	currentRoom := pathOfInterest.Roadsign
	for !currentRoom.IsEnd {
		ResultingPath = append(ResultingPath, currentRoom.Name)
		currentRoom = currentRoom.Roadsign
	}
	return append(ResultingPath, currentRoom.Name)
}

func resetVisits(Graph *Nest) {
	for i := range Graph.Vertices {
		if !Graph.Vertices[i].IsOnPath {
			Graph.Vertices[i].SetUnvisited()
		}
	}
}

// Finds all paths with a Depth First search
func AllPossiblePaths(Graph *Nest, currentRoom *Room, path []*Room, allPaths *[][]*Room) {
	newPath := make([]*Room, len(path)+1)
	copy(newPath, path)
	newPath[len(path)] = currentRoom

	if currentRoom.IsEnd {
		*allPaths = append(*allPaths, newPath)
		return
	}

	for _, connection := range currentRoom.Connections {
		if connection.Visited {
			continue
		}
		connection.SetVisited()
		AllPossiblePaths(Graph, connection, newPath, allPaths)
		connection.SetUnvisited()
	}
}

func Non_IntersectingPaths(allPaths [][]*Room) map[int][][]*Room {
	result := make(map[int][][]*Room)
	for i, path1 := range allPaths {
		result[i] = append(result[i], path1)
		for j, path2 := range allPaths {
			if i == j {
				continue
			}
			if !intersecting(path1, path2) {
				intersects := false
				for _, path3 := range result[i] {
					if intersecting(path2, path3) {
						intersects = true
						break
					}
				}
				if !intersects {
					result[i] = append(result[i], path2)
				}
			}
		}
	}
	return result
}

func intersecting(path1, path2 []*Room) bool { // helper for func getNonIntersectingPaths
	for _, room1 := range path1[1 : len(path1)-1] {
		for _, room2 := range path2[1 : len(path2)-1] {
			if room1.Name == room2.Name {
				return true
			}
		}
	}
	return false
}

func BestCombOfPaths(combOfPaths map[int][][]*Room, ants int) [][]*Room { // find best combinations of paths depending on the number of ants
	bestCombination := [][]*Room{}
	finalSteps := 0
	for _, comb := range combOfPaths {
		numOfSteps := ants
		for _, path := range comb {
			numOfSteps += len(path)
		}
		buffinalSteps := numOfSteps / len(comb)
		if finalSteps == 0 || finalSteps > buffinalSteps {
			finalSteps = buffinalSteps
			bestCombination = comb
		}
	}
	return bestCombination
}
