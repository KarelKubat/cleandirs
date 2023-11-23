package main

import (
	"flag"
	"strings"
	"syscall"
	"time"

	"github.com/KarelKubat/cleandirs/dir"
	"github.com/KarelKubat/cleandirs/l"
	"github.com/KarelKubat/flagnames"
)

var (
	dirsToClean = flag.String("dirs-to-clean", "/tmp", "comma-separated list of directories to clean out")
	ttl         = flag.Duration("ttl", time.Hour*24, "time to live: files older than this are removed")
	pruneDirs   = flag.Bool("prune-dirs", true, "when true, empty directories under --dirs-to-clean are removed")
	allFiles    = flag.Bool("all-files", false, "when true, remove also special files (pipes, fifo's etc.), default: regular files only")
	dryRun      = flag.Bool("dry-run", true, "when true, suggested removals are shown but not actuated")

	usageInfo = `
This is cleandirs, a tiny utility to remove stale files in temporary directories.
Try cleandirs -help for a listing of available flags.	
`
)

func main() {
	flagnames.Patch()
	flag.Parse()
	if len(flag.Args()) != 0 {
		l.Printf(l.FATAL, usageInfo)
	}

	cutoff := time.Now().Add(-*ttl)
	for _, dir := range strings.Split(*dirsToClean, ",") {
		cleanFilesIn(dir, cutoff)
	}
}

func cleanFilesIn(d string, cutoff time.Time) {
	l.Printf(l.INFO, "scanning %q, files modified before %v are considered stale", d, cutoff)
	entries, err := dir.List(d)
	if err != nil {
		l.Printf(l.WARN, "%q: failed to list entries: %v\n", d, err)
		return
	}
	var subdirs []string
	for _, e := range entries {
		// Entry is a subdir: queue for further processing
		if e.DirEntry.IsDir() {
			subdirs = append(subdirs, e.Fullname)
			continue
		}
		// Entry is too young: skip
		if e.FileInfo.ModTime().After(cutoff) {
			l.Printf(l.RECENT, "%q: keeping, age: %v", e.Fullname, e.Age)
			continue
		}
		// Entry is not a regular file: skip unless --all-files is given
		if !e.FileInfo.Mode().IsRegular() && !*allFiles {
			l.Printf(l.NOT_REGULAR, "%q: not a file", e.Fullname)
			continue
		}

		l.Printf(l.STALE, "%q: stale, age: %v", e.Fullname, e.Age)
		if *dryRun {
			continue
		}
		if err := syscall.Unlink(e.Fullname); err != nil {
			l.Printf(l.WARN, "%q: failed to unlink: %v\n", e.Fullname, err)
		}
	}

	for _, subd := range subdirs {
		cleanFilesIn(subd, cutoff)
		if *dryRun || !*pruneDirs {
			continue
		}
		syscall.Rmdir(subd) // discard error, dir is not necessarily empty
	}
}
