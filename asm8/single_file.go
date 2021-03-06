package asm8

import (
	"bytes"
	"io"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// BuildSingleFile builds a package named "main" from a single file.
func BuildSingleFile(f string, rc io.ReadCloser) ([]byte, []*lex8.Error) {
	pinfo := build8.SimplePkg("_", f, rc)

	compiled, es := Lang().Compile(pinfo)
	if es != nil {
		return nil, es
	}

	buf := new(bytes.Buffer)
	e := link8.LinkMain(compiled.Lib(), buf, "main")
	if e != nil {
		return nil, lex8.SingleErr(e)
	}

	return buf.Bytes(), nil
}
