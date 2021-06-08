package dirsync

import (
	"os"
	"path/filepath"
)

func dstsrc(srcdir string, dstdir string, prog *Progress) error {
	src, err := os.ReadDir(srcdir)
	if err != nil {
		return err
	}

	prog.AddTotal(int64(len(src)))

	// Go through files in directory
	for _, file := range src {
		// Check if file exists in corresponding folder
		_, err := os.Stat(filepath.Join(dstdir, file.Name()))
		if os.IsNotExist(err) {
			// If not, remove it
			f := filepath.Join(srcdir, file.Name())
			if file.IsDir() {
				err = os.RemoveAll(f)
				if err != nil {
					return err
				}
			} else {
				err = os.Remove(f)
				if err != nil {
					return err
				}
			}
		}
		prog.Add(1)
	}

	return nil
}
