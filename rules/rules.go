package rules

var Rules []*Rule

func init() {
	Rules = append(Rules, CPFRule)
	Rules = append(Rules, CNPJRule)
}

func Anonymize(s string) string {
	for _, r := range Rules {
		s = r.Anonymize(s)
	}

	return s
}
