package objectManager

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"opengl/opengl/object"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type Manager struct {
	modelsDir string
	models    map[string]modelInfo
	Objects   []*object.Object
}

type modelInfo struct {
	folder bool
}

func NewManager(coreDir string) *Manager {
	return &Manager{
		Objects:   make([]*object.Object, 0),
		models:    make(map[string]modelInfo),
		modelsDir: coreDir,
	}
}

func (m *Manager) Init() error {
	// search for all models in core folder
	coreFolder, err := os.Open(m.modelsDir)
	if err != nil {
		return err
	}
	files, err := coreFolder.ReadDir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := strings.Split(file.Name(), ".")
		m.models[fileName[0]] = modelInfo{file.IsDir()}
	}
	return nil
}

func (m *Manager) openDoc(name string) (*gltf.Document, error) {
	docFile := ""

	model, ok := m.models[name]
	if !ok {
		return nil, fmt.Errorf("object '%s' not found", name)
	}

	if model.folder {
		folder, err := os.Open(path.Join(m.modelsDir, name))
		if err != nil {
			return nil, err
		}
		files, err := folder.ReadDir(0)
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".gltf") {
				docFile = path.Join(name, f.Name())
				break
			}
		}
	} else {
		docFile = name + ".glb"
	}

	doc, err := gltf.Open(path.Join(m.modelsDir, docFile))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (m *Manager) NewObject(pos mgl32.Vec3, scale float32, name string) error {
	t := time.Now()

	doc, err := m.openDoc(name)
	if err != nil {
		return err
	}

	fmt.Printf("parsing model:%s\n", name)
	fmt.Println("\ttotal mesh:", len(doc.Meshes))
	fmt.Println("\ttotal textures:", len(doc.Textures))
	fmt.Println("\ttotal images:", len(doc.Images))

	obj := &object.Object{Pos: pos, Scale: scale}
	withErr := obj.Parse(doc, path.Join(m.modelsDir, name))

	fmt.Printf("Done in %d milliseconds. ", time.Since(t).Milliseconds())
	fmt.Printf("With errors: %t\n", withErr)
	runtime.GC()

	m.AddObject(obj)

	return nil
}

func (m *Manager) AddObject(object *object.Object) {
	m.Objects = append(m.Objects, object)
}

func (m *Manager) AddPlain() {
	plain := &object.Object{
		Pos:       mgl32.Vec3{0, 0, 0},
		Scale:     1,
		MainScene: nil,
		Scene:     nil,
		Nodes:     nil,
		Meshes:    nil,
	}
	m.Objects = append(m.Objects, plain)
}
