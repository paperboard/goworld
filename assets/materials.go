package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/johanhenriksson/goworld/render"
)

// MaterialDefinition file json representation
type MaterialDefinition struct {
	Shader   string
	Buffers  map[string][]*VertexPointerDefinition
	Textures map[string]*TextureDefinition
}

// VertexPointerDefinition json representation
type VertexPointerDefinition struct {
	Name      string
	Type      string
	GlType    render.GLType `json:"-"`
	Size      int
	Offset    int
	Count     int
	Normalize bool
	Integer   bool
}

// TextureDefinition json representation
type TextureDefinition struct {
	File   string
	Filter string
	Wrap   string
}

func LoadMaterialDefinition(file string) (*MaterialDefinition, error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	matf := &MaterialDefinition{}
	err = json.Unmarshal(jsonBytes, matf)
	if err != nil {
		return nil, err
	}

	return matf, nil
}

func LoadMaterial(name string, matf *MaterialDefinition) (*render.Material, error) {
	shader := GetShader(matf.Shader)

	mat := render.CreateMaterial(name, shader)

	// load textures
	for name, txtf := range matf.Textures {
		texture, err := render.TextureFromFile(txtf.File)
		if err != nil {
			return nil, err
		}
		if txtf.Filter == "nearest" {
			texture.SetFilter(render.NearestFilter)
		}
		mat.Textures.Add(name, texture)
	}

	return mat, nil
}

// GetMaterial returns a new instance of a material
func GetMaterial(name string) *render.Material {
	path := fmt.Sprintf("assets/materials/%s.json", name)
	def, err := LoadMaterialDefinition(path)
	if err != nil {
		panic(fmt.Errorf("failed to load material definition %s: %s", name, err))
	}

	if def.Shader == "" {
		def.Shader = name
	}

	mat, err := LoadMaterial(name, def)
	if err != nil {
		panic(fmt.Errorf("failed to load material %s: %s", name, err))
	}

	return mat
}

// GetMaterialShared returns a shared instance of a material
func GetMaterialShared(name string) *render.Material {
	if mat, exists := cache.Materials[name]; exists {
		return mat
	}

	mat := GetMaterial(name)
	cache.Materials[name] = mat

	return mat
}
