package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (a *App) ProcessANPFile(path string) {
	fileBuffer, err := os.ReadFile(path)

	if err != nil {
		panic("Ups")
	}

	a.path = path
	a.fileBuffer = fileBuffer

	content := string(fileBuffer)
	lines := strings.Split(content, "\r\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "posterunek") {
			checkpointData := strings.ToLower(strings.ReplaceAll(strings.Split(line, " ")[1], "_", " "))
			a.checkpoints = append(a.checkpoints, checkpointData)
			println("posterunek", checkpointData)
		}

		if strings.HasPrefix(line, "wersja") {
			a.version = strings.Split(line, " ")[1]
			println("wersja", a.version)
		}

		if strings.HasPrefix(line, "wjazd") {
			runData := strings.Split(line, " ")
			stopTime, err := strconv.Atoi(strings.Replace(runData[2], "-", "0", 1))

			if err != nil {
				println(fmt.Sprintf("Błąd podczas przetwarzania %d linijki", i), err.Error())
				continue
			}

			if len(runData) < 7 {
				println(fmt.Sprintf("Błąd podczas przetwarzania %d linijki", i), "< 7 argumentów")
				continue
			}

			stopType := strings.Replace(runData[3], "-", "", 1)

			if stopType == "" && stopTime > 0 {
				stopType = "pt"
			}

			runDef := &RunDefinition{
				arrivalLine:         runData[1],
				stopTime:            stopTime,
				stopType:            stopType,
				departureLine:       strings.Replace(runData[4], "-", "", 1),
				categories:          strings.Split(runData[5], ","),
				arrivalDefinition:   runData[6],
				departureDefinition: strings.Replace(runData[7], "-", "", 1),
				arrivalDelay:        runData[8],
				departureDelay:      runData[9],
			}

			a.runDefinitions = append(a.runDefinitions, *runDef)

			fmt.Printf("wjazd: %v\n", runDef)
		}

		if strings.HasPrefix(line, "mapuj") {
			mapData := strings.Split(line, " ")

			a.routeMappings[mapData[1]] = mapData[2]

			println("mapuj:", mapData[1], mapData[2])
		}

	}
}

func (a *App) SaveANPFile() {
	content := string(a.fileBuffer)
	lines := strings.Split(content, "\r\n")

	linesToSave := make([]string, 0)

	for _, line := range lines {
		linesToSave = append(linesToSave, line)

		if strings.HasPrefix(line, "###") {
			break
		}

	}

	for _, run := range a.activeRuns {
		if run.TimestampFrom-time.Now().UnixMilli() < 0 {
			linesToSave = append(linesToSave, run.Content)
		}
	}

	err := os.WriteFile(a.path, []byte(strings.Join(linesToSave, "\r\n")), 0666)

	if err != nil {
		panic("Kek")
	}
}
