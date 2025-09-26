package macrocontroller

import (
	"os"
	"strings"
	"time"

	"taz/modules/customlogger"
	"taz/modules/engine"
	"taz/modules/imagesearch"
	"taz/modules/macroinfo"
	"taz/modules/settingsmanager"
	"taz/modules/systeminfo"
	"taz/modules/timingsmanager"
	"taz/modules/windowmanager"

	"github.com/go-vgo/robotgo"
	"github.com/micmonay/keybd_event"
)

type MacroController struct {
	Engine      *engine.ScriptEngine
	Logger      *customlogger.Logger
	windowX     int
	windowY     int
	windowW     int
	windowH     int
	Keyboard    keybd_event.KeyBonding
	ImageSearch *imagesearch.ImageSearch
}

func NewController() *MacroController {
	mc := &MacroController{}
	mc.Engine = engine.NewEngine()
	mc.Engine.SetOnResume(func() {
		mc.RobloxWindowSetup()
		mc.Logger.Log("Macro Resumed", "", "success")
	})
	mc.Engine.SetOnStop(func() {
		pid := os.Getpid()
		windowmanager.ActivateWindow(pid)
		mc.ReleaseInputs()
		mc.Logger.Log("Macro Stopped", "", "error")
	})
	mc.Engine.SetOnPause(func() {
		mc.ReleaseInputs()
		mc.Logger.Log("Macro Paused", "", "warning")
	})
	mc.Keyboard, _ = keybd_event.NewKeyBonding()
	mc.ImageSearch = imagesearch.NewImageSearch(systeminfo.IsRetinaDisplay())
	mc.Logger = &customlogger.Logger{}
	return mc
}

func (mc *MacroController) PressKey(key int, duration time.Duration) {
	mc.Keyboard.SetKeys(key)
	mc.Keyboard.Press()
	mc.Engine.Sleep(duration)
	mc.Keyboard.Release()
}

func (mc *MacroController) ReleaseInputs() {
	mc.Keyboard.SetKeys(keybd_event.VK_W, keybd_event.VK_A, keybd_event.VK_S, keybd_event.VK_D, keybd_event.VK_E, keybd_event.VK_SPACE)
	mc.Keyboard.Release()
	robotgo.MouseUp()
}

func (mc *MacroController) RobloxWindowSetup() {
	robloxPid, _ := windowmanager.GetRobloxPID()
	mc.Engine.RunFuncNoReturn(func() { windowmanager.ActivateWindow(robloxPid) })
	mc.Engine.RunFuncNoReturn(func() { windowmanager.SetFullscreen(robloxPid, false) })
	mc.Engine.RunFuncNoReturn(func() { windowmanager.SetWindowBounds(robloxPid, 0, 0, 9999, 0) })
	mc.Engine.Sleep(1 * time.Second)
	mc.windowX, mc.windowY, mc.windowW, mc.windowH = windowmanager.GetRobloxWindowBounds()
}

func (mc *MacroController) ClickElement(filePath string, maxAttempts int, x int, y int, w int, h int) {
	for range maxAttempts {
		fx, fy := mc.ImageSearch.FindImageFileOnScreen(filePath, x, y, w, h, 0.2, true, false, true)
		robotgo.Move(fx, fy)
		if fx >= 0 && fy >= 0 {
			mc.Engine.RunFuncNoReturn(func() { robotgo.Move(fx, fy) })
			for range 2 {
				mc.Engine.RunFuncNoReturn(func() { robotgo.Click() })
			}
			break
		}
		mc.Engine.Sleep(300 * time.Millisecond)
	}
}

