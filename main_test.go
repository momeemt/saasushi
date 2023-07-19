package main

import "testing"

func TestAdd(*testing.T) {
	if add(1, 2) != 3 {
		panic("Something is terribly wrong.")
	}
}
