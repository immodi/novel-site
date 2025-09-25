package indexdtostructs

import "github.com/a-h/templ"

type LayoutData struct {
	Headers []HeaderItem
}

func (h HeaderItem) String() string {
	return h.Name
}

type HeaderItem struct {
	Name string
	Url  string
	Icon templ.Component
}
