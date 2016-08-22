package main

/*
 * print.go
 * Actually print the contents
 * By J. Stuart McMurray
 * Created 20160823
 * Last Modified 20160823
 */
import (
	"fmt"
	"strconv"
	"strings"
)

/* Print prints the contents of b */
func Print(
	b BIter,
	prefix, pathsep []byte,
	printAllPaths bool,
	indentWidth uint,
) error {
	c := b.Cursor()
	for k, v := c.First(); nil != k; k, v = c.Next() {
		/* Handle buckets */
		if nil == v {
			if err := printBucket(
				b,
				k,
				prefix,
				pathsep,
				printAllPaths,
				indentWidth,
			); nil != err {
				return err
			}
			continue
		}

		/* Handle keys */
		printSafe("%s%s -> %s", prefix, k, v)
	}
	return nil
}

/* printBucket handles printing a bucket's name and recursing */
func printBucket(
	b BIter,
	bname, prefix, pathsep []byte,
	printAllPaths bool,
	indentWidth uint,
) error {
	nb := b.Bucket(bname)
	if nil == nb {
		return fmt.Errorf("bucket %q not found", bname)
	}
	/* Print bucket name */
	bn := []byte(fmt.Sprintf("%s%s%s", prefix, bname, pathsep))
	printSafe("%s", bn)
	/* Work out next prefix */
	var np []byte
	if printAllPaths {
		np = bn
	} else {
		np = incSpace(prefix, indentWidth)
	}
	return Print(nb, np, pathsep, printAllPaths, indentWidth)
}

/* printSafe is like a safe printf which adds a trailing newline. */
func printSafe(f string, a ...interface{}) {
	/* Turn it into a string, add a newline if needed */
	s := fmt.Sprintf(f, a...)

	/* Safen it */
	s = strconv.QuoteToASCII(s)

	/* Remove surrounding quotes */
	s = s[1:]
	s = s[:len(s)-1]

	/* Add a trailing newline */
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}

	/* Send to stdout */
	fmt.Printf("%s", s)
}

/* incSpace adds n spaces to the end of b and returns the resulting slice */
func incSpace(b []byte, n uint) []byte {
	np := make([]byte, 0, uint(len(b))+n)
	np = append(np, b...)
	for i := uint(0); i < n; i++ {
		np = append(np, ' ')
	}
	return np
}
