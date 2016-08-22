package main

/*
 * dumpbolt.go
 * Print the contents of a bolt database
 * By J. Stuart McMurray
 * Created 20160821
 * Last Modified 20160822
 */

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func main() {
	var (
		pathsep = flag.String(
			"p",
			"/",
			"Bucket `separator`",
		)
		printAllPaths = flag.Bool(
			"a",
			false,
			"Print a path for every key/value pair",
		)
		indentWidth = flag.Uint(
			"i",
			8,
			"Number of `spaces` to indent each bucket level",
		)
		openTimeout = flag.Duration(
			"t",
			time.Second,
			"Database open `timeout`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options] <boltdb> [path [path...]]

Prints the contents of boltdb, optionally starting at the given path(s), which
should be a list of nested buckets formatted as a unix path (use -p if a bucket
contains a "/").

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Make sure we at least have a database name */
	if 0 == flag.NArg() {
		flag.Usage()
		os.Exit(1)
	}

	/* Open the database */
	db, err := bolt.Open(flag.Arg(0), 0600, &bolt.Options{
		ReadOnly: true,
		Timeout:  *openTimeout,
	})
	if nil != err {
		fmt.Fprintf(
			os.Stderr,
			"Unable to open %v: %v",
			flag.Arg(0),
			err,
		)
		os.Exit(2)
	}

	/* Work out the starting path(s) */
	starts := make([]string, 0, 1)
	if 1 == flag.NArg() {
		starts = append(starts, *pathsep)
	} else {
		starts = append(starts, flag.Args()[1:]...)
	}

	/* Dump the database starting at each starting path */
	for _, start := range starts {
		if err := Dump(
			db,
			[]byte(start),
			[]byte(*pathsep),
			*printAllPaths,
			*indentWidth,
		); nil != err {
			log.Printf("Error dumping from %v: %v", start, err)
		}
	}
}
