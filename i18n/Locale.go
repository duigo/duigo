package i18n

import (
	"fmt"
	"strings"
)

type Locale struct {
	Country  string
	Language string
}

func NewLocale(s string) Locale {
	v := strings.Split(s, "_")
	if len(v) >= 2 {
		return Locale{v[0], v[1]}
	}

	if len(v) >= 1 {
		return Locale{v[0], ""}
	}

	return Locale{}
}

func (l Locale) String() string {
	if "" == l.Language {
		return l.Country
	}

	return fmt.Sprintf("%s_%s", l.Country, l.Language)
}
