package app

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
)

func init() {
	starters = append(starters, startDestination)
}

func startDestination(ctx context.Context, a *Application) error {
	typesList, err := a.DB.LibraryTypeList(1, -1)
	if err != nil {
		return fae.Wrap(err, "getting library types")
	}
	templatesList, err := a.DB.LibraryTemplateList(1, -1)
	if err != nil {
		return fae.Wrap(err, "getting library templates")
	}
	librariesList, err := a.DB.LibraryList(1, -1)
	if err != nil {
		return fae.Wrap(err, "getting libraries")
	}

	types := map[string]*LibraryType{}
	for _, t := range typesList {
		types[t.Name] = t
	}
	templates := map[string]*LibraryTemplate{}
	for _, t := range templatesList {
		templates[t.Name] = t
	}
	libraries := map[string]*Library{}
	for _, l := range librariesList {
		libraries[l.Name] = l
	}

	a.Destinator = &Destinator{
		libraries: libraries,
		types:     types,
		templates: templates,
		// processedTemplates: map[string]string{},
	}

	// if err := a.Destinator.processTemplates(); err != nil {
	// 	return fae.Wrap(err, "processing templates")
	// }

	return nil
}

type Destinator struct {
	libraries map[string]*Library
	types     map[string]*LibraryType
	templates map[string]*LibraryTemplate

	// processedTemplates map[string]string
}

// var processTemplatesRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// func (d *Destinator) processTemplates() error {
// 	for _, t := range d.templates {
// 		match := processTemplatesRegex.FindAllStringSubmatch(t.Template, -1)
// 		for _, m := range match {
// 			if len(m) != 2 {
// 				continue
// 			}
// 			if v, ok := d.templates[m[1]]; ok {
// 				t.Template = strings.Replace(t.Template, "{{"+m[1]+"}}", v.Template, -1)
// 			}
// 		}
// 	}
//
// 	for _, t := range d.templates {
// 		d.processedTemplates[t.Name] = t.Template
// 	}
// 	return nil
// }

func (d *Destinator) Destination(kind primitive.Symbol, m *Medium) (string, error) {
	if string(kind) == "" {
		return "", fae.Errorf("kind is empty")
	}

	lib, ok := d.libraries[string(kind)]
	if !ok || lib == nil {
		return "", fae.Errorf("library not found for kind: %s", m.Kind)
	}

	if lib.Path == "" {
		return "", fae.Errorf("library path is empty for library: %s", lib.Name)
	}

	t, ok := d.templates[lib.LibraryTemplate.Name]
	if !ok || t == nil {
		return "", fae.Errorf("template not found for library: %s", lib.Name)
	}

	out := &strings.Builder{}
	data, err := NewDestinatorData(m)
	if err != nil {
		return "", fae.Wrap(err, "creating data")
	}
	data.path = lib.Path
	data.kind = string(kind)

	tmpl, err := template.New("destination").Parse(t.Template)
	if err != nil {
		return "", fae.Wrap(err, "parsing template")
	}
	err = tmpl.Execute(out, data)
	if err != nil {
		return "", fae.Wrap(err, "executing template")
	}

	return out.String(), nil
}

type DestinatorData struct {
	path      string
	kind      string
	directory string
	season    int
	episode   int
	absolute  int
}

func (d *DestinatorData) Path() string      { return d.path }
func (d *DestinatorData) Kind() string      { return d.kind }
func (d *DestinatorData) Directory() string { return d.directory }
func (d *DestinatorData) Season() string    { return fmt.Sprintf("%02d", d.season) }
func (d *DestinatorData) Episode() string   { return fmt.Sprintf("%02d", d.episode) }
func (d *DestinatorData) Absolute() string {
	if d.absolute == 0 {
		return fmt.Sprintf("%02dx%02d", d.season, d.episode)
	}
	return fmt.Sprintf("01x%03d", d.absolute)
}

func NewDestinatorData(m *Medium) (*DestinatorData, error) {
	d := &DestinatorData{
		directory: m.Directory,
		season:    m.SeasonNumber,
		episode:   m.EpisodeNumber,
		absolute:  m.AbsoluteNumber,
	}
	if m.Type == "Episode" {
		s := &Series{}
		err := app.DB.Series.FindByID(m.SeriesID, s)
		if err != nil {
			return nil, fae.Wrap(err, "finding series")
		}
		d.directory = s.Directory
	}

	return d, nil
}
