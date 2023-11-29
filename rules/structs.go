package rules

import (
	"fmt"
	"regexp"

	"github.com/cespare/xxhash/v2"
)

type Method string

const (
	Hash   Method = "Hash"
	Redact Method = "Redact"
)

type Rule struct {
	Regex     *regexp.Regexp    `json:"regex"`
	Validator func(string) bool `json:"validator"`
	Method    Method            `json:"method"`
	Prefix    string            `json:"prefix"`
}

func (r *Rule) Validate(s string) bool {
	if r.Validator == nil {
		return true
	}

	return r.Validator(s)
}

func (r *Rule) Anonymize(s string) string {
	switch r.Method {
	case Hash:
		return r.Regex.ReplaceAllStringFunc(s, func(m string) string {
			if r.Validate(m) {
				if len(r.Prefix) > 0 {
					return fmt.Sprintf("<%s:%x>", r.Prefix, xxhash.New().Sum([]byte(m)))
				}

				return fmt.Sprintf("<%x>", xxhash.New().Sum([]byte(m)))
			}

			return m
		})
	case Redact:
		return r.Regex.ReplaceAllStringFunc(s, func(m string) string {
			if r.Validate(m) {
				if len(r.Prefix) > 0 {
					return fmt.Sprintf("<%s:%s>", r.Prefix, "*REDACTED*")
				}

				return fmt.Sprintf("<%s>", "*REDACTED*")
			}

			return m
		})
	default:
		return s
	}
}
