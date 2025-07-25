package assets

import "embed"

//go:embed css/* img/* fonts/*
var Assets embed.FS
