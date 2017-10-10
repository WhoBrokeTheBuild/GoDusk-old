package dusk

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	Transform mgl32.Mat4

	glVao  uint32
	glVbos [4]uint32
}

func NewModel() (Model, error) {
	return Model{
		Transform: mgl32.Mat4{},
		glVao:     0,
		glVbos:    [4]uint32{0, 0, 0, 0},
	}, nil
}

func NewModelFromFile(filename string) (Model, error) {
	model, err := NewModel()
	if err != nil {
		return model, err
	}
	err = model.LoadFromFile(filename)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (model Model) LoadFromFile(filename string) error {
	dirname := path.Dir(filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024)

	var line, cmd string

	tmp, _, err := reader.ReadLine()
	for err == nil {
		line = string(tmp)

		fmt.Sscanf(line, "%s", &cmd)
        switch cmd {
        case "usemtl":

        case "mtllib":

        case "o": fallthrough
        case "g":

        }

        if cmd == ""
        if cmd == "mtllib" {

        }
        else if cmd == "usemtl" {

        }
		fmt.Println("Command: ", cmd)
		fmt.Println("Data: ", line)

		tmp, _, err = reader.ReadLine()
	}

	_ = dirname
	return nil
}
