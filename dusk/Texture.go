package dusk

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	glId uint32
}

func NewTexture(app *App, filename string) (*Texture, error) {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	LogLoad("Texture '%v'", filename)
	if filename == "" {
		return nil, fmt.Errorf("Filename cannot be empty")
	}

	data, err := app.AssetFunction(filename)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var glId uint32
	gl.GenTextures(1, &glId)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.ActiveTexture(gl.TEXTURE0)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.BindTexture(gl.TEXTURE_2D, glId)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	return &Texture{
		glId: glId,
	}, nil
}

func (tex *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, tex.glId)
}
