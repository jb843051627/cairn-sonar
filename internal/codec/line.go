package codec

import (
	"bufio"
	"io"
	"strings"
)

func Lines(r io.Reader) ([]string, error) {
	scan := bufio.NewScanner(r)
	out := make([]string, 0)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line != "" {
			out = append(out, line)
		}
	}
	return out, scan.Err()
}
