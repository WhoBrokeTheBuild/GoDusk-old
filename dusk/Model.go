package dusk

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	VERT_ATTRIB = 0
	NORM_ATTRIB = 1
	TXCD_ATTRIB = 2
)

type modelGroup struct {
	DrawMode uint32
	Start    int32
	Count    int32
	Material *Material
}

type Model struct {
	Transform mgl32.Mat4

	glVao  uint32
	glVbos [3]uint32
	groups []modelGroup
}

func NewModel(app *App) (*Model, error) {
	return &Model{
		Transform: mgl32.Ident4(),
		glVao:     0,
		glVbos:    [3]uint32{0, 0, 0},
	}, nil
}

func NewModelFromFile(app *App, filename string) (*Model, error) {
	model, err := NewModel(app)
	if err != nil {
		return model, err
	}
	err = model.LoadFromFile(app, filename)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (model *Model) Cleanup() {
	gl.DeleteBuffers(3, &model.glVbos[0])
	gl.DeleteVertexArrays(1, &model.glVao)
}

func (model *Model) LoadFromFile(app *App, filename string) error {
	LogLoad("Model '%v'", filename)

	// Holds a material
	type MatDef struct {
		Ambient     mgl32.Vec3
		Diffuse     mgl32.Vec3
		Specular    mgl32.Vec3
		Shininess   float32
		Dissolve    float32
		AmbientMap  string
		SpecularMap string
		DiffuseMap  string
		BumpMap     string
	}

	// Holds a single face
	type Face struct {
		VertInds [3]int
		NormInds [3]int
		TxcdInds [3]int
	}

	// Holds a group of faces and a material
	type Group struct {
		Name     string
		Material string
		Faces    []Face
	}

	LoadMaterials := func(filename string) (map[string]*MatDef, error) {
		LogLoad("Material '%v'", filename)

		materials := map[string]*MatDef{}

		dirname := path.Dir(filename)

		// Open the .mtl file
		data, err := app.AssetFunction(filename)
		if err != nil {
			return materials, err
		}

		// Create a reader with a specific buffer size, needed by reader.ReadLine()
		reader := bufio.NewReader(bytes.NewReader(data))

		var line string
		var curMat string

		tmp, _, err := reader.ReadLine()
		for ; err == nil; tmp, _, err = reader.ReadLine() {
			line = string(tmp)

			// Ignore empty lines and comments
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			// Split on the first ' ', ignore half lines
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 0 {
				continue
			}

			switch parts[0] {
			case "newmtl":

				curMat = parts[1]
				materials[curMat] = &MatDef{
					Ambient:     mgl32.Vec3{0, 0, 0},
					Diffuse:     mgl32.Vec3{0, 0, 0},
					Specular:    mgl32.Vec3{0, 0, 0},
					Shininess:   0.0,
					Dissolve:    0.0,
					AmbientMap:  "",
					SpecularMap: "",
					DiffuseMap:  "",
					BumpMap:     "",
				}

			case "Ka":
				fmt.Sscanf(parts[1], "%f %f %f",
					&materials[curMat].Ambient[0],
					&materials[curMat].Ambient[1],
					&materials[curMat].Ambient[2],
				)
			case "Kd":
				fmt.Sscanf(parts[1], "%f %f %f",
					&materials[curMat].Diffuse[0],
					&materials[curMat].Diffuse[1],
					&materials[curMat].Diffuse[2],
				)
			case "Ks":
				fmt.Sscanf(parts[1], "%f %f %f",
					&materials[curMat].Specular[0],
					&materials[curMat].Specular[1],
					&materials[curMat].Specular[2],
				)
			case "Ns":
				fmt.Sscanf(parts[1], "%f", &materials[curMat].Shininess)
			case "d":
				fmt.Sscanf(parts[1], "%f", &materials[curMat].Dissolve)
			case "map_Ka":
				materials[curMat].AmbientMap = path.Join(dirname, parts[1])
			case "map_Kd":
				materials[curMat].DiffuseMap = path.Join(dirname, parts[1])
			case "map_Ks":
				materials[curMat].SpecularMap = path.Join(dirname, parts[1])
			case "map_bump":
				materials[curMat].BumpMap = path.Join(dirname, parts[1])
			}
		}

		return materials, nil
	}

	// Open the .obj file
	data, err := app.AssetFunction(filename)
	if err != nil {
		return err
	}

	// Get the directory name for loading .mtl files
	dirname := path.Dir(filename)

	// Create a reader with a specific buffer size, needed by reader.ReadLine()
	reader := bufio.NewReader(bytes.NewReader(data))

	materials := map[string]*MatDef{}

	// Create a list of groups, and get a pointer to the first
	groups := []Group{{}}
	group := &groups[0]

	// Create the list of all Vertices, Normals, and Texture Coordinates
	allVerts := []mgl32.Vec3{}
	allNorms := []mgl32.Vec3{}
	allTxcds := []mgl32.Vec2{}

	var line string
	var count int

	tmpVec3 := mgl32.Vec3{}
	tmpVec2 := mgl32.Vec2{}
	tmpFace := Face{}

	tmp, _, err := reader.ReadLine()
	for ; err == nil; tmp, _, err = reader.ReadLine() {
		line = string(tmp)

		// Ignore empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Split on the first ' ', ignore half lines
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "usemtl":

			group.Material = parts[1]

		case "mtllib":

			newmats, err := LoadMaterials(path.Join(dirname, parts[1]))
			if err != nil {
				return err
			}
			for k, v := range newmats {
				materials[k] = v
			}

		case "o":
			fallthrough
		case "g":

			if group.Name == "" {
				group.Name = parts[1]
			} else {
				groups = append(groups, Group{
					Name: parts[1],
				})
				group = &groups[len(groups)-1]
			}

		case "f":

			// Test for and parse faces in the 'v//vn v//vn v//vn' format
			if strings.Contains(parts[1], "//") {
				count, err = fmt.Sscanf(parts[1],
					"%d//%d %d//%d %d//%d",
					&tmpFace.VertInds[0],
					&tmpFace.NormInds[0],
					&tmpFace.VertInds[1],
					&tmpFace.NormInds[1],
					&tmpFace.VertInds[2],
					&tmpFace.NormInds[2],
				)
				if err != nil || count != 6 {
					return fmt.Errorf("Malformed OBJ file '%v' %v", line)
				}
				// Test for and parse faces in the 'v/vt v/vt v/vt' format
			} else if strings.Count(parts[1], "/") == 3 {
				count, err = fmt.Sscanf(parts[1],
					"%d/%d %d/%d %d/%d",
					&tmpFace.VertInds[0],
					&tmpFace.TxcdInds[0],
					&tmpFace.VertInds[1],
					&tmpFace.TxcdInds[1],
					&tmpFace.VertInds[2],
					&tmpFace.TxcdInds[2],
				)
				if err != nil || count != 6 {
					return fmt.Errorf("Malformed OBJ file '%v'", line)
				}
				// Test for and parse faces in the 'v/vt/vn v/vt/vn v/vt/vn' format
			} else if strings.Count(parts[1], "/") == 6 {
				count, err = fmt.Sscanf(parts[1],
					"%d/%d/%d %d/%d/%d %d/%d/%d",
					&tmpFace.VertInds[0],
					&tmpFace.TxcdInds[0],
					&tmpFace.NormInds[0],
					&tmpFace.VertInds[1],
					&tmpFace.TxcdInds[1],
					&tmpFace.NormInds[1],
					&tmpFace.VertInds[2],
					&tmpFace.TxcdInds[2],
					&tmpFace.NormInds[2],
				)
				if err != nil || count != 9 {
					return fmt.Errorf("Malformed OBJ file '%v'", line)
				}
			} else {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			group.Faces = append(group.Faces, tmpFace)

		case "v":

			count, err = fmt.Sscanf(parts[1], "%f %f %f", &tmpVec3[0], &tmpVec3[1], &tmpVec3[2])
			if err != nil || count != 3 {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			allVerts = append(allVerts, tmpVec3)

		case "vn":

			count, err = fmt.Sscanf(parts[1], "%f %f %f", &tmpVec3[0], &tmpVec3[1], &tmpVec3[2])
			if err != nil || count != 3 {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			allNorms = append(allNorms, tmpVec3)

		case "vt":

			count, err = fmt.Sscanf(parts[1], "%f %f", &tmpVec2[0], &tmpVec2[1])
			if err != nil || count != 2 {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			allTxcds = append(allTxcds, tmpVec2)

		}
	}

	if group.Name == "" {
		group.Name = "default"
	}

	start := int32(0)
	verts := []float32{}
	norms := []float32{}
	txcds := []float32{}

	for g := range groups {
		group := &groups[g]
		for f := range group.Faces {
			face := &group.Faces[f]
			for i := 0; i < 3; i++ {
				// Adjust for negative indices
				if face.VertInds[i] < 0 {
					face.VertInds[i] += len(allVerts)
				}
				if face.NormInds[i] < 0 {
					face.NormInds[i] += len(allNorms)
				}
				if face.TxcdInds[i] < 0 {
					face.TxcdInds[i] += len(allTxcds)
				}

				// Adjust for zero-indexing
				face.VertInds[i] -= 1
				face.NormInds[i] -= 1
				face.TxcdInds[i] -= 1

				// Copy data to final arrays
				verts = append(verts,
					allVerts[face.VertInds[i]][0],
					allVerts[face.VertInds[i]][1],
					allVerts[face.VertInds[i]][2],
				)
				if face.NormInds[i] >= 0 {
					norms = append(norms,
						allNorms[face.NormInds[i]][0],
						allNorms[face.NormInds[i]][1],
						allNorms[face.NormInds[i]][2],
					)
				}
				if face.TxcdInds[i] >= 0 {
					txcds = append(txcds,
						allTxcds[face.TxcdInds[i]][0],
						allTxcds[face.TxcdInds[i]][1],
					)
				}
			}
		}

		var newMaterial *Material
		if mat, ok := materials[group.Material]; ok {
			newMaterial, err = NewMaterial(
				app,
				mat.Ambient,
				mat.Diffuse,
				mat.Specular,
				mat.Shininess,
				mat.Dissolve,
				mat.AmbientMap,
				mat.DiffuseMap,
				mat.SpecularMap,
				mat.BumpMap,
			)
			if err != nil {
				return err
			}
		}
		vertCount := int32(len(group.Faces) * 3 * 3)
		model.groups = append(model.groups, modelGroup{
			DrawMode: gl.TRIANGLES,
			Start:    start,
			Count:    vertCount,
			Material: newMaterial,
		})
		start += vertCount
	}

	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gl.GenVertexArrays(1, &model.glVao)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	gl.BindVertexArray(model.glVao)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	gl.GenBuffers(3, &model.glVbos[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, model.glVbos[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.VertexAttribPointer(VERT_ATTRIB, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.EnableVertexAttribArray(VERT_ATTRIB)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	if len(norms) == 0 {
		gl.DeleteBuffers(1, &model.glVbos[1])
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		model.glVbos[1] = 0
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, model.glVbos[1])
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(norms)*4, gl.Ptr(norms), gl.STATIC_DRAW)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		gl.VertexAttribPointer(NORM_ATTRIB, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		gl.EnableVertexAttribArray(NORM_ATTRIB)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
	}

	if len(txcds) == 0 {
		gl.DeleteBuffers(1, &model.glVbos[2])
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		model.glVbos[2] = 0
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, model.glVbos[2])
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(txcds)*4, gl.Ptr(txcds), gl.STATIC_DRAW)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		gl.VertexAttribPointer(TXCD_ATTRIB, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
		gl.EnableVertexAttribArray(TXCD_ATTRIB)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
	}

	return nil
}

func (model *Model) Render(shader *Shader) {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gl.BindVertexArray(model.glVao)
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	for g := range model.groups {
		group := &model.groups[g]
		if group.Material != nil {
			group.Material.Bind(shader)
		}
		gl.DrawArrays(group.DrawMode, group.Start, group.Count)
		if glerr = gl.GetError(); glerr > 0 {
			log.Printf("gl.GetError returned %v", glerr)
		}
	}
}
