package main

import (
	"context"
	"taz/modules/engine"
	"taz/modules/macrocontroller"
	"taz/modules/macroinfo"
	"taz/modules/settingsmanager"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	hook "github.com/robotn/gohook"
)

// App struct
type App struct {
	ctx context.Context
	mc  *macrocontroller.MacroController
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		mc: macrocontroller.NewController(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	go a.Hotkeys()
	a.mc.Engine.Stop()
	a.ctx = ctx
}

func (a *App) GetVersion() string {
	return macroinfo.Version
}

func (a *App) SaveSettings(settings settingsmanager.Settings) {
	settingsmanager.SaveSettings(settings)
}

func (a *App) LoadSettings() settingsmanager.Settings {
	settings, _ := settingsmanager.LoadSettings()
	return settings
}

func (a *App) Start() {
	a.mc.Engine.Start(a.mc.Macro)
}

func (a *App) Stop() {
	a.mc.Engine.Stop()
}

func (a *App) TogglePause() {
	a.mc.Engine.TogglePause()
}

func (a *App) GetEngineState() engine.EngineState {
	return a.mc.Engine.GetState()
}

func (a *App) Hotkeys() {
	hook.Register(hook.KeyUp, []string{"f1"}, func(e hook.Event) {
		a.Start()
		runtime.EventsEmit(a.ctx, "hotkey:start", nil)
	})
	hook.Register(hook.KeyUp, []string{"f2"}, func(e hook.Event) {
		a.TogglePause()
		runtime.EventsEmit(a.ctx, "hotkey:togglepause", nil)
	})
	hook.Register(hook.KeyUp, []string{"f3"}, func(e hook.Event) {
		a.Stop()
		runtime.EventsEmit(a.ctx, "hotkey:stop", nil)
	})

	s := hook.Start()
	<-hook.Process(s)
}
