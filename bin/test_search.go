// +build !appengine

package main

// Script to test add/search API

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// RoomInfo is information of a room that is returned to the client
type RoomInfo struct {
	Name       string
	CreateTime time.Time
}
type AddRoomRequest struct {
	Name string
}

type AddRoomResponse struct {
	RoomID string // ID of a room
}

type SearchRoomRequest struct {
	Name string
}

type SearchRoomResponse struct {
	Rooms []RoomInfo
}

func addRoom(req AddRoomRequest) (*AddRoomResponse, error) {
	addURL := "http://localhost:8080/api/room"
	bodyType := "application/json"
	buf, err := json.Marshal(&req)
	if err != nil {
		logger.Printf("json marshal failed: %v", err)
		return nil, err
	}
	res, err := http.Post(addURL, bodyType, bytes.NewReader(buf))
	if err != nil {
		logger.Printf("http post failed: %v", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		logger.Printf("Got an unexpected response: %v", res.StatusCode)
		return nil, err
	}
	buf, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Printf("Failed to read response: %v", err)
		return nil, err
	}
	var ares AddRoomResponse
	err = json.Unmarshal(buf, &ares)
	if err != nil {
		logger.Printf("json unmarshal failed: %v", err)
		return nil, err
	}
	return &ares, nil
}

func main() {
	req := AddRoomRequest{
		Name: "aardvark",
	}
	res, err := addRoom(req)
	if err != nil {
		logger.Fatalf("Failed to add room: %v", err)
		return
	}
	logger.Printf("Added room: %s", res.RoomID)
}
