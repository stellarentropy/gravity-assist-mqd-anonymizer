package rules

import (
	"regexp"

	"github.com/paemuri/brdoc"
)

var CNPJRule = &Rule{
	Regex:     regexp.MustCompile(`\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}`),
	Method:    Hash,
	Prefix:    "CNPJ",
	Validator: brdoc.IsCNPJ,
}
