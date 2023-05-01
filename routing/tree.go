package routing

import (
	"fmt"
	"path"
)

type routeRoot map[string]routeBranch

type routeBranch map[string]NodeRender

func getBaseAndSuffix(urlPath string) (string, string) {

	dir := path.Dir(urlPath)
	base := path.Base(urlPath)

	return dir, base
}

func (rr routeRoot) insert(k string, val NodeRender) {
	var (
		base, suffix = getBaseAndSuffix(k)
	)
	rb, ok := rr[base]
	if !ok {
		rr[base] = make(routeBranch)
		rb = rr[base]

	}
	rb.insert(suffix, val)
}

func (rr routeRoot) update(k string, val NodeRender) {
	var (
		base, suffix = getBaseAndSuffix(k)
	)
	rb, ok := rr[base]
	if !ok {
		rr[base] = make(routeBranch)
		rb = rr[base]
	}
	rb.update(suffix, val)
}

func (rr routeRoot) remove(k string) {
	var (
		base, suffix = getBaseAndSuffix(k)
	)
	rb, ok := rr[base]
	if !ok {
		return

	}
	rb.remove(suffix)
}

func (rr routeRoot) fetch(k string) (router NodeRender, err error) {
	var (
		base, suffix = getBaseAndSuffix(k)
	)
	rb, ok := rr[base]
	if !ok {
		rr[base] = make(routeBranch)
		return nil, fmt.Errorf("path not found %s", base)
	}
	router, err = rb.fetch(suffix)
	return router, err
}

func (rb routeBranch) insert(k string, val NodeRender) {
	if k == "" {
		println("invalid k")
		return
	}
	_, ok := rb[k]
	if ok {
		println("already found router at", k)
		return
	}
	rb[k] = val
}

func (rb routeBranch) update(k string, val NodeRender) {
	if k == "" {
		println("invalid k")
		return
	}
	rb[k] = val
}

func (rb routeBranch) remove(k string) {
	delete(rb, k)
}

func (rb routeBranch) fetch(k string) (r NodeRender, err error) {
	var (
		ok      = false
		hasWild = false
	)
	if k == "" {
		println("invalid k")
		return
	}
	r, ok = rb[k]
	if !ok {
		r, hasWild = rb[k]
		if !hasWild {
			err = fmt.Errorf("path not found %s", k)
		}
	}
	return r, err
}
