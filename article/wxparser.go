package article

import (
	"encoding/base64"
	"errors"
	"io"
	"reflect"
	"strconv"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"mutong.moe/go/utils/log"
)

var (
	errIncompleteWxArticle = errors.New("incomplete WeChat article")
)

const (
	_ConsumeIdle = iota
	_ConsumeInAuthor
	_ConsumeInBody
	_ConsumeDone
)

// NewFromWxStream consumes a HTML stream and generate a WxArticle
func NewFromWxStream(stream io.Reader) (*WxArticle, error) {
	// Initialize tokenizer
	tkz := html.NewTokenizer(stream)
	stage := _ConsumeIdle
	tags := []atom.Atom{0}
	atc := &WxArticle{}
	for stage != _ConsumeDone {
		switch tkz.Next() {
		case html.ErrorToken:
			if err := tkz.Err(); err == io.EOF {
				stage = _ConsumeDone
			} else {
				return nil, err
			}
		case html.StartTagToken:
			tagName, hasAttr := tkz.TagName()
			tags = append(tags, atom.Lookup(tagName))
			// Parse tag attrs
			if tags[len(tags)-1] == atom.Div {
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
				if tags[len(tags)-1] == atom.Script {
					consumeScript(atc, tkz.Text())
				}
			case _ConsumeInAuthor:
				atc.AuthorName = string(tkz.Text())
			}
		case html.EndTagToken:
			switch stage {
			case _ConsumeInAuthor:
				stage = _ConsumeIdle
			}
			tags = tags[:len(tags)-1]
		}
	}
	if len(atc.jsVarUnfilled) != 0 || atc.jsVarUnfilled == nil {
		log.Error("NewFromWxStream", "Unfilled variables: %v", atc.jsVarUnfilled)
		return nil, errIncompleteWxArticle
	}
	return atc, nil
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
		} else if err == errNotStringVariable {
			continue
		} else if err != nil {
			log.Error("consumeScript", "Unknwon error: %v", err)
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
