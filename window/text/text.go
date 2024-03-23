package text

import (
	"fmt"
	"github.com/flopp/go-findfont"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/nullboundary/glfont"
	"log"
)

type Text struct {
	*glfont.Font
	scale      int32
	activeText []*Item
}

type Item struct {
	Text  *string
	PosX  float32
	PosY  float32
	Scale float32
}

func Init(scale int32, windowWidth, windowHeight int, fontName string) *Text {
	// нужно для инициализации шейдерной программы этой либы
	// пакет должен совподать с пакетом внутири либы:github.com/go-gl/gl/all-core/gl
	if err := gl.Init(); err != nil {
		log.Printf("LoadFont: %v", err)
	}

	fontPath, err := findfont.Find(fontName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %s in '%s'\n", fontName, fontPath)

	font, err := glfont.LoadFont(fontPath, scale, windowWidth, windowHeight)
	if err != nil {
		log.Printf("LoadFont: %v", err)
	}
	return &Text{
		font,
		scale,
		make([]*Item, 0),
	}
}

func (t *Text) AddText(strings []*Item) {
	if t.activeText == nil {
		t.activeText = make([]*Item, 0)
	}
	for i, item := range strings {
		if item.PosY == 0 {
			item.PosY = float32(t.scale*int32(i+1)) * item.Scale
		}
	}
	t.activeText = append(t.activeText, strings...)
	fmt.Printf("added %d strings\n", len(strings))
}

func (t *Text) DrawText() {
	for _, item := range t.activeText {
		if item.Text == nil {
			continue
		}
		err := t.Printf(item.PosX, item.PosY, item.Scale, *item.Text)
		if err != nil {
			fmt.Println(err)
		}
	}
}