func (mc *MacroController) BuyFromShop(items []settingsmanager.BuyItemSettings) []string {
	var purchasedItems = []string{}
	//buy from bottom up
	mc.Engine.RunFuncNoReturn(func() { robotgo.Move(mc.windowX+mc.windowW/2, mc.windowY+mc.windowH/2) })
	for range 20 {
		mc.Engine.RunFuncNoReturn(func() { robotgo.ScrollDir(999, "down") })
	}
	mc.Engine.Sleep(500 * time.Millisecond)

	var mouseY int
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]

		//bottom 2 items are edge cases, mouse needs to specifically click on them
		//otherwise, click the top to open the item
		if i == len(items)-1 {
			mouseY = mc.windowY + int(float64(mc.windowH)*0.73)
		} else if i == len(items)-2 {
			mouseY = mc.windowY + int(float64(mc.windowH)*0.51)
		} else {
			mouseY = mc.windowY + int(float64(mc.windowH)*0.33)
		}
		mc.Engine.RunFuncNoReturn(func() { robotgo.Move(mc.windowX+mc.windowW/2, mouseY) })
		mc.Engine.RunFuncNoReturn(func() { robotgo.Click() })

		var sleepTime time.Duration
		if item.Enabled {
			sleepTime = 600
		} else {
			sleepTime = 200
		}
		//scroll up near the top to reach the remaining 2 items
		if i <= 2 {
			mc.Engine.Sleep(400 * time.Millisecond)
			robotgo.ScrollDir(80, "up")
			sleepTime = max(sleepTime-400, 100)
		}
		mc.Engine.Sleep(sleepTime * time.Millisecond)

		//check if item can be purchased
		if !item.Enabled {
			continue
		}
		fx, fy := mc.ImageSearch.FindImageFileOnScreen("images/buy_btn.png", mc.windowX, mc.windowY, mc.windowW, mc.windowH, 0.1, false, false, true)
		if fx >= 0 && fy >= 0 {
			mc.Engine.RunFuncNoReturn(func() { robotgo.Move(fx+20, fy+20) })
			for range 20 {
				mc.Engine.RunFuncNoReturn(func() { robotgo.Click() })
				mc.Engine.Sleep(30 * time.Millisecond)
			}
			purchasedItems = append(purchasedItems, item.Name)
		}
	}
	mc.ClickElement("images/close_btn.png", 5, mc.windowX+mc.windowW/2, mc.windowY, mc.windowW/2, mc.windowH/2)
	mc.Engine.Sleep(1000 * time.Millisecond)
	return purchasedItems
}

func (mc *MacroController) Macro() {
	settings, _ := settingsmanager.LoadSettings()
	mc.Logger.Settings = &settings
	mc.Logger.SendWebhook("Taz Macro Started",
		"**Macro Version:** "+macroinfo.Version+
			"\n\n**System Info:**"+
			"\nOS: "+systeminfo.GetOSInfo()+
			"\nCPU Architecture: "+systeminfo.GetCPUArchitecture()+
			"\nDisplay: "+systeminfo.GetDisplayType(),
		"success")
	mc.RobloxWindowSetup()
	mc.ClickElement("images/garden_btn.png", 3, mc.windowX, mc.windowY, mc.windowW, mc.windowH/2)

	for {
		now := time.Now()
		currMinute := now.Minute()
		currSecond := now.Second()
		seeds := settings.GetSeedsToBuy()
		if !settingsmanager.AllItemsToBuyAreDisabled(seeds) && ((currMinute%5 == 0 && currSecond > 10 && !timingsmanager.OnCooldown("SeedShop", 40*time.Second)) || !timingsmanager.OnCooldown("SeedShop", 5*time.Minute)) {
			mc.Logger.Log("", "Going to Seed Shop", "regular")
			mc.ClickElement("images/seeds_btn.png", 5, mc.windowX, mc.windowY, mc.windowW, mc.windowH/2)
			mc.Engine.Sleep(1 * time.Second)
			mc.Engine.RunFuncNoReturn(func() { mc.PressKey(keybd_event.VK_E, 1*time.Second) })
			mc.Engine.Sleep(1500 * time.Millisecond)
			purchasedSeeds := mc.BuyFromShop(seeds)
			if len(purchasedSeeds) > 0 {
				mc.Logger.Log("", "Bought Seeds:\n\n"+strings.Join(purchasedSeeds, "\n"), "complete")
			} else {
				mc.Logger.Log("", "Did not buy any seeds", "regular")
			}
			timingsmanager.UpdateObjectiveTime("SeedShop")
			mc.ClickElement("images/garden_btn.png", 3, mc.windowX, mc.windowY, mc.windowW, mc.windowH/2)
		}
		mc.Engine.Sleep(5 * time.Second)
		mc.ClickElement("images/close_btn.png", 1, mc.windowX+mc.windowW/2, mc.windowY, mc.windowW/2, mc.windowH/2)
		robotgo.Click()
	}
}
