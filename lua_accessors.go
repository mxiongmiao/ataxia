package main

import (
	"fmt"
	golua "github.com/aarzilli/golua/lua"
	"log"
	"strconv"
	//	"github.com/xenith-studios/ataxia/lua"
	luar "github.com/stevedonovan/luar"
)

func (server *Server) PublishAccessors(state *golua.State) {
	// register exported functions (this is a weird place, should be in main?  or called from there?)
	luar.Register(state, "", luar.Map{
		"GetPlayerData": server.GetPlayerData,
		"SendToPlayers": server.SendToPlayers,
	})
}

func (server *Server) SendToPlayers(msg string) {
	for _, player := range server.PlayerList.players {
		if player != nil {
			log.Println(msg)
			player.In <- fmt.Sprintf("%s\r\n", msg)
		}
	}
}

func (server *Server) GetPlayerData(id string, field string) (ret string) {
	player := server.PlayerList.Get(id)
	if field == "name" { // replace this with reflection on struct tags?
		ret = player.account.Name
	}
	return
}

func (world *World) PublishAccessors(state *golua.State) {
	luar.Register(state, "", luar.Map{
		"SendToAll":        world.SendToAll,
		"SendToOthers":     world.SendToOthers,
		"SendToChar":       world.SendToChar,
		"GetCharacterData": world.GetCharacterData,
		"SetCharacterData": world.SetCharacterData,
		"GetRoomData":      world.GetRoomData,
		"GetRoomExit":      world.GetRoomExit,
		"GetRoomExitData":  world.GetRoomExitData,
		"GetObjectData":    world.GetObjectData,
		"SetObjectData":    world.SetObjectData,
		"GetDictData":      world.GetDictData,
	})
}

func (world *World) SendToAll(msg string) {
	for _, ch := range world.Characters {
		if ch.Player != nil {
			log.Println(msg)
			ch.Player.In <- fmt.Sprintf("%s\r\n", msg)
		}
	}
}

func (world *World) SendToOthers(char_id string, msg string) {
	for id, ch := range world.Characters {
		if id == char_id {
			continue
		}

		if ch.Player != nil {
			log.Println(msg)
			ch.Player.In <- fmt.Sprintf("%s\r\n", msg)
		}
	}
}

func (world *World) SendToChar(id string, msg string) {
	ch := world.Characters[id]
	if ch != nil {
		if ch.Player != nil {
			ch.Player.In <- msg
		}
	}
}

func (world *World) GetObjectData(id string, field string) (ret string) {
	if world.Characters[id] != nil {
		return world.GetCharacterData(id, field)
	}
	if world.Rooms[id] != nil {
		return world.GetRoomData(id, field)
	}
	if world.RoomExits[id] != nil {
		return world.GetRoomExitData(id, field)
	}
	return ""
}

func (world *World) SetObjectData(id string, field string, value string) {
	if world.Characters[id] != nil {
		world.SetCharacterData(id, field, value)
	}
	//	if world.Rooms[id] != nil { world.SetRoomData(id, field, value) }
	//	if world.RoomExits[id] != nil { world.SetRoomExitData(id, field, value) }
}

func (world *World) GetDictData(id string, field string, key string) (ret string) {
	if world.Rooms[id] == nil {
		return ""
	}

	room := world.Rooms[id]
	if field == "exits" {
		keyv, _ := strconv.Atoi(key)
		exit := room.exits[keyv]
		if exit != nil {
			return exit.ID
		}
	}

	return ""
}

func (world *World) GetCharacterData(id string, field string) (ret string) {
	ch := world.Characters[id]
	if ch == nil {
		return ""
	}

	if field == "name" {
		return ch.Name
	}
	if field == "room" {
		return ch.Room.Id
	}
	return
}

func (world *World) SetCharacterData(id string, field string, value string) {
	ch := world.Characters[id]
	if ch == nil {
		return
	}

	if field == "room" {
		ch.Room = world.Rooms[value]
		return
	}
}

func (world *World) GetRoomData(id string, field string) (ret string) {
	ch := world.Rooms[id]
	if ch == nil {
		return ""
	}

	if field == "name" {
		return ch.Name
	}
	if field == "description" {
		return ch.Description
	}
	return ""
}

func (world *World) GetRoomExit(room_id string, dir int) (ret string) {
	room := world.Rooms[room_id]
	if room == nil {
		return ""
	}

	if room.exits[dir] != nil {
		return room.exits[dir].ID
	}

	return ""
}

func (world *World) GetRoomExitData(exit_id string, field string) (ret string) {
	exit := world.RoomExits[exit_id]
	if exit == nil {
		return ""
	}

	if field == "destination" {
		if exit.destination == nil {
			return ""
		}

		return exit.destination.Id
	}
	return ""
}