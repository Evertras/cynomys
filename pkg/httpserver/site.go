package httpserver

import "embed"

//go:embed site/*
var siteFiles embed.FS
