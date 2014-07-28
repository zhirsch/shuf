package main

import (
	"math/rand"
)

type numbers []int

func (n numbers) shuffle(*rand.Rand) {
	for i := 0; i < len(n); i++ {
		// j := r.Intn(i + 1)
		j := rand.Intn(i + 1)
		if i != j {
			n[i], n[j] = n[j], n[i]
		}
	}
}

func (n numbers) sorted() bool {
	for i, j := range n {
		if i != j {
			return false
		}
	}
	return true
}
