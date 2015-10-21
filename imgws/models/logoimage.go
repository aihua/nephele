package models

import (
	"fmt"
)

type LogoImage struct {
	size     int
	name     string
	logoPath string
	infoPath string
	filePath string
}

func NewLogoImage(size int, path string, name string, filename string) *LogoImage {
	l := &LogoImage{
		size,
		name,
		fmt.Sprintf("_logo_%d", size),
		fmt.Sprintf("_info_%d", size),
		filename,
	}
	return l
}

func (this LogoImage) Load() {
}
