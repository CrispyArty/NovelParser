package config

import (
	"encoding/json"
	"log"
	"os"
)

const settingsFile = "settings.json"

type ConfigNovel struct {
	LastChapterUrl string `json:"last_chapter_url"`
	BatchSize      int    `json:"batch_size"`
}

type ConfigData struct {
	Novels map[string]ConfigNovel `json:"novels"`
}

var appConfig ConfigData

func Init() {
	content, err := os.ReadFile(settingsFile)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(content, &appConfig)
}

func createNovel() ConfigNovel {
	return ConfigNovel{BatchSize: 10}
}

func Novel(name string) ConfigNovel {
	return appConfig.Novels[name]
}

func Config() ConfigData {
	return appConfig
}

func UpdateLastChapter(novelName string, url string) {
	novel, ok := appConfig.Novels[novelName]
	if !ok {
		novel = createNovel()
	}

	novel.LastChapterUrl = url

	appConfig.Novels[novelName] = novel
}

func Save() {
	check := func(err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}

	jsonData, err := json.MarshalIndent(appConfig, "", "  ")
	check(err)
	check(os.WriteFile(settingsFile, jsonData, 0644))
}
