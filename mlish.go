package mlish

import (
	"io"
	"os"
)

type Setting struct {
	DebugMode bool
	Out       io.Writer
}

var Settings = &Setting{
	Out: os.Stdout,
}
