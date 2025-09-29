package indexdtostructs

import "github.com/a-h/templ"

type LayoutData struct {
	Headers []HeaderItem
}

type HeaderItem struct {
	NavigationLink templ.Component
}
