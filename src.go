package dirsync

import (
	"io"
	"os"
	"path/filepath"

	"github.com/djherbis/times"
)

func srcdst(srcdir string, dstdir string, prog *Progress) error {
	src, err := os.ReadDir(srcdir)
	if err != nil {
		return err
	}

	prog.AddTotal(int64(len(src)))
	res := make([]chan error, len(src))

	// Go through files in directory
	for i, file := range src {
		// Get info
		errout := make(chan error)
		res[i] = errout

		srcname := filepath.Join(srcdir, file.Name())
		dstname := filepath.Join(dstdir, file.Name())

		info, err := file.Info()
		if err != nil {
			return err
		}

		// Directory or File?
		if file.IsDir() {
			// Make Destination Folder
			err = os.MkdirAll(dstname, info.Mode())
			if err != nil {
				return err
			}

			// Recursively do it on that folder
			go func(out chan error, srcdir string, dstdir string, prog *Progress) {
				err := srcdst(srcdir, dstdir, prog)
				out <- err
			}(res[i], srcname, dstname, prog)
		} else {
			// Are Equal? If so, don't copy
			stat, err := os.Stat(dstname)
			exist := os.IsExist(err)
			if exist {
				if err != nil {
					return err
				}
				if stat.Mode() == info.Mode() && stat.ModTime().Equal(info.ModTime()) && stat.Size() == info.Size() {
					go func() {
						errout <- nil
					}()
					continue
				}
			}

			// If Not, Copy
			sf, err := os.Open(srcname)
			if err != nil {
				return err
			}
			df, err := os.OpenFile(dstname, os.O_CREATE|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}

			// Copy files async
			go func(sf, df *os.File, errout chan error) {
				err = df.Truncate(0)
				if err != nil {
					errout <- err
					return
				}

				_, err = io.Copy(df, sf)
				if err != nil {
					errout <- err
					return
				}
				df.Close()
				sf.Close()

				// Change modified time for file to that of original
				inf, err := times.Stat(srcname)
				if err != nil {
					errout <- err
					return
				}
				err = os.Chtimes(dstname, inf.AccessTime(), inf.ModTime())
				if err != nil {
					errout <- err
					return
				}

				errout <- nil
			}(sf, df, errout)
		}
	}

	// Wait for it to be complete
	for _, out := range res {
		err := <-out
		prog.Add(1)
		if err != nil {
			return err
		}
	}

	return nil
}
