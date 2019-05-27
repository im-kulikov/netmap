package netmap

import (
	"sort"
)

type nw struct {
	n []uint32
	w []uint64
}

func (f nw) Len() int { return len(f.n) }
func (f nw) Swap(i, j int) {
	f.n[i], f.n[j] = f.n[j], f.n[i]
	f.w[i], f.w[j] = f.w[j], f.w[i]
}
func (f nw) Less(i, j int) bool {
	return f.n[i] < f.n[j]
}
func (f nw) Load() ([]uint32, []uint64) {
	return f.n, f.w
}

func getNodes(b Bucket, path []Bucket) (nodes []uint32) {
	if len(path) == 0 {
		return b.Nodelist()
	}
	for _, p := range b.Children() {
		if p.Equals(path[0]) {
			return getNodes(p, path[1:])
		}
	}
	return nil
}

func contains(nodes []uint32, n uint32) bool {
	for _, i := range nodes {
		if i == n {
			return true
		}
	}
	return false
}

func intersect(a, b []uint32) []uint32 {
	if a == nil {
		return b
	}

	var (
		la, lb = len(a), len(b)
		l      = min(la, lb)
		c      = make([]uint32, 0, l)
	)

	for i, j := 0, 0; i < la && j < lb; {
		switch true {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			c = append(c, a[i])
			i++
			j++
		}
	}

	return c
}

func diff(a []uint32, b map[uint32]struct{}, w []uint64) (c []uint32, d []uint64) {
	c = make([]uint32, 0, len(a))
	d = make([]uint64, 0, len(a))
	for i, e := range a {
		if _, ok := b[e]; !ok {
			c = append(c, e)
			d = append(d, w[i])

		}
	}
	return
}

func union(a, b []uint32) []uint32 {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	var (
		la, lb = len(a), len(b)
		l      = la + lb
		c      = make([]uint32, 0, l)
		i, j   int
	)

	for i, j = 0, 0; i < la && j < lb; {
		switch true {
		case a[i] < b[j]:
			c = append(c, a[i])
			i++
		case a[i] > b[j]:
			c = append(c, b[j])
			j++
		default:
			c = append(c, a[i])
			i++
			j++
		}
	}

	if i == la {
		c = append(c, b[j:]...)
	} else {
		c = append(c, a[i:]...)
	}

	return c
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func strawSort(a []uint32, b []uint64) {
	s := nw{a, b}
	sort.Sort(s)
	a, b = s.Load()
}
