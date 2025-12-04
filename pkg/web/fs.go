package web

import _ "embed"

//go:embed index.html
var LandingPage []byte

//go:embed health.html
var HealthPage []byte
