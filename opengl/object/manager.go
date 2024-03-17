package object

type Manager struct {
	Objects []*Object
}

func NewManager() *Manager {
	return &Manager{Objects: make([]*Object, 0)}
}

func (m *Manager) AddObject(object *Object) {
	m.Objects = append(m.Objects, object)
}
