package app

import (
	"fmt"
	"regexp"

	"github.com/samber/lo"
)

var pathTagRegex = regexp.MustCompile(`(?i)\s+\[(\w+)\]\.\w+$`)

func (p *Path) ParseTag() {
	matches := pathTagRegex.FindStringSubmatch(p.Local)
	if len(matches) > 1 {
		p.Tag = matches[1]
	}
}

func (p *Path) LocalPath() string {
	dir := app.Config.DirectoriesCompleted
	if p.IsCoverBackground() {
		dir = app.Config.DirectoriesImages
	}
	return fmt.Sprintf("%s/%s.%s", dir, p.Local, p.Extension)
}

func (p *Path) Exists() bool {
	return exists(p.LocalPath())
}

func (p *Path) IsCoverBackground() bool {
	return p.Type == "cover" || p.Type == "background" || lo.Contains(app.Config.ExtensionsImages, p.Extension)
}
func (p *Path) IsVideo() bool {
	return p.Type == "video"
}
