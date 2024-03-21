package logic

import (
	"log"
	"strconv"
	"strings"
)

type Room struct {
	Name             string
	Coord_x, Coord_y int
	IsStart          bool
	IsEnd            bool
	Connections      []*Room // Other Room objects that are connected to this one
	Roadsign         *Room   // Pointer towards the next room on the path
	Visited          bool    // Flag for the BFS to use
	IsOnPath         bool    // Flag for skipping Visited reset
	Occupants        []Ant   // Current ant(s) in the room
}

// Marks the room as visited
func (r *Room) SetVisited() {
	r.Visited = true
}

// Marks the room as unvisited
func (r *Room) SetUnvisited() {
	r.Visited = false
}

func (r *Room) MarkOnPath() {
	r.SetVisited()
	r.IsOnPath = true
}

// Adds a connection to another room.
func (r *Room) addConnection(link *Room) {
	r.Connections = append(r.Connections, link)
}

// Sets the room to be a start room
func (r *Room) SetStart() {
	r.IsStart = true
}

// Sets the room to be an end room
func (r *Room) SetEnd() {
	r.IsEnd = true
}

// Builds the rooms using roomify() function.
// roomsSlice: slice of rooms with their name and coordinates
func BobTheBuilder(roomsSlice []string) []Room {
	var listOfRooms []Room
	for _, room := range roomsSlice {
		listOfRooms = append(listOfRooms, roomBuilder(room))
	}
	return listOfRooms
}

// Links the rooms and creates a map with the room adresses.
// Map: <room_name> - <room_address> key-value pair
// Also adds connections between the rooms with linkRooms
func Tunneler(vertices []Room, connectionSlice []string) map[string]*Room {
	theMap := mapRooms(vertices)
	linkRooms(connectionSlice, theMap)
	return theMap
}

// Turns the slice of rooms objects into a map.
// Map: <room_name> - <room_address> key-value pair
func mapRooms(roomList []Room) map[string]*Room {
	theMap := make(map[string]*Room, len(roomList))
	for i := range roomList {
		theMap[roomList[i].Name] = &roomList[i] // assign each rooms name a value that is the address of the rooms pointer (like 0xc000106138)
	}
	return theMap
}

// Adds connections between all the rooms
func linkRooms(connectionSlice []string, theMap map[string]*Room) {
	for _, connection := range connectionSlice {
		rooms := strings.Split(connection, "-")

		// Adds a connection between the rooms
		// Both rooms get a pointer to the other room
		theMap[rooms[0]].addConnection(theMap[rooms[1]])
		theMap[rooms[1]].addConnection(theMap[rooms[0]])
	}
}

// Adds tags (of start and end) to corresponding rooms
func SetStartEnd(startString string, endString string, Graph *Nest) {
	startName := getName(startString)
	endName := getName(endString)
	Graph.RoomReference[startName].SetStart()
	Graph.Start = Graph.RoomReference[startName]
	Graph.RoomReference[endName].SetEnd()
	Graph.End = Graph.RoomReference[endName]
}

// Turns input room info (name, x y coords) to a Room struct/object.
func roomBuilder(roomData string) Room {
	var NewRoom Room
	var err error

	dataSlice := strings.Fields(roomData)

	NewRoom.Name = dataSlice[0]
	NewRoom.Coord_x, err = strconv.Atoi(dataSlice[1])
	if err != nil {
		log.Fatal()
	}
	NewRoom.Coord_y, err = strconv.Atoi(dataSlice[2])
	if err != nil {
		log.Fatal()
	}
	return NewRoom
}

// Returns the room name from the roomString
// String: "<room_name> <x_coord> <y_coord>"
func getName(roomString string) string {
	name, _, ok := strings.Cut(roomString, " ")
	if !ok {
		log.Fatal("INCORRECT STRING: no whitespace found.\n")
	}
	return name
}
