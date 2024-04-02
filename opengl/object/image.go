package object

import "github.com/qmuntal/gltf"

type Image struct {
	Name string
	URI  string
}

func (o *Object) parseImages(doc *gltf.Document) {
	o.Images = make([]*Image, len(doc.Images))
	for i, image := range doc.Images {
		img := &Image{
			Name: image.Name,
			URI:  image.URI,
		}
		o.Images[i] = img
	}
}
