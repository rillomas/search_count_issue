package sample

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"

	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestSearchRoom(t *testing.T) {
	ins, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer ins.Close()
	// request without lang
	areq := AddRoomRequest{
		Name: "aardvark",
	}
	buf, err := json.Marshal(areq)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}
	req, err := ins.NewRequest("POST", "/api/room", bytes.NewReader(buf))
	if err != nil {
		t.Errorf("Failed creating request: %v", err)
	}
	w := httptest.NewRecorder()
	handleAddRoom(w, req)
	if w.Code != 200 {
		t.Errorf("Should get OK: %d", w.Code)
	}
	var res AddRoomResponse
	err = json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
	t.Logf("RoomID: %s", res.RoomID)
	// check if room is really added
	c := appengine.NewContext(req)
	k, err := datastore.DecodeKey(res.RoomID)
	if err != nil {
		t.Errorf("Failed to decode key: %v", err)
	}
	var sr StoreRoom
	err = datastore.Get(c, k, &sr)
	if err != nil {
		t.Errorf("Failed to get stored room: %v", err)
	}
	if sr.Name != areq.Name {
		t.Errorf("Wrong room was added: %s", sr.Name)
	}
	// search for room
	searchName := "aardvark"
	sreq := SearchRoomRequest{
		Name: searchName,
	}
	buf, err = json.Marshal(sreq)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}
	req, err = ins.NewRequest("POST", "/api/room/search", bytes.NewReader(buf))
	if err != nil {
		t.Errorf("Failed creating request: %v", err)
	}
	w = httptest.NewRecorder()
	handleSearchRoom(w, req)
	if w.Code != 200 {
		t.Errorf("Should get OK: %d", w.Code)
	}
	var sres SearchRoomResponse
	err = json.Unmarshal(w.Body.Bytes(), &sres)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
	if len(sres.Rooms) != 1 {
		t.Fatalf("Should get single search result: %d", len(sres.Rooms))
	}
	rm := sres.Rooms[0]
	if rm.Name != searchName {
		t.Fatalf("Got wrong room: %v", rm.Name)
	}
}
