package httpserver

import "embed"

//go:embed site/*
var siteFilesRaw embed.FS
