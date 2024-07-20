package pkg

import "embed"

//go:embed httputil
var HttpServerFs embed.FS

//go:embed database
var DatabaseFs embed.FS

//go:embed identity
var IdentityFs embed.FS

//go:embed cmd
var CmdFs embed.FS
