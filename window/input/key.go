package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type keyAttr struct {
	Release bool
	Press   bool
	Repeat  bool
}

type Axis int

const (
	Horizontal Axis = iota // A,D,ArrowLeft,ArrowRight
	Vertical               // W,S,ArrowUp,ArrowDown
	MouseX
	MouseY
)

var (
	keys = map[glfw.Key]keyAttr{
		glfw.KeyUnknown:      {},
		glfw.KeySpace:        {},
		glfw.KeyApostrophe:   {},
		glfw.KeyComma:        {},
		glfw.KeyMinus:        {},
		glfw.KeyPeriod:       {},
		glfw.KeySlash:        {},
		glfw.Key0:            {},
		glfw.Key1:            {},
		glfw.Key2:            {},
		glfw.Key3:            {},
		glfw.Key4:            {},
		glfw.Key5:            {},
		glfw.Key6:            {},
		glfw.Key7:            {},
		glfw.Key8:            {},
		glfw.Key9:            {},
		glfw.KeySemicolon:    {},
		glfw.KeyEqual:        {},
		glfw.KeyA:            {},
		glfw.KeyB:            {},
		glfw.KeyC:            {},
		glfw.KeyD:            {},
		glfw.KeyE:            {},
		glfw.KeyF:            {},
		glfw.KeyG:            {},
		glfw.KeyH:            {},
		glfw.KeyI:            {},
		glfw.KeyJ:            {},
		glfw.KeyK:            {},
		glfw.KeyL:            {},
		glfw.KeyM:            {},
		glfw.KeyN:            {},
		glfw.KeyO:            {},
		glfw.KeyP:            {},
		glfw.KeyQ:            {},
		glfw.KeyR:            {},
		glfw.KeyS:            {},
		glfw.KeyT:            {},
		glfw.KeyU:            {},
		glfw.KeyV:            {},
		glfw.KeyW:            {},
		glfw.KeyX:            {},
		glfw.KeyY:            {},
		glfw.KeyZ:            {},
		glfw.KeyLeftBracket:  {},
		glfw.KeyBackslash:    {},
		glfw.KeyRightBracket: {},
		glfw.KeyGraveAccent:  {},
		glfw.KeyWorld1:       {},
		glfw.KeyWorld2:       {},
		glfw.KeyEscape:       {},
		glfw.KeyEnter:        {},
		glfw.KeyTab:          {},
		glfw.KeyBackspace:    {},
		glfw.KeyInsert:       {},
		glfw.KeyDelete:       {},
		glfw.KeyRight:        {},
		glfw.KeyLeft:         {},
		glfw.KeyDown:         {},
		glfw.KeyUp:           {},
		glfw.KeyPageUp:       {},
		glfw.KeyPageDown:     {},
		glfw.KeyHome:         {},
		glfw.KeyEnd:          {},
		glfw.KeyCapsLock:     {},
		glfw.KeyScrollLock:   {},
		glfw.KeyNumLock:      {},
		glfw.KeyPrintScreen:  {},
		glfw.KeyPause:        {},
		glfw.KeyF1:           {},
		glfw.KeyF2:           {},
		glfw.KeyF3:           {},
		glfw.KeyF4:           {},
		glfw.KeyF5:           {},
		glfw.KeyF6:           {},
		glfw.KeyF7:           {},
		glfw.KeyF8:           {},
		glfw.KeyF9:           {},
		glfw.KeyF10:          {},
		glfw.KeyF11:          {},
		glfw.KeyF12:          {},
		glfw.KeyF13:          {},
		glfw.KeyF14:          {},
		glfw.KeyF15:          {},
		glfw.KeyF16:          {},
		glfw.KeyF17:          {},
		glfw.KeyF18:          {},
		glfw.KeyF19:          {},
		glfw.KeyF20:          {},
		glfw.KeyF21:          {},
		glfw.KeyF22:          {},
		glfw.KeyF23:          {},
		glfw.KeyF24:          {},
		glfw.KeyF25:          {},
		glfw.KeyKP0:          {},
		glfw.KeyKP1:          {},
		glfw.KeyKP2:          {},
		glfw.KeyKP3:          {},
		glfw.KeyKP4:          {},
		glfw.KeyKP5:          {},
		glfw.KeyKP6:          {},
		glfw.KeyKP7:          {},
		glfw.KeyKP8:          {},
		glfw.KeyKP9:          {},
		glfw.KeyKPDecimal:    {},
		glfw.KeyKPDivide:     {},
		glfw.KeyKPMultiply:   {},
		glfw.KeyKPSubtract:   {},
		glfw.KeyKPAdd:        {},
		glfw.KeyKPEnter:      {},
		glfw.KeyKPEqual:      {},
		glfw.KeyLeftShift:    {},
		glfw.KeyLeftControl:  {},
		glfw.KeyLeftAlt:      {},
		glfw.KeyLeftSuper:    {},
		glfw.KeyRightShift:   {},
		glfw.KeyRightControl: {},
		glfw.KeyRightAlt:     {},
		glfw.KeyRightSuper:   {},
		glfw.KeyMenu:         {},
	}
	mKey glfw.ModifierKey
)

func KeyCallback(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, modKey glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
	mKey = modKey
	keys[key] = keyAttr{action == 0, action == 1, action == 2}
	// fmt.Printf("%s:\"%s\"%+v, %d\n", time.Now().Format("04:05.000"), string(rune(key)), keys[key], modKey)
}

func GetKey(key glfw.Key) bool {
	return keys[key].Repeat || keys[key].Press
}

func GetKeyPressWithModKey(key glfw.Key, modKey glfw.ModifierKey) bool {
	if mKey == modKey {
		return keys[key].Press
	}
	return false
}

func GetKeyDown(key glfw.Key) bool {
	return keys[key].Press
}

func GetKeyUp(key glfw.Key) bool {
	return keys[key].Release
}

func GetDefaultAxis(axis Axis) float64 {
	switch axis {
	case Horizontal:
		return boolToFloat(GetKey(glfw.KeyA) || GetKey(glfw.KeyLeft))*-1 +
			boolToFloat(GetKey(glfw.KeyD) || GetKey(glfw.KeyRight))*1
	case Vertical:
		return boolToFloat(GetKey(glfw.KeyS) || GetKey(glfw.KeyDown))*-1 +
			boolToFloat(GetKey(glfw.KeyW) || GetKey(glfw.KeyUp))*1
	case MouseX:
		return posX
	case MouseY:
		return posY
	default:
		return 0
	}
}

func GetAxis(upKey, downKey glfw.Key) float64 {
	return boolToFloat(GetKey(downKey))*-1 +
		boolToFloat(GetKey(upKey))*1
}
