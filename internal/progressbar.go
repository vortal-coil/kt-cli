package internal

import "github.com/schollz/progressbar/v3"

// NewProgressBar @todo integrate with download and upload
func NewProgressBar(max int64) *progressbar.ProgressBar {
	bar := progressbar.Default(max)
	return bar
}
