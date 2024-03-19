package object

import "github.com/qmuntal/gltf"

type Node struct {
	nodeType    NodeType
	children    []*Node
	mesh        *Mesh
	translation [3]float32
	scale       [3]float32
	rotation    [4]float32
}

func parseNodes(doc *gltf.Document, meshes []*Mesh) []*Node {
	res := make([]*Node, 0)
	children := make([]uint32, len(doc.Nodes))
	recursiveNodesChildren(doc.Nodes, 0, children)
	for i, node := range doc.Nodes {
		n := &Node{
			nodeType:    Root,
			translation: node.Translation,
			scale:       node.Scale,
			rotation:    node.Rotation,
		}
		if node.Mesh != nil {
			n.mesh = meshes[*node.Mesh]
		}
		if children[i] > 0 {
			n.nodeType = Child
		}
		res = append(res, n)
	}
	return res
}

func recursiveNodesChildren(nodes []*gltf.Node, parentNode uint32, res []uint32) []uint32 {
	// key - child, val - parent
	if res == nil {
		res = make([]uint32, len(nodes))
	}
	if len(nodes[parentNode].Children) > 0 {
		for _, child := range nodes[parentNode].Children {
			res[child] = parentNode
			if len(nodes[child].Children) > 0 {
				recursiveNodesChildren(nodes, child, res)
			}
		}
	}
	return res
}
