package dir

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Entry struct {
	Fullname string
	DirEntry os.DirEntry
	FileInfo os.FileInfo
	Linkdest *Entry
	Age      time.Duration
}

func List(dir string) (entries []Entry, err error) {
	now := time.Now()

	dirf, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("error opening dir %q: %v", dir, err)
	}
	defer dirf.Close()

	var dirEntries []os.DirEntry
	for {
		e, err := dirf.ReadDir(128)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading dir %q: %v\n", dir, err)
		}
		dirEntries = append(dirEntries, e...)
	}

	for _, e := range dirEntries {
		entry := Entry{
			Fullname: dir + "/" + e.Name(),
			DirEntry: e,
		}

		info, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("error getting info for %q: %v\n", entry.Fullname, err)
		}
		entry.FileInfo = info
		entry.Age = now.Sub(info.ModTime())

		entries = append(entries, entry)
	}
	return entries, nil
}
