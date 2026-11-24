// Package migrations embeds the SQL migration files into the binary,
// so the API can apply them on startup without any external files.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
