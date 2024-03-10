package text

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
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
	// нужно для инициализации шейдерной программы этой либы
	// пакет должен совподать с пакетом внутири либы:github.com/go-gl/gl/all-core/gl
	if err := gl.Init(); err != nil {
		log.Printf("LoadFont: %v", err)
	}

	font, err := glfont.LoadFont("./window/text/Roboto-Light.ttf", int32(scale), windowWidth, windowHeight)
	if err != nil {
		log.Printf("LoadFont: %v", err)
	}
	return &Text{font}
}

func (t Text) DrawText(strings []Item) {
	for _, item := range strings {
		if item.Text == nil {
			continue
		}
		err := t.Printf(item.PosX, item.PosY, item.Scale, *item.Text)
		if err != nil {
			fmt.Println(err)
		}
	}
}
