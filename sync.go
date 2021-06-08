package dirsync

// Sync syncs two directories, and returns a progress indicator
func Sync(srcdir string, dstdir string, prog **Progress) error {
	p := NewProgress(1)
	*prog = p

	// Sync src -> dst
	err := srcdst(srcdir, dstdir, p)
	if err != nil {
		return err
	}

	return nil
}
