package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
	_ "github.com/gernest/laisense/statik"
	"github.com/olekukonko/tablewriter"
	"github.com/rakyll/statik/fs"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

//go:generate statik -f -m -src=./license-list-data/json/details/

// Default index directory
const defIdx = ".licence.search"

type D struct {
	Name string `json:"name"`
	Text string `json:"licenseText"`
	ID   string `json:"licenseId"`
}

func main() {
	x, err := zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		log.Fatal(err)
	}
	defer x.Sync()
	a := cli.NewApp()
	a.Name = "laisense"
	a.Usage = "Guard your go.mod with the right LICENCE dependencies"
	a.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "index,i",
			Usage:  "Directory in which the search index is/should be stored",
			EnvVar: "LAISENSE_INDEX_DIR",
		},
		cli.BoolFlag{
			Name:  "json,j",
			Usage: "Renders output as json",
		},
	}
	a.Action = run(x)
	if err := a.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(z *zap.Logger) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		idx, err := setup(ctx.GlobalString("index"))
		if err != nil {
			return err
		}
		defer idx.Close()
		wd := ctx.Args().First()
		if wd == "" {
			wd, err = os.Getwd()
			if err != nil {
				return err
			}
		}
		return do(z, idx, wd, ctx.GlobalBool("json"))
	}
}

func match(n string) bool {
	switch n {
	case "licence", "license":
		return true
	default:
		return false
	}
}

func createIndex(o string) {
	log.Printf("creating index file into :%s \n", o)
	x, err := fs.New()
	if err != nil {
		log.Fatalf("Failed to initialize embeded files: %v", err)
	}
	var m []*D
	err = fs.Walk(x, "/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		f, err := x.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		var n D
		err = json.NewDecoder(f).Decode(&n)
		if err != nil {
			return err
		}
		m = append(m, &n)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to load licence  details: %v", err)
	}
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(o, mapping)
	if err != nil {
		log.Fatalf("Failed to initialize  index file at %s: %v", o, err)
	}
	defer index.Close()
	for _, v := range m {
		err := index.Index(v.ID, v)
		if err != nil {
			log.Fatalf("Failed to index  %s:%q  index file: %v", o, v.ID, err)
		}
	}
	log.Printf("successful created index file into :%s \n", o)
}

func search(z *zap.Logger, index bleve.Index, data string) (string, error) {
	z.Debug("Opening beleve index")

	query := bleve.NewQueryStringQuery(data)
	searchRequest := bleve.NewSearchRequest(query)
	z.Debug("Searching the index")
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return "", err
	}
	if len(searchResult.Hits) > 0 {
		return searchResult.Hits[0].ID, nil
	}
	return "", nil
}

func escape(s string) string {
	var buf bytes.Buffer
	for _, v := range s {
		switch v {
		case '+', '-', '=', '&', '|', '>', ';', '<', '!', '(', ')',
			'{', '}', '[', ']', '^', '"', '~', '*', '?', ':', '\\', '/':
			buf.WriteByte('\\')
			buf.WriteRune(v)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}
func setup(indexDir string) (bleve.Index, error) {
	if indexDir == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Failed to get user home directory %e ", err)
		}
		indexDir = filepath.Join(h, defIdx)
	}
	_, err := os.Stat(indexDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("Failed to check stats for %s %e ", indexDir, err)
		}
		createIndex(indexDir)
	}
	return bleve.Open(indexDir)
}

func do(z *zap.Logger, idx bleve.Index, dir string, toJSON bool) error {
	v, err := build.ImportDir(dir, 0)
	if err != nil {
		return err
	}
	w := z.With(zap.String("PkgPath", v.Dir))
	w.Debug("Detected package")
	all, err := list(v.Dir)
	if err != nil {
		w.Error("Failed to list dependencies", zap.String("err", err.Error()))
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(all))
	for _, v := range all {
		go func(m *Module) {
			find(w, idx, m)
			wg.Done()
		}(v)
	}
	wg.Wait()
	if toJSON {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		return e.Encode(all)
	}
	summary(all)
	output(all)
	return nil
}

func summary(all []*Module) {
	var notLicensed int
	var notDetected int
	for _, v := range all {
		id := v.Licence.ID
		if !v.Licence.Found {
			id = "n/a"
			notLicensed++
		}
		if id == "" {
			notDetected++
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NOT LICENSED", "UNKNOWN", "TOTAL"})
	table.Append([]string{
		fmt.Sprint(notLicensed),
		fmt.Sprint(notDetected),
		fmt.Sprint(len(all)),
	})
	table.Render()
}

func output(all []*Module) {
	var data [][]string
	for _, v := range all {
		id := v.Licence.ID
		if !v.Licence.Found {
			id = "n/a"
		}
		if id == "" {
			id = v.Licence.Hint
		}
		data = append(data, []string{v.Path, id})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PACKAGE", "LICENSE"})
	table.AppendBulk(data)
	table.Render()
}

func find(w *zap.Logger, idx bleve.Index, v *Module) error {
	m := w.With(
		zap.String("Module", v.Dir),
		zap.String("Version", v.Version),
	)
	m.Debug("detected module")
	// step 0 check for any prospect license files
	m.Debug("checking files in module dir")
	o, err := ioutil.ReadDir(v.Dir)
	if err != nil {
		return err
	}
	for _, file := range o {
		if file.IsDir() {
			continue
		}
		k := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		k = strings.ToLower(k)
		full := filepath.Join(v.Dir, file.Name())
		if match(k) {
			v.Licence.Found = true
			m.Debug("Found licence", zap.String("LicencePath", full))
			l, err := ioutil.ReadFile(full)
			if err != nil {
				v.Licence.Hint = "Failed to read file"
				return err
			}
			ls := string(l)
			id, err := search(m, idx, ls)
			if err != nil {
				id, err = search(m, idx, escape(ls))
			}
			if err != nil {
				return err
			}
			if id != "" {
				m.Debug("Matched Licence", zap.String("LicenceID", id))
				v.Licence.ID = id
			} else {

				v.Licence.Hint = trunc(ls)
				m.Debug("No licence match found")
			}
			return nil
		}
	}
	m.Debug("No license found")
	return nil
}

func trunc(s string) string {
	if len(s) < 255 {
		return s
	}
	return s[:255]
}

type Module struct {
	Path     string // module path
	Version  string // module version
	Error    string // error loading module
	Info     string // absolute path to cached .info file
	GoMod    string // absolute path to cached .mod file
	Zip      string // absolute path to cached .zip file
	Dir      string // absolute path to cached source root directory
	Sum      string // checksum for path, version (as in go.sum)
	GoModSum string // checksum for go.mod (as in go.sum)
	Licence  Licence
}

type Licence struct {
	Found bool   `json:",omitempty"`
	ID    string `json:",omitempty"`
	Hint  string `json:",omitempty"`
}

func list(dir string) ([]*Module, error) {
	x := exec.Command("go", "list", "-json", "-m", "all")
	var buf bytes.Buffer
	x.Dir = dir
	x.Stdout = &buf
	x.Run()
	dec := json.NewDecoder(&buf)
	var o []*Module
	for dec.More() {
		var m Module
		err := dec.Decode(&m)
		if err != nil {
			return nil, err
		}
		o = append(o, &m)
	}
	return o, nil
}
