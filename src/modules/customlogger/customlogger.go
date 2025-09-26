package customlogger

import (
	"encoding/json"
	"log"
	"taz/modules/settingsmanager"
	"time"

	discordwebhook "github.com/bensch777/discord-webhook-golang"
)

type Logger struct {
	Settings *settingsmanager.Settings
}

type LogType struct {
	Color int
}

var logTypes = map[string]LogType{
	"regular":  {Color: 0x808080},
	"complete": {Color: 0x66ccff},
	"success":  {Color: 0x2ecc71},
	"warning":  {Color: 0xf1c40f},
	"error":    {Color: 0xe74c3c},
}

func (logger *Logger) SendWebhook(title string, description string, logType string) {
	embed := discordwebhook.Embed{
		Title:       title,
		Description: description,
		Color:       logTypes[logType].Color,
		Timestamp:   nil,
	}
	hook := discordwebhook.Hook{
		Username:   "TAZ GAG Macro",
		Avatar_url: "https://static.wikia.nocookie.net/growagarden/images/a/a4/GoldenPeachProduce.png/revision/latest?cb=20250913035013",
		Content:    "",
		Embeds:     []discordwebhook.Embed{embed},
	}
	payload, err := json.Marshal(hook)
	if err != nil {
		log.Fatal(err)
	}
	discordwebhook.ExecuteWebhook(logger.Settings.DiscordWebhookURL, payload)
}

func (logger *Logger) Log(title string, description string, logType string) {
	currTime := time.Now().Format("[15:04:05]")
	if title != "" {
		title = currTime + " " + title
	} else {
		description = currTime + " " + description
	}
	if logger.Settings.EnableDiscordWebhook && logger.Settings.DiscordWebhookURL != "" {
		logger.SendWebhook(title, description, logType)
	}
}
