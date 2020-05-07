//+build !vulkan

package common

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/gl"
)

// TextureResource is the resource used by the RenderSystem. It uses .jpg, .gif, and .png images
type TextureResource struct {
	Texture *gl.Texture
	Width   float32
	Height  float32
	url     string
}

// URL is the file path of the TextureResource
func (t TextureResource) URL() string {
	return t.url
}

// UploadTexture sends the image to the GPU, to be kept in GPU RAM
func UploadTexture(img Image) *gl.Texture {
	var id *gl.Texture
	if !engo.Headless() {
		id = engo.Gl.CreateTexture()

		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, id)

		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_S, engo.Gl.CLAMP_TO_EDGE)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_T, engo.Gl.CLAMP_TO_EDGE)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MIN_FILTER, engo.Gl.LINEAR)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MAG_FILTER, engo.Gl.NEAREST)

		if img.Data() == nil {
			panic("Texture image data is nil.")
		}

		engo.Gl.TexImage2D(engo.Gl.TEXTURE_2D, 0, engo.Gl.RGBA, engo.Gl.RGBA, engo.Gl.UNSIGNED_BYTE, img.Data())
	}
	return id
}

// NewTextureResource sends the image to the GPU and returns a `TextureResource` for easy access
func NewTextureResource(img Image) TextureResource {
	id := UploadTexture(img)
	return TextureResource{Texture: id, Width: float32(img.Width()), Height: float32(img.Height())}
}

// NewTextureSingle sends the image to the GPU and returns a `Texture` with a viewport for single-sprite images
func NewTextureSingle(img Image) Texture {
	id := UploadTexture(img)
	return Texture{id, float32(img.Width()), float32(img.Height()), engo.AABB{Max: engo.Point{X: 1.0, Y: 1.0}}}
}

// Texture represents a texture loaded in the GPU RAM (by using OpenGL), which defined dimensions and viewport
type Texture struct {
	ID       *gl.Texture
	width    float32
	height   float32
	viewport engo.AABB
}

// Width returns the width of the texture.
func (t Texture) Width() float32 {
	return t.width
}

// Height returns the height of the texture.
func (t Texture) Height() float32 {
	return t.height
}

// Texture returns the OpenGL ID of the Texture.
func (t Texture) Texture() *Texture {
	return &t
}

// View returns the viewport properties of the Texture. The order is Min.X, Min.Y, Max.X, Max.Y.
func (t Texture) View() (float32, float32, float32, float32) {
	return t.viewport.Min.X, t.viewport.Min.Y, t.viewport.Max.X, t.viewport.Max.Y
}

// Close removes the Texture data from the GPU.
func (t Texture) Close() {
	if !engo.Headless() {
		engo.Gl.DeleteTexture(t.ID)
	}
}
