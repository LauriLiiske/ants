package main

import (
	"lem-in/logic"
	"log"
)

func main() {
	var TheNest logic.Nest // This is the Graph object with all the rooms and connections

	// This section constructs the graph from data:
	// builds rooms, adds connections and sets the start and end room
	lines := logic.ReadFile()
	RoomsSlice, ConnectionSlice, StartRoom, EndRoom := logic.GetData(lines, &TheNest)

	// Initialize the graphs nodes
	TheNest.Vertices = logic.BobTheBuilder(RoomsSlice)
	TheNest.RoomReference = logic.Tunneler(TheNest.Vertices, ConnectionSlice)

	logic.SetStartEnd(StartRoom, EndRoom, &TheNest)

	// Find possible paths
	logic.Pathfinding(&TheNest)

	if len(TheNest.ConfirmedPaths) == 0 {
		log.Fatal("ERROR: No possible paths found.")
	}

	// Deal with the ants movements
	logic.AddAntsToStart(&TheNest, TheNest.Start)
	logic.GiveAntsRoadsigns(&TheNest)
	antsRoomsResult := logic.AntSwapper(&TheNest)
	
	logic.CreateAndPrintFinalAnswer(&TheNest, RoomsSlice, ConnectionSlice, antsRoomsResult)
}
