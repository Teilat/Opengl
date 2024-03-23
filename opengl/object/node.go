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

func (n *Node) Draw(program uint32, absPos mgl32.Vec3) {
	if n.mesh != nil {
		gl.BindTexture(gl.TEXTURE_2D, n.mesh.Texture1Id)
		gl.BindVertexArray(n.mesh.Vao)

		mMatrix := n.matrix

		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str(modelData)), 1, false, n.getDataP())

		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str(modelMatrix)), 1, false, &mMatrix[0])
		gl.DrawElements(gl.TRIANGLES, int32(len(n.mesh.Indices)), gl.UNSIGNED_INT, nil)
	}
}

func (n *Node) getDataP() *float32 {
	data := mgl32.Mat4x3FromCols(mgl32.Vec4{n.translation[0], n.translation[1], n.translation[2], 0}, n.rotation, mgl32.Vec4{n.scale[0], n.scale[1], n.scale[2], 0})
	return &data[0]
}

func parseNodes(doc *gltf.Document, meshes []*Mesh) []*Node {
	res := make([]*Node, len(doc.Nodes))
	parents := make([]uint32, len(doc.Nodes))
	recursiveNodeParent(doc.Nodes, 0, parents)
	for i, node := range doc.Nodes {
		n := &Node{
			Name:        node.Name,
			translation: node.TranslationOrDefault(),
			scale:       node.ScaleOrDefault(),
			rotation:    node.RotationOrDefault(),
			matrix:      node.MatrixOrDefault(),
			children:    make([]*Node, 0),
		}

		if node.Mesh != nil {
			n.mesh = meshes[*node.Mesh]
		}
		res[i] = n
	}
	for i, node := range res {
		// parent
		if parents[i] != uint32(i) {
			node.parent = res[parents[i]]
		}
		// children
		for _, child := range doc.Nodes[i].Children {
			node.children = append(node.children, res[child])
		}
	}
	return res
}

func recursiveNodeParent(nodes []*gltf.Node, parentNode uint32, res []uint32) []uint32 {
	// key - child, val - parent
	if res == nil {
		res = make([]uint32, len(nodes))
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
