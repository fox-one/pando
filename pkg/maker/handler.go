package maker

import (
	"fmt"
	"sort"
)

type HandlerFunc func(r *Request) error

type handlerWithVersion struct {
	f HandlerFunc
	v int
}

func HandlerVersion(m map[int]HandlerFunc) HandlerFunc {
	versions := make([]handlerWithVersion, 0, len(m))
	for v, f := range m {
		versions = append(versions, handlerWithVersion{f: f, v: v})
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].v > versions[j].v
	})

	return func(r *Request) error {
		for _, h := range versions {
			if r.SysVersion >= h.v {
				return h.f(r)
			}
		}

		return fmt.Errorf("%s request with version %d cannot be handled", r.Action.String(), r.SysVersion)
	}
}
