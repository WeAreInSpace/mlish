package mlish

import (
	"regexp"
)

func NewString() *String {
	return new(String)
}

type Querier interface {
	Find(s string) string
	Replace(src []byte, repl []byte) []byte
}

type String string

func (s *String) Query(pattern string) *QueryString {
	regex := regexp.MustCompile(pattern)
	return &QueryString{
		data:  s,
		regex: regex,
	}
}

func (s *String) Get() string {
	return string(*s)
}

type QueryString struct {
	data  *String
	regex *regexp.Regexp
}

func (qs QueryString) Replace(replaceTo string) {
	qs.regex.ReplaceAllString(qs.data.Get(), replaceTo)
}
