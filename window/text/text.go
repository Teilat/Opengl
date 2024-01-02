package text

import (
	"fmt"
	"github.com/nullboundary/glfont"
	"log"
)

type Text struct {
	*glfont.Font
}

type Item struct {
	Text  *string
	PosX  float32
	PosY  float32
	Scale float32
}

func Init(scale, windowWidth, windowHeight int) *Text {
	font, err := glfont.LoadFont("./window/text/Roboto-Light.ttf", int32(scale), windowWidth, windowHeight)
	if err != nil {
		log.Printf("LoadFont: %v", err)
	}
	return &Text{font}
}

func (t Text) DrawText(strings []Item) {
	// TODO передавать массив структур со строуками, позиуиями, цветами и тд...
	for _, item := range strings {
		err := t.Printf(item.PosX, item.PosY, item.Scale, *item.Text)
		if err != nil {
			fmt.Println(err)
		}
	}
}
