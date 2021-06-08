package dirsync

import "fmt"

type empty struct{}

// Sync syncs two directories, and returns a progress indicator
func Sync(srcdir string, dstdir string, limit int) (*Progress, chan error) {
	p := NewProgress(1)
	errs := make(chan error)

	haslimit := false
	var maxprocs chan empty
	if limit > 0 {
		maxprocs = make(chan empty, limit)
		haslimit = true
	}

	go func() {
		// Sync src -> dst
		srcdst(srcdir, dstdir, p, errs, haslimit, maxprocs)

		// Print for debugging
		fmt.Println("dstsrc")

		// Delete hanging files in dst
		dstsrc(dstdir, srcdir, p, errs, haslimit, maxprocs)

		// Finish
		errs <- nil
	}()

	return p, errs
}
