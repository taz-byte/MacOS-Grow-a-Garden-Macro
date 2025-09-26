package systeminfo

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
#include <CoreGraphics/CoreGraphics.h>
#include <stdio.h>

int isRetinaDisplay() {
    uint32_t displayCount;
    CGGetActiveDisplayList(0, NULL, &displayCount);
    CGDirectDisplayID displays[displayCount];
    CGGetActiveDisplayList(displayCount, displays, &displayCount);

    for (uint32_t i = 0; i < displayCount; i++) {
        CGSize size = CGDisplayScreenSize(displays[i]);
        CGRect bounds = CGDisplayBounds(displays[i]);
        float scale = (float)bounds.size.width / (size.width / 25.4 * 72);
        if (scale > 1.0) {
            return 1;
        }
    }
    return 0;
}
*/
import "C"
import (
	"os/exec"
	"strings"
)

func IsRetinaDisplay() bool {
	isRetina := C.isRetinaDisplay() == 1
	if isRetina {
		println("retina")
	}
	return isRetina
}

func GetOSInfo() string {
	//get the codename
	cmd := `sed -nE '/SOFTWARE LICENSE AGREEMENT FOR/s/([A-Za-z]+ ){5}|\\$//gp' /System/Library/CoreServices/Setup\ Assistant.app/Contents/Resources/en.lproj/OSXSoftwareLicense.rtf`
	codeNameOut, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "unknown"
	}
	codeName := strings.TrimSpace(string(codeNameOut))
	codeName = strings.TrimLeft(codeName, `\f0123456789`)

	//get the version
	versionOut, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "unknown"
	}
	productVersion := strings.TrimSpace(string(versionOut))

	return codeName + " " + productVersion
}

func GetCPUArchitecture() string {
	archOut, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(archOut))
}

func GetDisplayType() string {
	if IsRetinaDisplay() {
		return "Retina Display"
	}
	return "Non-Retina Display"
}
