//+build !vulkan

package common

import "github.com/EngoEngine/gl"

type BufferData struct {
	Buffer        *gl.Buffer
	BufferContent []float32
}
