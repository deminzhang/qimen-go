package assets

import (
	"embed"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

//go:embed **/*
var FS embed.FS
