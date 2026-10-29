package assets

import (
	"embed"
)

// This line below tells Go to store files from the /static folder in an 
// *embedded filesystem* referenced by the glob var "Files"

//go:embed "static"
var Files embed.FS