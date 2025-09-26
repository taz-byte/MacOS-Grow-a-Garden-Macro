package displaydata

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

func IsRetinaDisplay() bool {
	isRetina := C.isRetinaDisplay() == 1
	if isRetina {
		println("retina")
	}
	println("nah")
	return isRetina
}
