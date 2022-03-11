package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/KarelKubat/cleandirs/dir"
)

var (
	dirsToClean = flag.String("dirs-to-clean", "/tmp", "comma-separated list of directories to clean out")
	ttl         = flag.Duration("ttl", time.Hour*24, "time to live: files older than this are removed")
	pruneDirs   = flag.Bool("prune-dirs", true, "when true, empty directories under --dirs-to-clean are removed")
	dryRun      = flag.Bool("dry-run", true, "when true, suggested removals are shown but not actuated")
	version     = flag.Bool("version", false, "show version ID, then stop")

	usageInfo = fmt.Sprintf(`

This is cleandirs V%s, a tiny utility to remove stale files in temporary directories.
Try cleandirs -help for a listing of available flags.	
	`, versionID)
)

const (
	versionID = "1.00"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatal(usageInfo)
	}

	if *version {
		fmt.Println(versionID)
		os.Exit(1)
	}

	cutoff := time.Now().Add(-*ttl)
	log.Printf("files modified before %v are considered stale", cutoff)
	for _, dir := range strings.Split(*dirsToClean, ",") {
		cleanFilesIn(dir, cutoff)
	}
}

func cleanFilesIn(d string, cutoff time.Time) {
	entries, err := dir.List(d)
	if err != nil {
		log.Printf("warning: failed to list entries in %q: %v\n", d, err)
	}
	var subdirs []string
	for _, e := range entries {
		if e.DirEntry.IsDir() {
			subdirs = append(subdirs, e.Fullname)
			continue
		}
		if e.FileInfo.ModTime().After(cutoff) {
			continue
		}
		log.Println(e.Fullname, e.Age)
		if *dryRun {
			continue
		}
		if err := syscall.Unlink(e.Fullname); err != nil {
			fmt.Printf("warning: failed to unlink %q: %v\n", e.Fullname, err)
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
