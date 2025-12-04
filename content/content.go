package content

import "embed"

//go:embed blog/*
var Content embed.FS
