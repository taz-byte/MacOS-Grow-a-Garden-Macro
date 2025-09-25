package settingsmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	BuyCarrot          bool `json:"buy_carrot"`
	BuyStrawberry      bool `json:"buy_strawberry"`
	BuyBlueberry       bool `json:"buy_blueberry"`
	BuyOrangeTulip     bool `json:"buy_orange_tulip"`
	BuyTomato          bool `json:"buy_tomato"`
	BuyCorn            bool `json:"buy_corn"`
	BuyDaffodil        bool `json:"buy_daffodil"`
	BuyWatermelon      bool `json:"buy_watermelon"`
	BuyPumpkin         bool `json:"buy_pumpkin"`
	BuyApple           bool `json:"buy_apple"`
	BuyBamboo          bool `json:"buy_bamboo"`
	BuyCoconut         bool `json:"buy_coconut"`
	BuyCactus          bool `json:"buy_cactus"`
	BuyDragonFruit     bool `json:"buy_dragon_fruit"`
	BuyMango           bool `json:"buy_mango"`
	BuyGrape           bool `json:"buy_grape"`
	BuyMushroom        bool `json:"buy_mushroom"`
	BuyPepper          bool `json:"buy_pepper"`
	BuyCacao           bool `json:"buy_cacao"`
	BuyBeanstalk       bool `json:"buy_beanstalk"`
	BuyEmberLily       bool `json:"buy_ember_lily"`
	BuySugarApple      bool `json:"buy_sugar_apple"`
	BuyBurningBud      bool `json:"buy_burning_bud"`
	BuyGiantPinecone   bool `json:"buy_giant_pinecone"`
	BuyElderStrawberry bool `json:"buy_elder_strawberry"`
	BuyRomanesco       bool `json:"buy_romanesco"`

	EnableDiscordWebhook bool   `json:"enable_discord_webhook"`
	DiscordWebhookURL    string `json:"discord_webhook_url"`
}

type BuyItemSettings struct {
	Name    string
	Enabled bool
}

func (s *Settings) GetSeedsToBuy() []BuyItemSettings {
	return []BuyItemSettings{
		{"Carrot", s.BuyCarrot},
		{"Strawberry", s.BuyStrawberry},
		{"Blueberry", s.BuyBlueberry},
		{"Orange Tulip", s.BuyOrangeTulip},
		{"Tomato", s.BuyTomato},
		{"Corn", s.BuyCorn},
		{"Daffodil", s.BuyDaffodil},
		{"Watermelon", s.BuyWatermelon},
		{"Pumpkin", s.BuyPumpkin},
		{"Apple", s.BuyApple},
		{"Bamboo", s.BuyBamboo},
		{"Coconut", s.BuyCoconut},
		{"Cactus", s.BuyCactus},
		{"Dragon Fruit", s.BuyDragonFruit},
		{"Mango", s.BuyMango},
		{"Grape", s.BuyGrape},
		{"Mushroom", s.BuyMushroom},
		{"Pepper", s.BuyPepper},
		{"Cacao", s.BuyCacao},
		{"Beanstalk", s.BuyBeanstalk},
		{"Ember Lily", s.BuyEmberLily},
		{"Sugar Apple", s.BuySugarApple},
		{"Burning Bud", s.BuyBurningBud},
		{"Giant Pinecone", s.BuyGiantPinecone},
		{"Elder Strawberry", s.BuyElderStrawberry},
		{"Romanesco", s.BuyRomanesco},
	}
}

func getSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Library", "Application Support", "Taz GAG Macro")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "settings.json"), nil
}

func SaveSettings(settings Settings) error {
	//save settings into a json file
	data, err := json.MarshalIndent(settings, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	path, err := getSettingsPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadSettings() (Settings, error) {
	var settings Settings

	path, err := getSettingsPath()
	if err != nil {
		return settings, fmt.Errorf("failed to get settings path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// If the file does not exist, return default settings without error
		if os.IsNotExist(err) {
			return settings, nil
		}
		return settings, fmt.Errorf("failed to read settings.json: %w", err)
	}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return settings, fmt.Errorf("failed to parse settings.json: %w", err)
	}

	return settings, nil
}
