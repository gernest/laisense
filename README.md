# laisense

Gives you insight on `go.mod` dependencies licence

# installation

```
go get github.com/gernest/laisense
```

# Usage

Just call it with no arg on the go project with `go.mod`

```
laisense
```

This might take a while running for the first time as, we need to generate search index. Subsequent calls will be faster as the index is created once.

By default the index will be stored in `$HOME/.licence.search` directory. If you want the index to be stored somewhere else (Note you will have to either set env var `$LAISENSE_INDEX_DIR` to point to this new path or specify it every time you invoke the binary with the `--index` flag) pass the `--index` flag.

# help

```
NAME:
   laisense - Guard your go.mod with the right LICENCE dependencies

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --index value, -i value  Directory in which the search index is/should be stored [$LAISENSE_INDEX_DIR]
   --json, -j               Renders output as json
   --help, -h               show help
```

# example  output

```
+--------------+---------+-------+
| NOT LICENSED | UNKNOWN | TOTAL |
+--------------+---------+-------+
|            3 |       0 |    94 |
+--------------+---------+-------+
+-------------------------------------------+--------------------+
|                  PACKAGE                  |      LICENSE       |
+-------------------------------------------+--------------------+
| github.com/gernest/laisense               | MIT                |
| github.com/BurntSushi/toml                | n/a                |
| github.com/RoaringBitmap/roaring          | Apache-2.0         |
| github.com/armon/consul-api               | ODbL-1.0           |
| github.com/blevesearch/bleve              | Apache-2.0         |
| github.com/blevesearch/blevex             | n/a                |
| github.com/blevesearch/go-porterstemmer   | MIT                |
| github.com/blevesearch/mmap-go            | RPL-1.5            |
| github.com/blevesearch/segment            | Apache-2.0         |
| github.com/blevesearch/snowballstem       | n/a                |
| github.com/blevesearch/zap/v11            | Apache-2.0         |
| github.com/blevesearch/zap/v12            | Apache-2.0         |
| github.com/blevesearch/zap/v13            | Apache-2.0         |
| github.com/blevesearch/zap/v14            | Apache-2.0         |
| github.com/blevesearch/zap/v15            | Apache-2.0         |
| github.com/coreos/etcd                    | Apache-2.0         |
| github.com/coreos/go-etcd                 | Apache-2.0         |
| github.com/coreos/go-semver               | Apache-2.0         |
| github.com/couchbase/ghistogram           | Apache-2.0         |
| github.com/couchbase/moss                 | Apache-2.0         |
| github.com/couchbase/vellum               | Apache-2.0         |
| github.com/cpuguy83/go-md2man             | MIT                |
| github.com/cpuguy83/go-md2man/v2          | MIT                |
| github.com/davecgh/go-spew                | ISC                |
| github.com/fsnotify/fsnotify              | RPL-1.5            |
| github.com/glycerine/go-unsnap-stream     | MIT                |
| github.com/glycerine/goconvey             | MIT                |
| github.com/golang/protobuf                | RPL-1.5            |
| github.com/golang/snappy                  | RPL-1.5            |
| github.com/google/renameio                | Apache-2.0         |
| github.com/gopherjs/gopherjs              | BSD-2-Clause       |
| github.com/hashicorp/hcl                  | ODC-By-1.0         |
| github.com/hpcloud/tail                   | MIT                |
| github.com/inconshreveable/mousetrap      | Apache-2.0         |
| github.com/jtolds/gls                     | MIT                |
| github.com/kisielk/gotool                 | MIT                |
| github.com/kljensen/snowball              | MIT                |
| github.com/kr/pretty                      | MIT                |
| github.com/kr/pty                         | MIT                |
| github.com/kr/text                        | MIT                |
| github.com/magiconair/properties          | BSD-2-Clause       |
| github.com/mattn/go-runewidth             | MIT                |
| github.com/mitchellh/go-homedir           | MIT                |
| github.com/mitchellh/mapstructure         | MIT                |
| github.com/mschoch/smat                   | Apache-2.0         |
| github.com/olekukonko/tablewriter         | MIT                |
| github.com/onsi/ginkgo                    | MIT                |
| github.com/onsi/gomega                    | MIT                |
| github.com/pelletier/go-toml              | MIT                |
| github.com/philhofer/fwd                  | MIT                |
| github.com/pkg/errors                     | BSD-2-Clause       |
| github.com/pmezard/go-difflib             | BSD-3-Clause       |
| github.com/rakyll/statik                  | Apache-2.0         |
| github.com/rcrowley/go-metrics            | BSD-2-Clause-Views |
| github.com/rogpeppe/go-internal           | RPL-1.5            |
| github.com/russross/blackfriday           | BSD-2-Clause       |
| github.com/russross/blackfriday/v2        | BSD-2-Clause       |
| github.com/shurcooL/sanitized_anchor_name | MIT                |
| github.com/spf13/afero                    | Apache-2.0         |
| github.com/spf13/cast                     | MIT                |
| github.com/spf13/cobra                    | Apache-2.0         |
| github.com/spf13/jwalterweatherman        | MIT                |
| github.com/spf13/pflag                    | RPL-1.5            |
| github.com/spf13/viper                    | MIT                |
| github.com/steveyen/gtreap                | MIT                |
| github.com/stretchr/objx                  | MIT                |
| github.com/stretchr/testify               | MIT                |
| github.com/syndtr/goleveldb               | RPL-1.5            |
| github.com/tinylib/msgp                   | MIT                |
| github.com/ugorji/go/codec                | MIT                |
| github.com/urfave/cli                     | MIT                |
| github.com/willf/bitset                   | RPL-1.5            |
| github.com/xordataexchange/crypt          | MIT                |
| github.com/yuin/goldmark                  | MIT                |
| go.etcd.io/bbolt                          | MIT                |
| go.uber.org/atomic                        | MIT                |
| go.uber.org/multierr                      | MIT                |
| go.uber.org/tools                         | MIT                |
| go.uber.org/zap                           | MIT                |
| golang.org/x/crypto                       | RPL-1.5            |
| golang.org/x/lint                         | RPL-1.5            |
| golang.org/x/mod                          | RPL-1.5            |
| golang.org/x/net                          | RPL-1.5            |
| golang.org/x/sync                         | RPL-1.5            |
| golang.org/x/sys                          | RPL-1.5            |
| golang.org/x/text                         | RPL-1.5            |
| golang.org/x/tools                        | RPL-1.5            |
| golang.org/x/xerrors                      | RPL-1.5            |
| gopkg.in/check.v1                         | BSD-2-Clause       |
| gopkg.in/errgo.v2                         | RPL-1.5            |
| gopkg.in/fsnotify.v1                      | RPL-1.5            |
| gopkg.in/tomb.v1                          | RPL-1.5            |
| gopkg.in/yaml.v2                          | Apache-2.0         |
| honnef.co/go/tools                        | MIT                |
+-------------------------------------------+--------------------+
```