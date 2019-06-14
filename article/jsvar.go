package article

import (
	"bytes"
	"errors"
	"io"
	"regexp"
)

type naiveJS bytes.Reader

var (
	_varDefinition = regexp.MustCompile(`var\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*`)
)

func newNaiveJS(buf []byte) *naiveJS {
	return (*naiveJS)(bytes.NewReader(buf))
}

func matchRegExp(re *regexp.Regexp, buf *bytes.Reader) ([][]byte, error) {
	// Fina a match
	startLen := buf.Len()
	matchLoc := re.FindReaderSubmatchIndex(buf)
	if matchLoc == nil {
		return nil, io.EOF
	}
	endLen := buf.Len()
	// Seek back to match start
	if _, err := buf.Seek(int64(endLen-startLen+matchLoc[0]), io.SeekCurrent); err != nil {
		return nil, err
	}
	// Retrieve all results
	result := make([][]byte, len(matchLoc)/2)
	for i := 0; i < len(result); i++ {
		// Seek to match start
		matchSize := matchLoc[i*2+1] - matchLoc[i*2]
		currentLen := buf.Len()
		if _, err := buf.Seek(int64(currentLen-startLen+matchLoc[i*2]), io.SeekCurrent); err != nil {
			return nil, err
		}
		result[i] = make([]byte, matchSize)
		if _, err := buf.Read(result[i]); err != nil {
			return nil, err
		}
	}
	// Seek to match end
	currentLen := buf.Len()
	if _, err := buf.Seek(int64(currentLen-startLen+matchLoc[1]), io.SeekCurrent); err != nil {
		return nil, err
	}
	return result, nil
}

func (js *naiveJS) nextVariable() (string, string, error) {
	// Find match of variable definition
	buf := (*bytes.Reader)(js)
	// Fine a variable definition
	matches, err := matchRegExp(_varDefinition, buf)
	if err != nil {
		return "", "", err
	}
	varName := string(matches[1])
	// Match variable value (Currently only strings are supported)
	var varValue []rune
	for len(varValue) == 0 {
		var r, quote rune
		var escaped bool
		for r, _, err = buf.ReadRune(); ; r, _, err = buf.ReadRune() {
			if err != nil {
				return "", "", err
			}
			if quote == '\x00' {
				if r == '"' || r == '\'' {
					quote = r
					varValue = make([]rune,0)
				} else if r != ' ' && r != '|' { // Naively skip spaces and || sign
					return "", "", errors.New("not a string variable")
				}
			} else {
				if escaped {
					escaped = false
				} else if r == quote {
					break
				} else if r == '\\' {
					escaped = true
				}
				varValue = append(varValue, r)
			}
		}
	}
	return varName, string(varValue), nil
}
