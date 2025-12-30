package config

import (
	"encoding/json"
	"log"
	"os"
)

const (
	settingsFile = "settings.json"
)

type ConfigNovel struct {
	LastChapterUrl string `json:"last_chapter_url"`
	BatchSize      int    `json:"batch_size"`
}

type ConfigData struct {
	HttpHeaders struct {
		UserAgent string            `json:"User-Agent"`
		Cookies   map[string]string `json:"cookies"`
	} `json:"http_headers"`
	Novels map[string]ConfigNovel `json:"novels"`
}

var appConfig ConfigData

func init() {
	Init()
}

func Init() {
	content, err := os.ReadFile(AssetPath(settingsFile))

	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(content, &appConfig)
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

	// log.Println(novel, appConfig.Novels)

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

	check(os.WriteFile(AssetPath(settingsFile), jsonData, 0644))
}

// func AssetPath(file string) string {
// 	return filepath.Join(WorkingDir, file)
// }

func createNovel() ConfigNovel {
	return ConfigNovel{BatchSize: 10}
}
