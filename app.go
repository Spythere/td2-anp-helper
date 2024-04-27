package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	path    string
	version string

	activeData     *ActiveData
	routeMappings  map[string]string
	currentRuns    map[int]string
	runDefinitions []RunDefinition
	activeRuns     []ActiveRun
	checkpoints    []string
	fileBuffer     []byte

	ticker *time.Ticker
}

var TICKER_SECONDS = time.Second * 10

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		version:        "0.0.0",
		routeMappings:  make(map[string]string),
		runDefinitions: make([]RunDefinition, 0),
		activeRuns:     make([]ActiveRun, 0),
		currentRuns:    make(map[int]string),
		fileBuffer:     make([]byte, 0),
		ticker:         time.NewTicker(TICKER_SECONDS),
	}
}

func (a *App) ShowActiveRuns() {
	fmt.Printf("\nChosen file: %s\n", a.path)
	fmt.Printf("\n============= Active runs (%s): =============\n", time.Now().Local())

	if len(a.activeRuns) == 0 {
		println("No active runs.")
	}

	for _, run := range a.activeRuns {
		if run.TimestampFrom > 0 {
			active := "not active"
			if run.TimestampFrom-time.Now().UnixMilli() < 0 {
				active = "active"
			}

			fmt.Printf("%s (%s, %s)\n", run.Content, time.UnixMilli(run.TimestampFrom), active)
		} else {
			fmt.Printf("%s\n", run.Content)

		}
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.ticker.Stop()

	go func() {
		for {
			select {
			case <-a.ticker.C:
				a.refreshANP()
			}
		}
	}()
}

func (a *App) refreshANP() {
	// mocks/Otwocko.anp
	// C:\\Users\\spyth\\Documents\\TTSK\\TrainDriver2\\SavedStations\\ANP\\Głęboszów.anp

	a.processData()
	a.SaveANPFile()
	a.ShowActiveRuns()

}

func (a *App) ResetANP() {
	a.path = ""
	a.ticker.Stop()
	a.activeData = nil
	a.runDefinitions = make([]RunDefinition, 0)
	a.checkpoints = make([]string, 0)
	a.currentRuns = make(map[int]string)
	a.routeMappings = map[string]string{}
}

func (a *App) GetActiveRuns() string {
	b, err := json.Marshal(a.activeRuns)

	if err != nil {
		return ""
	}

	return string(b)
}

func (a *App) GetFilePath() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Wybierz plik ANP do aktualizacji",
		Filters: []runtime.FileFilter{{Pattern: "*.anp"}},
	})

	if err != nil {
		fmt.Println(err)
		return ""
	}

	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Uwaga",
		Message: "Uwaga! Wybrany plik ANP zostanie nadpisany!",
		Type:    runtime.WarningDialog,
	})

	a.ProcessANPFile(path)
	a.ticker.Reset(TICKER_SECONDS)
	a.refreshANP()

	return path
}
