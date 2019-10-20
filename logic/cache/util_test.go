package cache

import (
	"fmt"
	"testing"
)

type A struct {
	A int
}

func TestGetHashAll(t *testing.T) {
	fmt.Println(getHashAll("2"))
}

func TestSetHash(t *testing.T) {
	fmt.Println(setHash("1", "1", A{A: 1}))
}

func TestGetHash(t *testing.T) {
	var a A
	fmt.Println(getHash("1", "1", &a))
	fmt.Println(a)
}

func TestDelHash(t *testing.T) {
	fmt.Println(delHash("1", "1"))
}
