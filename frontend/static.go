package frontend

import (
	"embed"
	"io/fs"
	"log"
)

// NOTE: go:embed でフロントエンドをGoのバイナリに埋め込んで配信している
// 同様にembed.FSを増やすことで複数のフロントエンドプロジェクトを同時に埋め込むことが可能

//go:embed placeholder/**
var uiDist embed.FS

var UI fs.FS

func init() {
	// Prefer built UI (if embedded); fallback to placeholder
	if sub, err := fs.Sub(uiDist, "app-ui/dist"); err == nil {
		UI = sub
		return
	}
	ph, err := fs.Sub(uiDist, "placeholder")
	if err != nil {
		log.Printf("frontend: placeholder fs missing: %v", err)
	}
	UI = ph
}
