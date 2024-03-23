package object

import "github.com/qmuntal/gltf"

type Image struct {
	Name string
	URI  string
}

func parseImages(doc *gltf.Document) []*Image {
	res := make([]*Image, len(doc.Images))
	for i, image := range doc.Images {
		img := &Image{
			Name: image.Name,
			URI:  image.URI,
		}
		res[i] = img
	}
	return res
}
