package anonymizer

import (
	"bytes"
	"io"
	"strings"

	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/rules"
)

func AnonymizePayload(r io.Reader) io.Reader {
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, r); err != nil {
		return nil
	}

	s := buf.String()

	return strings.NewReader(rules.Anonymize(s))
}
