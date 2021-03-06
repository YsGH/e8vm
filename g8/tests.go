package g8

import (
	"strings"

	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

func isTestName(name string) bool {
	if len(name) < len("TestX") {
		return false
	}
	if !strings.HasPrefix(name, "Test") {
		return false
	}
	lead := name[4]
	if lead >= 'a' && lead <= 'z' {
		return false
	}
	return true
}

func listTests(tops *sym8.Table) []*objFunc {
	var list []*objFunc

	syms := tops.List()
	for _, s := range syms {
		if s.Type != symFunc {
			continue
		}
		f := s.Item.(*objFunc)
		if f.isMethod {
			panic("bug") // a top level function should never be a method
		}
		if !types.SameType(f.ref.Type(), types.VoidFunc) {
			continue
		}
		name := s.Name()
		if isTestName(name) {
			list = append(list, f)
		}
	}

	return list
}
