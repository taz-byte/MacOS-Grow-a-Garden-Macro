package imagesearch

import (
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/bitmap"
)

// FindImageFileOnScreen searches for an image file (needle) inside a captured screen region.
//
// Parameters:
//   - path: path to the image file to search for.
//   - x, y: top-left coordinates of the screen region to search.
//   - width, height: dimensions of the screen region.
//   - tolerance: similarity threshold (0.0–1.0). Larger values allow looser matches.
//   - returnCenter: whether to return the center of the match instead of the top-left corner.
//   - screenCoords: if true, returns absolute screen coordinates. If false, scales down coords by 2 (for Retina).
//   - includeStartingPoint: whether to include the search region’s starting point in returned coords.
//
// Returns:
//   - fx, fy: the found coordinates, or (-1, -1) if no match was found.

type ImageSearch struct {
	RetinaDisplay bool
}

func NewImageSearch(retinaDisplay bool) *ImageSearch {
	is := &ImageSearch{}
	is.RetinaDisplay = retinaDisplay
	return is
}

func (*ImageSearch) FindImageFileOnScreen(path string, x int, y int, width int, height int, tolerance float64, returnCenter bool, screenCoords bool, includeStartingPoint bool) (int, int) {
	needle := bitmap.Open(path)
	defer robotgo.FreeBitmap(needle)

	screen := robotgo.CaptureScreen(x, y, width, height)
	defer robotgo.FreeBitmap(screen)

	// img := robotgo.ToImage(screen)
	// imgo.Save("test.png", img)

	fx, fy := bitmap.Find(needle, screen, tolerance)

	if fx < 0 || fy < 0 {
		return -1, -1
	}

	if returnCenter {
		gbit := robotgo.ToBitmap(needle)
		fx += gbit.Width / 2
		fy += gbit.Height / 2
	}

	if !screenCoords {
		fx /= 2
		fy /= 2
	}

	if includeStartingPoint {
		fx += x
		fy += y
	}

	return fx, fy
}
