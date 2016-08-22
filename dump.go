package main

/*
 * dump.go
 * Dump a bolt database
 * By J. Stuart McMurray
 * Created 20160823
 * Last Modified 20160823
 */

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
)

/* BIter is anything which has buckets and can be iterated over */
type BIter interface {
	Bucket([]byte) *bolt.Bucket
	Cursor() *bolt.Cursor
}

/* Dump starts a dump starting at start */
func Dump(
	db *bolt.DB,
	start, bsep []byte,
	printAllPaths bool,
	indentWidth uint,
) error {
	/* Work out the series of buckets */
	bs := bytes.Split(start, bsep)
	/* Remove blanks */
	cur := 0
	for _, b := range bs {
		if 0 == len(b) {
			continue
		}
		bs[cur] = b
		cur++
	}
	bs = bs[:cur]

	/* Remake path without //'s */
	start = []byte(fmt.Sprintf("%s%s", bsep, bytes.Join(bs, bsep)))

	var prefix []byte
	if printAllPaths {
		prefix = start
		if !bytes.HasSuffix(prefix, bsep) {
			prefix = append(start, []byte(bsep)...)
		}
	} else {
		prefix = []byte{}
	}

	/* Print out this bit */
	return db.View(func(tx *bolt.Tx) error {
		startBucket, err := Dive(tx, bs, []byte{}, bsep)
		if nil != err {
			return err
		}
		return Print(
			startBucket,
			prefix,
			bsep,
			printAllPaths,
			indentWidth,
		)
	})
}

/* Dive dives through layers of buckets until it runs out of buckets.  It
returns the last bucket in the list, if it exists. */
func Dive(
	h BIter,
	buckets [][]byte,
	where []byte,
	bsep []byte,
) (BIter, error) {
	/* If we've no more buckets, we're here */
	if 0 == len(buckets) {
		return h, nil
	}
	/* Path of bucket we want */
	tgt := bytes.Join([][]byte{where, buckets[0]}, bsep)
	/* Get the next bucket in the list */
	b := h.Bucket(buckets[0])
	if nil == b {
		/* Bucket not found */
		return nil, fmt.Errorf(
			"bucket %q not found",
			tgt,
		)
	}
	/* Get the next level */
	return Dive(b, buckets[1:], tgt, bsep)
}
