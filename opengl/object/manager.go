package object

import "github.com/go-gl/mathgl/mgl32"

type Manager struct {
	Objects []*Object
}

func NewManager() *Manager {
	return &Manager{Objects: make([]*Object, 0)}
}

func (m *Manager) AddObject(object *Object) {
	m.Objects = append(m.Objects, object)
}

func (m *Manager) AddPlain() {
	plain := &Object{
		Pos:       mgl32.Vec3{0, 0, 0},
		MainScene: nil,
		Scene:     nil,
		Nodes:     nil,
		Meshes:    nil,
	}
	m.Objects = append(m.Objects, plain)
}
