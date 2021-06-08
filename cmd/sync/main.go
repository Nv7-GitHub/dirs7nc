package main

import (
	dirsync "github.com/Nv7-Github/dirs7nc"
)

func main() {
	var prog *dirsync.Progress

	err := dirsync.Sync("testing/a", "testing/b", &prog)
	if err != nil {
		panic(err)
	}
}
