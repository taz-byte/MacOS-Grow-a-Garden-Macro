import './style.css';
import {SaveSettings, LoadSettings, GetVersion, Start, TogglePause, Stop, GetEngineState} from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';

function toTitleCase(str) {
    return str.toLowerCase().replace(/\b\w/g, char => char.toUpperCase());
}

//Setup sidebar page switching

document.querySelectorAll(".sidebar-item").forEach(x => {
    x.addEventListener("click", () => {
        document.querySelectorAll(".sidebar-item").forEach(y => {
            y.classList.remove("active");
            document.querySelector(`.tab-pane[data-tab-name="${y.dataset.tabFor}"]`).classList.remove("active");
        })
        x.classList.add("active");
        const tab_name = x.dataset.tabFor;
        document.querySelector(`.tab-pane[data-tab-name="${tab_name}"]`).classList.add("active");
        document.querySelector("#tab-title").innerText = toTitleCase(tab_name);
    })
})

//setup event handlers for inputs
document.querySelectorAll(".toggle").forEach(x => {
    x.addEventListener("click", () => {
        x.classList.toggle("active");
    })
})

//Handle setting inputs


function getInputValue(element){
    if (element.matches('input[type="checkbox"]')) return element.checked;
    if (element.matches('.toggle')) return element.classList.contains("active");
    return element.value;
}

function loadInputElement(element, value){
    if (element.matches('input[type="checkbox"]')){
        element.checked = value;
    }else if (element.matches('.toggle')){
        if (value){
            element.classList.add("active");
        }else{
            element.classList.remove("active");
        }
    }else{
        element.value = value;
    }
}

function saveInputs(){
    //save inputs into the setting file
    let settings = {}
    document.querySelectorAll(".setting-input").forEach(x => {
        const name = x.closest("[data-setting-name]").dataset.settingName;
        const value = getInputValue(x);
        settings[name] = value;
    })
    SaveSettings(settings);
}

function loadInputs(settings){
    for (const [key, value] of Object.entries(settings)) {
        const input_element = document.querySelector(`[data-setting-name="${key}"]`).querySelector(".setting-input");
        loadInputElement(input_element, value);
    }
}

//setup event handlers for inputs
document.querySelectorAll(".toggle.setting-input").forEach(x => {
    x.addEventListener("click", () => {
        saveInputs();
    })
})

document.querySelectorAll('input[type="text"].setting-input').forEach(x => {
    x.addEventListener("change", () => {
        saveInputs();
    })
})

//start, pause and stop handling
function UIStart(){
    document.querySelector(".running-container").classList.add("show");
    document.querySelector("#start-btn").classList.remove("show")
    document.querySelector("#toggle-pause-btn").classList = "stop-btn"
    document.querySelector("#toggle-pause-btn").innerText = "Pause [F2]";
}

function UIStop(){
    document.querySelector(".running-container").classList.remove("show");
    document.querySelector("#start-btn").classList.add("show");
}

function UITogglePause(){
    const toggle_pause_btn = document.querySelector("#toggle-pause-btn")
    if (toggle_pause_btn.classList.contains("stop-btn")){
        toggle_pause_btn.classList = "start-btn";
        toggle_pause_btn.innerText = "Resume [F2]";
    }else{
        toggle_pause_btn.classList = "stop-btn";
        toggle_pause_btn.innerText = "Pause [F2]";
    }
}

document.querySelector("#start-btn").addEventListener("click", () => {
    UIStart();
    Start();
})

document.querySelector("#stop-btn").addEventListener("click", () => {
    UIStop();
    Stop();
})

document.querySelector("#toggle-pause-btn").addEventListener("click", () => {
    UITogglePause();
    TogglePause();
})

EventsOn("hotkey:start", () => {
    UIStart();
})

EventsOn("hotkey:togglepause", () => {
    UITogglePause();
})

EventsOn("hotkey:stop", () => {
    UIStop();
})

//onload handler
window.addEventListener("load", async () => {
    const settings = await LoadSettings();
    loadInputs(settings);

    const version_number = await GetVersion();
    document.querySelector("#version-number").innerText = version_number;
})