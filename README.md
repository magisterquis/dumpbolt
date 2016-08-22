dumpbolt
========
Dumps the contents of a [Bolt database](https://github.com/boltdb/bolt).

Give it a database file and, optionally, a starting point, and it will print
what's in the database.

Installation
------------
```bash
go get github.com/magisterquis/dumpbolt
```

Quickstart
----------
Print the contents of `example.db`
```bash
dumpbolt example.db
```

Print the contents of `example.db`, starting in the bucket `foo`, which itself
is in the bucket `bar`
```bash
dumpbolt example.db /bar/foo
```

Print `example.db` with the full bucket path for each key/value pair
```bash
dumpbolt -a example.db
```

Print starting at the bucket named `bar/foo` inside the bucket `tridge`
```bash
dumpbolt -p '^' exampledb '^tridge^bar/foo'

Output
------
Output comes in two formats.  By default, a more human-friendly layout is used.
```
foo/
        bar/
                tridge -> Hello,
                baaz -> World!
        quux/
                curly -> 4
```
For ease of grepping and using other tools, the entire path can be put on each
line with `-a`
```
/foo/
/foo/bar/
/foo/bar/tridge -> Hello,
/foo/bar/baaz -> World!
/foo/quux/
/foo/quux/curly -> 4
```
