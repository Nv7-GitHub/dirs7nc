package main

import (
	"fmt"

	dirsync "github.com/Nv7-Github/dirs7nc"
)

func main() {
	_, errs := dirsync.Sync("/Users/nishant", "/Volumes/Files Sync/nishant")
	//_, errs := dirsync.Sync("testing/a", "testing/b")

	err := <-errs
	for err != nil {
		fmt.Println(err)
		err = <-errs
	}
}
