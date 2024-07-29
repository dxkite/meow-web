package pkg

import "embed"

//go:embed httputil
var HttpUtilFs embed.FS

//go:embed database
var DatabaseFs embed.FS

//go:embed crypto
var CryptoFs embed.FS

//go:embed cmd
var CmdFs embed.FS

//go:embed errors
var ErrorsFs embed.FS

//go:embed log
var LogFs embed.FS
