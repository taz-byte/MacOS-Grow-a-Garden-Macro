package windowmanager

// #cgo LDFLAGS: -framework CoreGraphics
// #cgo LDFLAGS: -framework CoreFoundation
// #include <CoreGraphics/CoreGraphics.h>
// #include <CoreFoundation/CoreFoundation.h>
// #cgo LDFLAGS: -framework ApplicationServices
// #include <ApplicationServices/ApplicationServices.h>
import (
	// #cgo LDFLAGS: -framework CoreGraphics
	// #cgo LDFLAGS: -framework CoreFoundation
	// #include <CoreGraphics/CoreGraphics.h>
	// #include <CoreFoundation/CoreFoundation.h>
	// #cgo LDFLAGS: -framework ApplicationServices
	// #include <ApplicationServices/ApplicationServices.h>
	"C"

	"strings"
	"unsafe"

	"github.com/go-vgo/robotgo"
)

func cfStringToGoString(cfstr C.CFStringRef) string {
	if cfstr == 0 {
		return ""
	}
	length := C.CFStringGetLength(cfstr)
	if length == 0 {
		return ""
	}
	// allocate buffer
	buf := make([]C.char, (length*4)+1) // UTF8 may expand
	if C.CFStringGetCString(cfstr, &buf[0], C.CFIndex(len(buf)), C.kCFStringEncodingUTF8) == 0 {
		return ""
	}
	return C.GoString(&buf[0])
}

func goStringToCFString(gostr string) C.CFStringRef {
	cstr := C.CString(gostr)
	defer C.free(unsafe.Pointer(cstr))
	return C.CFStringCreateWithCString(0, cstr, C.kCFStringEncodingUTF8)
}

func getIntFromDict(dict C.CFDictionaryRef, key string) (int, bool) {
	cfKey := goStringToCFString(key)
	defer C.CFRelease(C.CFTypeRef(cfKey))

	var value unsafe.Pointer
	if C.CFDictionaryGetValueIfPresent(dict, unsafe.Pointer(cfKey), &value) == 0 {
		return 0, false
	}

	var num C.CFNumberRef = (C.CFNumberRef)(value)
	var out C.int
	if C.CFNumberGetValue(num, C.kCFNumberIntType, unsafe.Pointer(&out)) == 0 {
		return 0, false
	}
	return int(out), true
}

func GetRobloxWindowBounds() (int, int, int, int) {
	// go hotkeys()
	// engine.setState(Idle)

	// // keep program alive
	// select {}
	windows := C.CGWindowListCopyWindowInfo(C.kCGWindowListExcludeDesktopElements|C.kCGWindowListOptionOnScreenOnly, C.kCGNullWindowID)
	count := C.CFArrayGetCount(windows)
	for i := C.CFIndex(0); i < count; i++ {
		elem := C.CFArrayGetValueAtIndex(windows, i)
		dict := (C.CFDictionaryRef)(unsafe.Pointer(elem))

		var titleVal unsafe.Pointer
		if C.CFDictionaryGetValueIfPresent(dict, unsafe.Pointer(C.kCGWindowName), &titleVal) != 0 {
			title := cfStringToGoString((C.CFStringRef)(titleVal))
			if title != "" && strings.Contains(title, "Roblox") {
				// get bounds dictionary
				var boundsVal unsafe.Pointer
				key := goStringToCFString("kCGWindowBounds")
				defer C.CFRelease(C.CFTypeRef(key))
				if C.CFDictionaryGetValueIfPresent(dict, unsafe.Pointer(key), &boundsVal) != 0 {
					bounds := (C.CFDictionaryRef)(boundsVal)
					x, okX := getIntFromDict(bounds, "X")
					y, okY := getIntFromDict(bounds, "Y")
					w, okW := getIntFromDict(bounds, "Width")
					h, okH := getIntFromDict(bounds, "Height")

					if okX && okY && okW && okH {
						return x, y + 30, w, h - 30
					}
				}
			}
		}
	}

	return 0, 0, 0, 0

}

func GetRobloxPID() (int, bool) {
	fpid, err := robotgo.FindIds("Roblox")
	if err == nil {
		if len(fpid) > 0 {
			return fpid[0], true
		}
	}
	return 0, false
}

func MinimiseWindow(pid int) {
	robotgo.MinWindow(pid)
}

func ActivateWindow(pid int) {
	robotgo.ActivePid(pid)
}

func SetFullscreen(pid int, fullscreen bool) {
	appElement := C.AXUIElementCreateApplication(C.pid_t(pid))

	var windowCFType C.CFTypeRef
	cfWindowAttribute := goStringToCFString("AXMainWindow")
	defer C.CFRelease(C.CFTypeRef(cfWindowAttribute))
	C.AXUIElementCopyAttributeValue(appElement, cfWindowAttribute, &windowCFType)
	windowElement := C.AXUIElementRef(windowCFType)

	cfFullscreenAttribute := goStringToCFString("AXFullScreen")
	defer C.CFRelease(C.CFTypeRef(cfFullscreenAttribute))

	var cfBool C.CFBooleanRef
	if fullscreen {
		cfBool = C.kCFBooleanTrue
	} else {
		cfBool = C.kCFBooleanFalse
	}
	C.AXUIElementSetAttributeValue(windowElement, cfFullscreenAttribute, C.CFTypeRef(cfBool))
}

func SetWindowBounds(pid int, x int, y int, width int, height int) {
	appElement := C.AXUIElementCreateApplication(C.pid_t(pid))

	var windowCFType C.CFTypeRef
	cfWindowAttribute := goStringToCFString("AXMainWindow")
	defer C.CFRelease(C.CFTypeRef(cfWindowAttribute))
	C.AXUIElementCopyAttributeValue(appElement, cfWindowAttribute, &windowCFType)
	windowElement := C.AXUIElementRef(windowCFType)

	var pos C.CGPoint
	pos.x = C.double(x)
	pos.y = C.double(y)
	posValue := C.CFTypeRef(unsafe.Pointer(C.AXValueCreate(C.kAXValueCGPointType, unsafe.Pointer(&pos))))
	defer C.CFRelease(posValue)

	var size C.CGSize
	size.width = C.double(width)
	size.height = C.double(height)
	sizeValue := C.CFTypeRef(unsafe.Pointer(C.AXValueCreate(C.kAXValueCGSizeType, unsafe.Pointer(&size))))
	defer C.CFRelease(sizeValue)

	cfWindowPosAttribute := goStringToCFString("AXPosition")
	defer C.CFRelease(C.CFTypeRef(cfWindowPosAttribute))
	C.AXUIElementSetAttributeValue(windowElement, cfWindowPosAttribute, posValue)

	cfWindowSizeAttribute := goStringToCFString("AXSize")
	defer C.CFRelease(C.CFTypeRef(cfWindowSizeAttribute))
	C.AXUIElementSetAttributeValue(windowElement, cfWindowSizeAttribute, sizeValue)

}

func MaximizeWindow(pid int) {
	screenW, screenH := robotgo.GetScreenSize()
	SetWindowBounds(pid, 0, 0, screenW, screenH)
}
