package graphics

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	// Used only to initialize the JPEG subsystem
	_ "image/jpeg"
	// Used only to initialize the PNG subsystem
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Texture a representation of an image file in memory
type Texture struct {
	id     uint32
	width  int32
	height int32
}

// NewTextureFromFile loads the image from a file into a texture
func NewTextureFromFile(filePath string) *Texture {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error loading texture. %s", err)
		return nil
	}
	defer file.Close()

	decodedImage, format, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding <%s> image: '%s'", format, filePath)
		return nil
	}
	return NewTextureFromImage(decodedImage)
}

// NewTextureFromImage uses the data from an Image struct to create a texture
func NewTextureFromImage(imageData image.Image) *Texture {
	switch imageData.(type) {
	case *image.RGBA:
	default:
		rgba := image.NewRGBA(imageData.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			fmt.Printf("Error creating texture: unsupported stride")
			return nil
		}
		draw.Draw(rgba, rgba.Bounds(), imageData, image.Point{0, 0}, draw.Src)
		imageData = rgba
	}

	texture := &Texture{
		width:  int32(imageData.Bounds().Dx()),
		height: int32(imageData.Bounds().Dy()),
	}
	pixelData := imageData.(*image.RGBA).Pix

	gl.GenTextures(1, &texture.id)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, texture.width, texture.height,
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixelData),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture
}

// NewEmptyTexture creates an empty texture with a specified size
func NewEmptyTexture(width int, height int) (*Texture, error) {
	bounds := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	}
	imageData := image.NewRGBA(bounds)

	texture := &Texture{
		width:  int32(imageData.Bounds().Dx()),
		height: int32(imageData.Bounds().Dy()),
	}
	gl.GenTextures(1, &texture.id)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, texture.width, texture.height,
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imageData.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture, nil
}

// ID returns the unique OpenGL ID of this texture
func (t *Texture) ID() uint32 {
	return t.id
}

// Width returns the texture width in pixels
func (t *Texture) Width() int32 {
	return t.width
}

// Height returns the texture width in pixels
func (t *Texture) Height() int32 {
	return t.height
}
