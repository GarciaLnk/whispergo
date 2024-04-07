package main

import (
	"context"
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	runtime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type Model struct {
	Name     string  `json:"name"`
	Size     float32 `json:"size"`
	Download bool    `json:"download"`
	Active   bool    `json:"active"`
}

const builtinModel = "ggml-tiny-q5_1"

var homeDir, _ = os.UserHomeDir()
var appDir = filepath.Join(homeDir, ".whispergo")
var models []Model

//go:embed all:models
var f embed.FS
var model whisper.Model

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	appDirPath := filepath.Join(appDir)
	if _, err := os.Stat(appDirPath); os.IsNotExist(err) {
		err := os.Mkdir(appDirPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	modelsDirPath := filepath.Join(appDir, "models")
	if _, err := os.Stat(modelsDirPath); os.IsNotExist(err) {
		err := os.Mkdir(modelsDirPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	modelListPath := filepath.Join(appDir, "models", "list.json")
	var modelList []byte
	if _, err := os.Stat(modelListPath); os.IsNotExist(err) {
		modelList, err = f.ReadFile("models/list.json")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(modelListPath, modelList, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		modelList, err = os.ReadFile(modelListPath)
		if err != nil {
			panic(err)
		}
	}

	err := json.Unmarshal(modelList, &models)
	if err != nil {
		panic(err)
	}

	var activeModel = ""
	for i, model := range models {
		modelFilePath := filepath.Join(modelsDirPath, model.Name+".bin")
		if _, err := os.Stat(modelFilePath); err == nil {
			models[i].Download = true
		} else {
			models[i].Download = false
			if model.Active {
				models[i].Active = false
				models[0].Active = true
			}
		}
	}

	models[0].Download = true
	for _, model := range models {
		if model.Active {
			activeModel = model.Name
		}
	}

	modelList, err = json.Marshal(models)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(modelListPath, modelList, 0644)
	if err != nil {
		panic(err)
	}

	modelFile, err := getModelFile(activeModel, modelsDirPath)
	if err != nil {
		panic(err)
	}

	model, err = whisper.New(modelFile.Name())
	if err != nil {
		panic(err)
	}
	defer modelFile.Close()
}

func (a *App) shutdown(ctx context.Context) {
	model.Close()
}

func (a *App) LoadFile() string {
	var dialogOptions runtime.OpenDialogOptions
	dialogOptions.Title = "Choose audio file"
	dialogOptions.DefaultDirectory, _ = os.UserHomeDir()
	dialogOptions.ShowHiddenFiles = false
	dialogOptions.Filters = []runtime.FileFilter{
		{
			DisplayName: "Audio Files (*.wav)",
			Pattern:     "*.wav",
		},
	}

	filepath, err := runtime.OpenFileDialog(a.ctx, dialogOptions)
	if err != nil {
		return "error opening file"
	}

	return filepath
}

func (a *App) ProcessFile(filepath string) string {
	// Wait for the model to be loaded
	for model == nil {
		time.Sleep(100 * time.Millisecond)
	}

	// Process the file
	var cb string
	Process(model, filepath, &cb)

	return cb
}

func (a *App) GetModels() []Model {
	return models
}

func (a *App) DownloadModel(modelName string) bool {
	modelUrl, err := URLForModel(modelName)
	if err != nil {
		return false
	}

	Download(a.ctx, os.Stdout, modelUrl, filepath.Join(appDir, "models"))

	// update the list of models
	for i, m := range models {
		if m.Name == modelName {
			models[i].Download = true
		}
	}
	modelList, _ := json.Marshal(models)
	os.WriteFile(filepath.Join(appDir, "models", "list.json"), modelList, 0644)

	return true
}

func (a *App) SetActiveModel(modelName string) {
	// update the list of models
	for i, m := range models {
		if m.Name == modelName {
			models[i].Active = true
		} else {
			models[i].Active = false
		}
	}
	modelList, _ := json.Marshal(models)
	os.WriteFile(filepath.Join(appDir, "models", "list.json"), modelList, 0644)

	// load the model
	modelFile, err := getModelFile(modelName, filepath.Join(appDir, "models"))
	if err != nil {
		panic(err)
	}

	model, err = whisper.New(modelFile.Name())
	if err != nil {
		panic(err)
	}
	defer modelFile.Close()
}

func getModelFile(activeModel, modelsDirPath string) (*os.File, error) {
	var modelFile *os.File
	if activeModel == builtinModel {
		modelData, err := f.ReadFile("models/" + builtinModel + ".bin")
		if err != nil {
			return nil, err
		}

		modelFile, err = os.CreateTemp("", builtinModel+".bin")
		if err != nil {
			return nil, err
		}

		_, err = modelFile.Write(modelData)
		if err != nil {
			return nil, err
		}

	} else {
		modelFilePath := filepath.Join(modelsDirPath, activeModel+".bin")
		modelFile, _ = os.Open(modelFilePath)
	}

	return modelFile, nil
}
