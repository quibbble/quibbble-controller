package parse

import "embed"

//go:embed items/*.yaml
//go:embed spells/*.yaml
//go:embed units/*.yaml
var Library embed.FS
