package life

import (
	"math/rand"
)

const (
	alive = true
	dead  = false
)

type Life [][]bool

func NewLife(x, y int) *Life {
	l := make(Life, x)
	for i := range l {
		l[i] = make([]bool, y)
	}
	return &l
}

func (l *Life) Clear() {
	for i := range *l {
		for j := range (*l)[i] {
			(*l)[i][j] = dead
		}
	}
}

func (l *Life) Fill() {
	for i := range *l {
		for j := range (*l)[i] {
			(*l)[i][j] = alive
		}
	}
}

func (l *Life) Rand() {
	for i := range *l {
		for j := range (*l)[i] {
			if rand.Intn(6) == 0 {
				(*l)[i][j] = alive
			} else {
				(*l)[i][j] = dead
			}
		}
	}
}

func (l *Life) Next() {
	next := make([][]bool, len(*l))
	for i := range next {
		next[i] = make([]bool, len((*l)[i]))
		for j := range (*l)[i] {
			next[i][j] = (*l).Check(i, j)
		}
	}
	*l = next
}

func (l *Life) Check(x, y int) bool {
	w := len(*l)
	h := len((*l)[0])
	state := func(x, y int) bool {
		x += w
		x %= w
		y += h
		y %= h
		return (*l)[x][y]
	}
	population := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if !(j == 0 && i == 0) && state(x+i, y+j) {
				population++
			}
		}
	}
	return population == 3 || population == 2 && (*l)[x][y]
}
