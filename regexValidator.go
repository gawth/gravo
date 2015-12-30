package main

import (
	"regexp"
)

type regexValidator struct {
	reg *regexp.Regexp
}

func (rv *regexValidator) IsValid(target []byte) bool {
	return rv.reg.Match(target)
}
