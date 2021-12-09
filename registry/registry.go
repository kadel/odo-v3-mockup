package registry

import (
	"github.com/devfile/registry-support/index/generator/schema"
	indexSchema "github.com/devfile/registry-support/index/generator/schema"
	registryLibrary "github.com/devfile/registry-support/registry-library/library"
)

type projectType struct {
	name                    string
	positionInRegistryIndex int
}

type DevfileIndex struct {
	registryUrl  string
	index        []indexSchema.Schema
	projectTypes map[string][]projectType
}

func (p *projectType) GetName() string {
	return p.name
}

func (p *projectType) GetPositionInIndex() int {
	return p.positionInRegistryIndex
}

func NewDevfileIndex(registryUrl string) *DevfileIndex {
	index := DevfileIndex{
		registryUrl: registryUrl,
	}
	index.projectTypes = make(map[string][]projectType)
	index.populateIndex()
	index.populateProjectTypes()

	index.GetProjectTypes("javascript")
	return &index
}

func (r *DevfileIndex) GetLanguages() []string {
	languages := []string{}
	for l := range r.projectTypes {
		languages = append(languages, l)
	}
	return languages
}

func (r *DevfileIndex) GetProjectTypes(lang string) []string {
	names := []string{}

	if projectType, ok := r.projectTypes[lang]; ok {
		for _, projectType := range projectType {
			names = append(names, projectType.GetName())
		}
	}

	return names
}

func (r *DevfileIndex) GetDevfileByIndex(number int) *indexSchema.Schema {
	for _, language := range r.projectTypes {
		for _, projectType := range language {
			if projectType.positionInRegistryIndex == number {
				return &r.index[number]
			}
		}
	}
	return nil
}

func (r *DevfileIndex) GetDevfileByName(name string) *indexSchema.Schema {
	for _, d := range r.index {
		if d.Name == name {
			return &d
		}
	}
	return nil
}

func (r *DevfileIndex) populateIndex() {
	if len(r.index) != 0 {
		return
	}
	registryIndex, err := registryLibrary.GetRegistryIndex(r.registryUrl, false, "", schema.StackDevfileType)

	if err != nil {
		panic(err)
	}

	r.index = registryIndex
}

func (r *DevfileIndex) populateProjectTypes() {
	for i, d := range r.index {
		if r.projectTypes[d.Language] == nil {
			r.projectTypes[d.Language] = make([]projectType, 0)
		}
		r.projectTypes[d.Language] = append(r.projectTypes[d.Language],
			projectType{
				name:                    d.DisplayName,
				positionInRegistryIndex: i,
			},
		)
	}
}

// return devfile information based on the language and index in its projectType slice
func (r *DevfileIndex) GetDevfile(language string, projectTypeIndex int) indexSchema.Schema {
	indexPosition := r.projectTypes[language][projectTypeIndex].positionInRegistryIndex
	return r.index[indexPosition]

}
