package ratelimit

type List map[string]bool

func (l List) Add(uri string) {
	l[uri] = true
}

func (l List) HasPrefix(uri string) {
	l[uri] = false
}

func (l List) Remove(uri string) {
	delete(l, uri)
}

var WhiteList List = make(List)
