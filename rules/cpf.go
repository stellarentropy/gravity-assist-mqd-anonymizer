package rules

import (
	"regexp"

	"github.com/paemuri/brdoc"
)

var CPFRule = &Rule{
	Regex:     regexp.MustCompile(`\d{3}\.?\d{3}\.?\d{3}-?\d{2}`),
	Method:    Hash,
	Prefix:    "CPF",
	Validator: brdoc.IsCPF,
}
