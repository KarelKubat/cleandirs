// Package l is a wrapper around log.
package l

import (
	"fmt"
	"log"
	"os"
)

type Purpose int

const (
	INFO        Purpose = iota // General info statement
	RECENT                     // File is recent, won't remove
	STALE                      // File is stale, will remove
	NOT_REGULAR                // File is not a regular file, won't remove
	WARN                       //  General warning
	FATAL                      // General failure
)

func (p Purpose) String() string {
	return []string{
		"INFO        ",
		"RECENT      ",
		"STALE       ",
		"NOT_REGULAR ",
		"WARN        ",
		"FATAL       "}[p]
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
}

func Printf(p Purpose, f string, args ...interface{}) {
	log.Print(p.String(), fmt.Sprintf(f, args...))
	if p == FATAL {
		os.Exit(1)
	}
}
