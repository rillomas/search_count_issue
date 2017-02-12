// +build !appengine

package main

// Script to test add/search API

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	host = flag.String("host", "", "host domain of the targt server (ex. sample.appspot.com)")
	tls  = flag.Bool("secure", false, "true if target server requires https")
)

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

func addRoom(req AddRoomRequest, base string) (*AddRoomResponse, error) {
	addURL := fmt.Sprintf("%s/api/room", base)
	bodyType := "application/json"
	buf, err := json.Marshal(&req)
	if err != nil {
		log.Printf("json marshal failed: %v", err)
		return nil, err
	}
	res, err := http.Post(addURL, bodyType, bytes.NewReader(buf))
	if err != nil {
		log.Printf("http post failed: %v", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Got an unexpected response: %v", res.StatusCode)
		return nil, err
	}
	buf, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return nil, err
	}
	var ares AddRoomResponse
	err = json.Unmarshal(buf, &ares)
	if err != nil {
		log.Printf("json unmarshal failed: %v", err)
		return nil, err
	}
	return &ares, nil
}

func searchRoom(req SearchRoomRequest, base string) (*SearchRoomResponse, error) {
	searchURL := fmt.Sprintf("%s/api/room/search", base)
	bodyType := "application/json"
	buf, err := json.Marshal(&req)
	if err != nil {
		log.Printf("json marshal failed: %v", err)
		return nil, err
	}
	res, err := http.Post(searchURL, bodyType, bytes.NewReader(buf))
	if err != nil {
		log.Printf("http post failed: %v", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Got an unexpected response: %v", res.StatusCode)
		return nil, err
	}
	buf, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return nil, err
	}
	var sres SearchRoomResponse
	err = json.Unmarshal(buf, &sres)
	if err != nil {
		log.Printf("json unmarshal failed: %v", err)
		return nil, err
	}
	return &sres, nil

}

func main() {
	flag.Parse()
	if *host == "" {
		log.Fatalf("-host flag required (ex. sample.appspot.com)")
	}
	protocol := "http"
	if *tls {
		protocol = "https"
	}
	base := fmt.Sprintf("%s://%s", protocol, *host)
	req := AddRoomRequest{
		Name: "aardvark",
	}
	res, err := addRoom(req, base)
	if err != nil {
		log.Fatalf("Failed to add room: %v", err)
		return
	}
	log.Printf("Added room: %s", res.RoomID)
	sreq := SearchRoomRequest{
		Name: "aardvark",
	}
	sres, err := searchRoom(sreq, base)
	if err != nil {
		log.Fatalf("Failed to search room: %v", err)
		return
	}
	log.Printf("Got %d search match.", len(sres.Rooms))
	for _, r := range sres.Rooms {
		log.Printf("%s", r.Name)
	}
}
