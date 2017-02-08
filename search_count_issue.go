package sample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
)

const (
	roomName = "room"
	roomKind = "room"
)

// SearchRoom is information of a room used for searching
type SearchRoom struct {
	Name       string
	CreateTime time.Time
}

// StoreRoom is information of a room that is stored in the datastore
type StoreRoom struct {
	Name       string
	CreateTime time.Time
}

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

func init() {
	http.HandleFunc("/api/room", handleAddRoom)
	http.HandleFunc("/api/room/search", handleSearchRoom)
}

func handleAddRoom(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "Failed to ReadAll: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := new(AddRoomRequest)
	err = json.Unmarshal(buf, req)
	if err != nil {
		log.Errorf(c, "Failed to unmarshal: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// store room in datastore
	ct := time.Now()
	sr := StoreRoom{
		Name:       req.Name,
		CreateTime: ct,
	}
	rk := datastore.NewIncompleteKey(c, roomKind, nil)
	rk, err = datastore.Put(c, rk, &sr)
	if err != nil {
		log.Errorf(c, "Failed to store room: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// store room in search index
	rm := SearchRoom{
		Name:       req.Name,
		CreateTime: ct,
	}
	idx, err := search.Open(roomName)
	if err != nil {
		log.Errorf(c, "Failed to open search index: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	roomID := rk.Encode()
	_, err = idx.Put(c, roomID, &rm)
	if err != nil {
		log.Errorf(c, "Failed to store search index: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := AddRoomResponse{
		RoomID: roomID,
	}
	// send response
	outBuf, err := json.Marshal(res)
	if err != nil {
		log.Errorf(c, "Failed to marshal response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(outBuf)
	if err != nil {
		log.Errorf(c, "Failed to write output: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSearchRoom(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "Failed to ReadAll: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := new(SearchRoomRequest)
	err = json.Unmarshal(buf, req)
	if err != nil {
		log.Errorf(c, "Failed to unmarshal: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := fmt.Sprintf("Name: %s", req.Name)
	idx, err := search.Open(roomName)
	if err != nil {
		log.Errorf(c, "Failed to open search index: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	so := search.SearchOptions{
		IDsOnly: true,
	}
	itr := idx.Search(c, q, &so)
	log.Infof(c, "Search count: %d", itr.Count())
	var res SearchRoomResponse
	for {
		id, err := itr.Next(nil)
		if err == search.Done {
			break
		}
		if err != nil {
			log.Errorf(c, "Failed to iterate search result: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rk, err := datastore.DecodeKey(id)
		if err != nil {
			log.Errorf(c, "Failed to iterate search result: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var sr StoreRoom
		err = datastore.Get(c, rk, &sr)
		if err != nil {
			log.Errorf(c, "Failed to get room: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ri := RoomInfo{
			Name: sr.Name,
		}
		res.Rooms = append(res.Rooms, ri)
	}
	// send response
	outBuf, err := json.Marshal(res)
	if err != nil {
		log.Errorf(c, "Failed to marshal response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(outBuf)
	if err != nil {
		log.Errorf(c, "Failed to write output: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
