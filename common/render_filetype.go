package common

import (
	"fmt"
	"image"
	"image/draw"
	// imported to decode jpegs and upload them to the GPU.
	_ "image/jpeg"
	// imported to decode .pngs and upload them to the GPU.
	_ "image/png"
	// imported to decode .gifs and uppload them to the GPU.
	_ "image/gif"
	"io"

	// these are for svg support

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	"github.com/EngoEngine/engo"
)

type imageLoader struct {
	images map[string]TextureResource
}

func (i *imageLoader) Load(url string, data io.Reader) error {
	var res TextureResource
	if getExt(url) == ".svg" {
		icon, err := oksvg.ReadIconStream(data, oksvg.WarnErrorMode)
		if err != nil {
			return err
		}
		w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		gv := rasterx.NewScannerGV(w, h, img, img.Bounds())
		r := rasterx.NewDasher(w, h, gv)
		icon.Draw(r, 1.0)
		b := img.Bounds()
		newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)
		res = NewTextureResource(&ImageObject{newm})
	} else {
		img, _, err := image.Decode(data)
		if err != nil {
			return err
		}
		b := img.Bounds()
		newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)
		res = NewTextureResource(&ImageObject{newm})
	}
	res.url = url
	i.images[url] = res

	return nil
}

func (i *imageLoader) Unload(url string) error {
	delete(i.images, url)
	return nil
}

func (i *imageLoader) Resource(url string) (engo.Resource, error) {
	texture, ok := i.images[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return texture, nil
}

// Image holds data and properties of an .jpg, .gif, or .png file
type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

// ImageToNRGBA takes a given `image.Image` and converts it into an `image.NRGBA`. Especially useful when transforming
// image.Uniform to something usable by `engo`.
func ImageToNRGBA(img image.Image, width, height int) *image.NRGBA {
	newm := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newm, newm.Bounds(), img, image.Point{0, 0}, draw.Src)

	return newm
}

// ImageObject is a pure Go implementation of a `Drawable`
type ImageObject struct {
	data *image.NRGBA
}

// NewImageObject creates a new ImageObject given the image.NRGBA reference
func NewImageObject(img *image.NRGBA) *ImageObject {
	return &ImageObject{img}
}

// Data returns the entire image.NRGBA object
func (i *ImageObject) Data() interface{} {
	return i.data
}

// Width returns the maximum X coordinate of the image
func (i *ImageObject) Width() int {
	return i.data.Rect.Max.X
}

// Height returns the maximum Y coordinate of the image
func (i *ImageObject) Height() int {
	return i.data.Rect.Max.Y
}

// LoadedSprite loads the texture-reference from `engo.Files`, and wraps it in a `*Texture`.
// This method is intended for image-files which represent entire sprites.
func LoadedSprite(url string) (*Texture, error) {
	res, err := engo.Files.Resource(url)
	if err != nil {
		return nil, err
	}

	img, ok := res.(TextureResource)
	if !ok {
		return nil, fmt.Errorf("resource not of type `TextureResource`: %s", url)
	}

	return &Texture{img.Texture.ID, img.Width, img.Height, engo.AABB{Max: engo.Point{X: 1.0, Y: 1.0}}}, nil
}

func init() {
	engo.Files.Register(".jpg", &imageLoader{images: make(map[string]TextureResource)})
	engo.Files.Register(".png", &imageLoader{images: make(map[string]TextureResource)})
	engo.Files.Register(".gif", &imageLoader{images: make(map[string]TextureResource)})
	engo.Files.Register(".svg", &imageLoader{images: make(map[string]TextureResource)})
}
