package article

import (
	"io"
	"reflect"
	"strconv"
	"encoding/base64"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	_ConsumeIdle = iota
	_ConsumeInAuthor
	_ConsumeInBody
	_ConsumeDone
)

// NewFromWxStream consumes a HTML stream and generate a WxArticle
func NewFromWxStream(stream io.Reader) (*WxArticle, error) {
	tkz := html.NewTokenizer(stream)
	atc := WxArticle{}
	stage := _ConsumeIdle
	tags := []atom.Atom{0}
	for stage != _ConsumeDone {
		switch tkz.Next() {
		case html.ErrorToken:
			if tkz.Err() == io.EOF {
				stage = _ConsumeDone
			} else {
				return nil, tkz.Err()
			}
		case html.StartTagToken:
			tagName, hasAttr := tkz.TagName()
			tags = append(tags, atom.Lookup(tagName))
			// Parse tag attrs
			if tags[len(tags) - 1] == atom.Div {
				var key, value []byte // define them first since we don't want to redefine hasAttr
				for hasAttr {
					key, value, hasAttr = tkz.TagAttr()
					switch string(key) {
					case "class":
						switch string(value) {
						case "reward-author":
							stage = _ConsumeInAuthor
						}
					}
				}
			}
		case html.TextToken:
			switch stage {
			case _ConsumeIdle:
				if tags[len(tags) - 1] == atom.Script {
					consumeScript(&atc, tkz.Text())
				}
			case _ConsumeInAuthor:
				atc.AuthorName = string(tkz.Text())
			}
		case html.EndTagToken:
			switch stage {
			case _ConsumeInAuthor:
				stage = _ConsumeIdle
			}
			tags = tags[:len(tags) - 1]
		}
	}
	return &atc, nil
}

func consumeScript(atc *WxArticle, script []byte) {
	atcType := reflect.TypeOf(*atc)
	// Build jsVarUnfilled
	if atc.jsVarUnfilled == nil {
		atc.jsVarUnfilled = make(map[string]int)
		for i := 0; i < atcType.NumField(); i++ {
			if varName := atcType.Field(i).Tag.Get("jsvar"); varName != "" {
				atc.jsVarUnfilled[varName] = i
			}
		}
	}
	// Scan script string for variable definitions
	actValue := reflect.Indirect(reflect.ValueOf(atc))
	buffer := newNaiveJS(script)
	for {
		varName, varValue, err := buffer.nextVariable()
		if err == io.EOF {
			break
		}
		if fieldID, ok := atc.jsVarUnfilled[varName]; ok {
			typeField := atcType.Field(fieldID)
			// Decode the field first
			if fieldEncoding := typeField.Tag.Get("encoding"); fieldEncoding == "base64" {
				val, _ := base64.StdEncoding.DecodeString(varValue)
				varValue = string(val)
			}
			// Assign to variable
			if fieldKind := typeField.Type.Kind(); fieldKind == reflect.String {
				actValue.Field(fieldID).SetString(varValue)
			} else if fieldKind == reflect.Int64 {
				varNumValue, _ := strconv.ParseInt(varValue, 10, 64)
				actValue.Field(fieldID).SetInt(varNumValue)
			}
			delete(atc.jsVarUnfilled, varName)
		}
	}
}
