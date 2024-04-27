package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type ActiveData struct {
	Trains []struct {
		ID                 string `json:"id"`
		TrainNo            int    `json:"trainNo"`
		Mass               int    `json:"mass"`
		Speed              int    `json:"speed"`
		Length             int    `json:"length"`
		Distance           int    `json:"distance"`
		StockString        string `json:"stockString"`
		DriverName         string `json:"driverName"`
		DriverID           int    `json:"driverId"`
		DriverIsSupporter  bool   `json:"driverIsSupporter"`
		DriverLevel        int    `json:"driverLevel"`
		CurrentStationHash string `json:"currentStationHash"`
		CurrentStationName string `json:"currentStationName"`
		Signal             string `json:"signal"`
		ConnectedTrack     string `json:"connectedTrack"`
		Online             int    `json:"online"`
		LastSeen           int64  `json:"lastSeen"`
		Region             string `json:"region"`
		IsTimeout          bool   `json:"isTimeout"`
		Timetable          struct {
			Skr      bool   `json:"SKR"`
			Twr      bool   `json:"TWR"`
			Category string `json:"category"`
			StopList []struct {
				StopName               string  `json:"stopName"`
				StopNameRAW            string  `json:"stopNameRAW"`
				StopType               string  `json:"stopType"`
				StopDistance           float32 `json:"stopDistance"`
				PointID                string  `json:"pointId"`
				Comments               string  `json:"comments"`
				MainStop               bool    `json:"mainStop"`
				ArrivalLine            string  `json:"arrivalLine"`
				ArrivalTimestamp       int64   `json:"arrivalTimestamp"`
				ArrivalRealTimestamp   int64   `json:"arrivalRealTimestamp"`
				ArrivalDelay           int     `json:"arrivalDelay"`
				DepartureLine          string  `json:"departureLine"`
				DepartureTimestamp     int64   `json:"departureTimestamp"`
				DepartureRealTimestamp int64   `json:"departureRealTimestamp"`
				DepartureDelay         int     `json:"departureDelay"`
				BeginsHere             bool    `json:"beginsHere"`
				TerminatesHere         bool    `json:"terminatesHere"`
				Confirmed              int     `json:"confirmed"`
				Stopped                int     `json:"stopped"`
				StopTime               int     `json:"stopTime"`
			} `json:"stopList"`
			Route       string   `json:"route"`
			TimetableID int      `json:"timetableId"`
			Sceneries   []string `json:"sceneries"`
		} `json:"timetable,omitempty"`
	} `json:"trains"`
	ActiveSceneries []struct {
		DispatcherID            int    `json:"dispatcherId"`
		DispatcherName          string `json:"dispatcherName"`
		DispatcherIsSupporter   bool   `json:"dispatcherIsSupporter"`
		StationName             string `json:"stationName"`
		StationHash             string `json:"stationHash"`
		Region                  string `json:"region"`
		MaxUsers                int    `json:"maxUsers"`
		CurrentUsers            int    `json:"currentUsers"`
		Spawn                   int    `json:"spawn"`
		LastSeen                int64  `json:"lastSeen"`
		DispatcherExp           int    `json:"dispatcherExp"`
		NameFromHeader          string `json:"nameFromHeader"`
		SpawnString             string `json:"spawnString"`
		NetworkConnectionString string `json:"networkConnectionString"`
		IsOnline                int    `json:"isOnline"`
		DispatcherRate          int    `json:"dispatcherRate"`
		DispatcherStatus        int    `json:"dispatcherStatus"`
	} `json:"activeSceneries"`
	APIStatuses struct {
		StationsAPI            string `json:"stationsAPI"`
		TrainsAPI              string `json:"trainsAPI"`
		DispatchersAPI         string `json:"dispatchersAPI"`
		SceneryRequirementsAPI string `json:"sceneryRequirementsAPI"`
	} `json:"apiStatuses"`
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

func fetchActiveData(url string, target interface{}) error {
	r, err := httpClient.Get(url)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
