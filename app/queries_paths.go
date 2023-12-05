package app

import "fmt"

func (p *Path) LocalPath() string {
	dir := cfg.DirectoriesCompleted
	if p.IsImage() {
		dir = cfg.DirectoriesImages
	}
	return fmt.Sprintf("%s/%s.%s", dir, p.Local, p.Extension)
}

func (p *Path) Exists() bool {
	return exists(p.LocalPath())
}

func (p *Path) IsImage() bool {
	return p.Type == "cover" || p.Type == "background"
}
func (p *Path) IsVideo() bool {
	return p.Type == "video"
}
