package text

import (
	"github.com/nullboundary/glfont"
	"log"
)

func Init(windowWidth, windowHeight int) *glfont.Font {
	font, err := glfont.LoadFont("./window/text/Roboto-Light.ttf", int32(52), windowWidth, windowHeight)
	if err != nil {
		log.Printf("LoadFont: %v", err)
	}
	return font
}
