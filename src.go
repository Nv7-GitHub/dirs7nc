package dirsync

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/djherbis/times"
)

func srcdst(srcdir string, dstdir string, prog *Progress, errs chan error) {
	src, err := os.ReadDir(srcdir)
	if err != nil {
		errs <- err
		return
	}

	prog.AddTotal(int64(len(src)))
	wg := &sync.WaitGroup{}

	// Go through files in directory
	for _, file := range src {
		srcname := filepath.Join(srcdir, file.Name())
		dstname := filepath.Join(dstdir, file.Name())

		info, err := file.Info()
		if err != nil {
			errs <- err
		}

		// Directory or File?
		if file.IsDir() {
			// Make Destination Folder
			err = os.MkdirAll(dstname, info.Mode())
			if err != nil {
				errs <- err
			}

			// Recursively do it on that folder
			wg.Add(1)
			go func(srcdir string, dstdir string, prog *Progress) {
				srcdst(srcdir, dstdir, prog, errs)
				prog.Add(1)
				wg.Done()
			}(srcname, dstname, prog)
		} else {
			// Are Equal? If so, don't copy
			stat, err := os.Stat(dstname)
			exist := os.IsExist(err)
			if exist {
				if err != nil {
					errs <- err
					continue
				} else if stat.Mode() == info.Mode() && stat.ModTime().Equal(info.ModTime()) && stat.Size() == info.Size() {
					continue
				}
			}

			// If Not, Copy
			sf, err := os.Open(srcname)
			if err != nil {
				errs <- err
			}
			df, err := os.OpenFile(dstname, os.O_CREATE|os.O_WRONLY, info.Mode())
			if err != nil {
				errs <- err
			}

			// Copy files async
			wg.Add(1)
			go func(sf, df *os.File) {
				err = df.Truncate(0)
				if err != nil {
					errs <- err
					return
				}

				_, err = io.Copy(df, sf)
				if err != nil {
					errs <- err
					return
				}
				df.Close()
				sf.Close()

				// Change modified time for file to that of original
				inf, err := times.Stat(srcname)
				if err != nil {
					errs <- err
					return
				}
				err = os.Chtimes(dstname, inf.AccessTime(), inf.ModTime())
				if err != nil {
					errs <- err
					return
				}
				prog.Add(1)
				wg.Done()
			}(sf, df)
		}
	}

	wg.Wait()
}
