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

var (
	errNotStringVariable = errors.New("not a string variable")
)

func newNaiveJS(buf []byte) *naiveJS {
	return (*naiveJS)(bytes.NewReader(buf))
}

func matchRegExp(re *regexp.Regexp, buf *bytes.Reader) (result [][]byte, err error) {
	// Fina a match
	startLen := buf.Len()
	matchLoc := re.FindReaderSubmatchIndex(buf)
	if matchLoc == nil {
		return nil, io.EOF
	}
	endLen := buf.Len()
	// Seek back to match start
	if _, err = buf.Seek(int64(endLen-startLen+matchLoc[0]), io.SeekCurrent); err != nil {
		return
	}
	// Retrieve all results
	result = make([][]byte, len(matchLoc)/2)
	for i := 0; i < len(result); i++ {
		// Seek to match start
		matchSize := matchLoc[i*2+1] - matchLoc[i*2]
		currentLen := buf.Len()
		if _, err = buf.Seek(int64(currentLen-startLen+matchLoc[i*2]), io.SeekCurrent); err != nil {
			return
		}
		result[i] = make([]byte, matchSize)
		if _, err = buf.Read(result[i]); err != nil {
			return
		}
	}
	// Seek to match end
	currentLen := buf.Len()
	_, err = buf.Seek(int64(currentLen-startLen+matchLoc[1]), io.SeekCurrent)
	return
}

func (js *naiveJS) nextVariable() (varName string, varValue string, err error) {
	// Find a match of variable definition
	buf := (*bytes.Reader)(js)
	varNameMatch, err := matchRegExp(_varDefinition, buf)
	if err != nil {
		return
	}
	varName = string(varNameMatch[1])
	// Match variable value (Currently only strings are supported)
	var varValueRunes []rune
	for len(varValueRunes) == 0 {
		var r, quote rune
		var escaped bool
		for r, _, err = buf.ReadRune(); ; r, _, err = buf.ReadRune() {
			if err != nil {
				return
			}
			// not in a quote
			if quote == '\x00' {
				// quote handeling
				if r == '"' || r == '\'' {
					quote = r
					varValueRunes = make([]rune, 0) // reset result
				} else if r != ' ' && r != '|' { // Naively skip spaces and || sign
					if varValueRunes == nil { // nothing is found
						err = errNotStringVariable
					}
					return
				}
			} else {
				// escaping
				if escaped {
					escaped = false
				} else if r == quote {
					break
				} else if r == '\\' {
					escaped = true
				}
				varValueRunes = append(varValueRunes, r)
			}
		}
	}
	varValue = string(varValueRunes)
	return
}
