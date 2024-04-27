package main

import (
	"fmt"
	"slices"
	"strings"
)

type RunDefinition struct {
	arrivalLine         string
	stopTime            int
	stopType            string
	departureLine       string
	categories          []string
	arrivalDefinition   string
	departureDefinition string
	arrivalDelay        string
	departureDelay      string
}

type ActiveRun struct {
	Content       string
	TimestampFrom int64
}

func (a *App) processData() {
	err := fetchActiveData("https://stacjownik.spythere.eu/api/getActiveData", &a.activeData)

	if err != nil {
		println("Błąd podczas ładowania danych z API", err.Error())
	}

	a.activeRuns = nil

	for _, train := range a.activeData.Trains {
		if train.Timetable.Category == "" || train.Region != "eu" {
			continue
		}

		fmt.Printf("%d %s\n", train.TrainNo, train.Timetable.Category)

		for _, stop := range train.Timetable.StopList {
			if slices.Contains(a.checkpoints, strings.ToLower(stop.StopNameRAW)) {
				arrivalRoute := a.routeMappings[stop.ArrivalLine+"_Wjazd"]
				departureRoute := a.routeMappings[stop.DepartureLine+"_Wyjazd"]

				for _, run := range a.runDefinitions {
					if run.arrivalLine != arrivalRoute ||
						run.departureLine != departureRoute ||
						stop.Confirmed == 1 ||
						(stop.StopTime < run.stopTime) ||
						(run.stopType != "*" && !strings.Contains(stop.StopType, run.stopType)) ||
						(stop.StopType != "" && run.stopType == "") {
						continue
					}

					isCategoryFound := false

					for _, category := range run.categories {
						if strings.HasPrefix(train.Timetable.Category, category) || category == "*" {
							isCategoryFound = true
							break
						}
					}

					if !isCategoryFound {
						continue
					}

					a.activeRuns = append(a.activeRuns, ActiveRun{
						Content:       fmt.Sprintf("przebieg %d - - %s %s", train.TrainNo, run.arrivalDelay, run.arrivalDefinition),
						TimestampFrom: stop.ArrivalTimestamp - 1000*60*10,
					})

					// println(stop.StopNameRAW, arrivalRoute, departureRoute, run.arrivalDefinition)

					if run.departureDefinition != "" {
						a.activeRuns = append(a.activeRuns, ActiveRun{
							Content:       fmt.Sprintf("przebieg %d - - %s %s", train.TrainNo, run.departureDelay, run.departureDefinition),
							TimestampFrom: stop.DepartureTimestamp - 1000*60*1,
						})
					}

					break
				}

			}
		}
	}
}
