package web

import (
	_ "embed"
)

var (
	//go:embed login/index.html
	LoginHTML []byte
	//go:embed index.html
	HomeHTML []byte
)
