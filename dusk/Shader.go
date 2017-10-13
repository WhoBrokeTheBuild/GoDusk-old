package dusk

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	glId uint32
}

func NewShader(app *App, filenames ...string) (*Shader, error) {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	glProgId := gl.CreateProgram()
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	defer gl.DeleteShader(glProgId)

	glIds := []uint32{}
	defer func() {
		for _, glId := range glIds {
			gl.DeleteShader(glId)
			if glerr = gl.GetError(); glerr > 0 {
				log.Printf("gl.GetError returned %v", glerr)
			}
		}
	}()

	for _, f := range filenames {
		glId, err := compileShader(app, f)
		if err != nil {
			return nil, err
		}

		gl.AttachShader(glProgId, glId)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
	}

	gl.LinkProgram(glProgId)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	var status int32
	gl.GetProgramiv(glProgId, gl.LINK_STATUS, &status)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(glProgId, gl.INFO_LOG_LENGTH, &logLength)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}

		programLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(glProgId, logLength, nil, gl.Str(programLog))
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}

		return nil, fmt.Errorf("Failed to link shader program: %v", programLog)
	}

	shader := &Shader{
		glId: glProgId,
	}
	glProgId = 0

	return shader, nil
}

func (shader *Shader) Cleanup() {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gl.DeleteShader(shader.glId)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
}

func (shader *Shader) Use() {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gl.UseProgram(shader.glId)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
}

func (shader *Shader) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(shader.glId, gl.Str(name+"\x00"))
}

func compileShader(app *App, filename string) (uint32, error) {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var shaderType uint32
	if strings.HasSuffix(filename, ".vs.glsl") {
		shaderType = gl.VERTEX_SHADER
	} else if strings.HasSuffix(filename, ".fs.glsl") {
		shaderType = gl.FRAGMENT_SHADER
	} else if strings.HasSuffix(filename, ".gs.glsl") {
		shaderType = gl.GEOMETRY_SHADER
	}

	data, err := app.AssetFunction(filename)
	if err != nil {
		return 0, err
	}

	source := string(data) + "\x00"

	glId := gl.CreateShader(shaderType)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	glSources, free := gl.Strs(source)
	gl.ShaderSource(glId, 1, glSources, nil)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	free()

	gl.CompileShader(glId)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	var status int32
	gl.GetShaderiv(glId, gl.COMPILE_STATUS, &status)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(glId, gl.INFO_LOG_LENGTH, &logLength)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}

		shaderLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(glId, logLength, nil, gl.Str(shaderLog))
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}

		return glId, fmt.Errorf("Failed to compile shader '%v': %v", filename, shaderLog)
	}

	return glId, nil
}
