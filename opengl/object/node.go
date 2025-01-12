package object

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
)

const (
	modelData   = "modelData\x00"
	modelMatrix = "modelMatrix\x00"
)

type Node struct {
	Name        string
	children    []*Node
	parent      *Node
	mesh        *Mesh
	matrix      mgl32.Mat4
	translation mgl32.Vec3
	scale       mgl32.Vec3
	rotation    mgl32.Vec4
}

func (n *Node) Draw(program uint32) {
	if n.mesh != nil {
		for _, p := range n.mesh.Primitives {
			gl.BindTexture(gl.TEXTURE_2D, p.TextureId)
			gl.BindVertexArray(p.Vao)

			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str(modelData)), 1, false, n.getDataP())
			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str(modelMatrix)), 1, false, &n.matrix[0])
			gl.DrawElements(p.PrimitiveMode, int32(len(p.Indices)), gl.UNSIGNED_INT, nil)
		}
	}
}

func (n *Node) getDataP() *float32 {
	data := mgl32.Mat4x3FromCols(
		mgl32.Vec4{n.translation[0], n.translation[1], n.translation[2], 0},
		n.rotation,
		mgl32.Vec4{n.scale[0], n.scale[1], n.scale[2], 0},
	)
	return &data[0]
}

func (o *Object) parseNodes(doc *gltf.Document) {
	o.Nodes = make([]*Node, len(doc.Nodes))
	parents := recursiveNodeParent(doc.Nodes, 0, nil)
	for i, node := range doc.Nodes {
		n := &Node{
			Name:        node.Name,
			translation: toVec3f32(node.TranslationOrDefault()),
			scale:       toVec3f32(node.ScaleOrDefault()),
			rotation:    toVec4f32(node.RotationOrDefault()),
			matrix:      toMat4f32(node.MatrixOrDefault()),
			children:    make([]*Node, 0),
		}
		n.translation = n.translation.Add(o.Pos)
		n.scale = n.scale.Mul(o.Scale)

		if node.Mesh != nil {
			n.mesh = o.Meshes[*node.Mesh]
		}
		o.Nodes[i] = n
	}
	for i, node := range o.Nodes {
		// parent
		if parents[i] != i {
			node.parent = o.Nodes[parents[i]]
		}
		// children
		for _, child := range doc.Nodes[i].Children {
			node.children = append(node.children, o.Nodes[child])
		}
	}
}

func recursiveNodeParent(nodes []*gltf.Node, parentNode int, res []int) []int {
	// key - child, val - parent
	if res == nil {
		res = make([]int, len(nodes))
	}
	if len(nodes[parentNode].Children) > 0 {
		for _, child := range nodes[parentNode].Children {
			res[child] = parentNode
			if len(nodes[child].Children) > 0 {
				recursiveNodeParent(nodes, child, res)
			}
		}
	}
	return res
}
