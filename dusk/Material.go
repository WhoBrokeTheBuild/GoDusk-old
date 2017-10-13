package dusk

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	AMBIENT_TEXID  = 0
	DIFFUSE_TEXID  = 1
	SPECULAR_TEXID = 2
	BUMP_TEXID     = 3
)

const (
	AMBIENT_MAP_FLAG  = 1
	DIFFUSE_MAP_FLAG  = 2
	SPECULAR_MAP_FLAG = 4
	BUMP_MAP_FLAG     = 8
)

type Material struct {
	ambient     mgl32.Vec3
	diffuse     mgl32.Vec3
	specular    mgl32.Vec3
	shininess   float32
	dissolve    float32
	ambientMap  *Texture
	diffuseMap  *Texture
	specularMap *Texture
	bumpMap     *Texture
	mapFlags    uint32
}

func NewMaterial(
	app *App,
	ambient, diffuse, specular mgl32.Vec3,
	shininess, dissolve float32,
	ambientMap, diffuseMap, specularMap, bumpMap string,
) (*Material, error) {
	var err error
	var ambientTex *Texture
	var diffuseTex *Texture
	var specularTex *Texture
	var bumpTex *Texture

	var flags uint32

	if ambientMap != "" {
		ambientTex, err = NewTexture(app, ambientMap)
		if err != nil {
			return nil, err
		}
		flags |= AMBIENT_MAP_FLAG
	}

	if diffuseMap != "" {
		diffuseTex, err = NewTexture(app, diffuseMap)
		if err != nil {
			return nil, err
		}
		flags |= DIFFUSE_MAP_FLAG
	}

	if specularMap != "" {
		specularTex, err = NewTexture(app, specularMap)
		if err != nil {
			return nil, err
		}
		flags |= SPECULAR_MAP_FLAG
	}

	if bumpMap != "" {
		bumpTex, err = NewTexture(app, bumpMap)
		if err != nil {
			return nil, err
		}
		flags |= BUMP_MAP_FLAG
	}

	return &Material{
		ambient:     ambient,
		diffuse:     diffuse,
		specular:    specular,
		shininess:   shininess,
		dissolve:    dissolve,
		ambientMap:  ambientTex,
		diffuseMap:  diffuseTex,
		specularMap: specularTex,
		bumpMap:     bumpTex,
		mapFlags:    flags,
	}, nil
}

func (mat *Material) Bind(shader *Shader) {
	gl.Uniform1ui(shader.GetUniformLocation("_MapFlags"), mat.mapFlags)
	gl.Uniform3fv(shader.GetUniformLocation("_Ambient"), 1, &mat.ambient[0])
	gl.Uniform3fv(shader.GetUniformLocation("_Diffuse"), 1, &mat.diffuse[0])
	gl.Uniform3fv(shader.GetUniformLocation("_Specular"), 1, &mat.specular[0])
	gl.Uniform1f(shader.GetUniformLocation("_Shininess"), mat.shininess)
	gl.Uniform1f(shader.GetUniformLocation("_Dissolve"), mat.dissolve)

	if mat.ambientMap != nil {
		gl.Uniform1i(shader.GetUniformLocation("_AmbientMap"), AMBIENT_TEXID)
		gl.ActiveTexture(gl.TEXTURE0 + AMBIENT_TEXID)
		mat.ambientMap.Bind()
	}

	if mat.diffuseMap != nil {
		gl.Uniform1i(shader.GetUniformLocation("_DiffuseMap"), DIFFUSE_TEXID)
		gl.ActiveTexture(gl.TEXTURE0 + DIFFUSE_TEXID)
		mat.diffuseMap.Bind()
	}

	if mat.specularMap != nil {
		gl.Uniform1i(shader.GetUniformLocation("_SpecularMap"), SPECULAR_TEXID)
		gl.ActiveTexture(gl.TEXTURE0 + SPECULAR_TEXID)
		mat.specularMap.Bind()
	}

	if mat.bumpMap != nil {
		gl.Uniform1i(shader.GetUniformLocation("_BumpMap"), BUMP_TEXID)
		gl.ActiveTexture(gl.TEXTURE0 + BUMP_TEXID)
		mat.bumpMap.Bind()
	}
}
