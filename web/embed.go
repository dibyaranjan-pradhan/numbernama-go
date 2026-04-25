package web

import "embed"

// Assets holds the single-page client for the Go gameplay service.
//
//go:embed index.html app.js
var Assets embed.FS
