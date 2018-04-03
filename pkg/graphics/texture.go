package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"os"
	"image"
	"image/draw"
	"log"
	_ "image/jpeg"
	_ "image/png"
)

type Texture struct {
	id     uint32
	width  int32
	height int32
}

func NewTextureFromFile(filePath string) (*Texture) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panicf("cannot find file: '%s'", filePath)
		return nil
	}
	defer file.Close()

	decodedImage, format, err := image.Decode(file)
	if err != nil {
		log.Panicf("cannot decode image <%s>: '%s'", format, filePath)
		return nil
	}
	return NewTextureFromImage(decodedImage)
}

func NewTextureFromImage(imageData image.Image) (*Texture) {
	switch imageData.(type) {
	case *image.RGBA:
	default:
		rgba := image.NewRGBA(imageData.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			log.Panicf("unsupported stride")
			return nil
		}
		draw.Draw(rgba, rgba.Bounds(), imageData, image.Point{0, 0}, draw.Src)
		imageData = rgba
	}

	texture := Texture{
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

	return &texture
}

func NewEmptyTexture(width int, height int) (*Texture, error) {
	bounds := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	}
	imageData := image.NewRGBA(bounds)

	texture := Texture{
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

	return &texture, nil
}

func (t *Texture) Id() (uint32) {
	return t.id
}
