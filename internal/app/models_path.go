package app

import (
	"fmt"

	"github.com/samber/lo"
)

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
