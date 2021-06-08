package dirsync

// Sync syncs two directories, and returns a progress indicator
func Sync(srcdir string, dstdir string) (*Progress, chan error) {
	p := NewProgress(1)
	errs := make(chan error)

	go func() {
		// Sync src -> dst
		srcdst(srcdir, dstdir, p, errs)

		// Delete hanging files in dst
		dstsrc(dstdir, srcdir, p, errs)

		// Finish
		errs <- nil
	}()

	return p, errs
}
