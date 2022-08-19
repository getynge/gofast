// arguments: extract <destination> <path> <file path>:<import path>
package main

import (
	"fmt"
	"github.com/traefik/yaegi/extract"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		panic("please provide at least a package destination, package path, and target")
		return
	}

	extractor := extract.Extractor{
		Dest:    os.Args[1],
		Include: []string{".*"},
	}

	targets := os.Args[3:]

	for _, target := range targets {
		var actualTargets []string
		parts := strings.Split(target, ":")
		if len(parts) != 2 {
			fmt.Printf("target %s must be of the format filepath:importpath, skipping", target)
			continue
		}
		actualTargets, err := filepath.Glob(parts[0])
		if err != nil {
			panic("glob failed: " + err.Error())
		}
		for _, at := range actualTargets {
			if !strings.HasPrefix(at, "./") {
				at = fmt.Sprintf("./%s", at)
			}
			base := filepath.Base(at)
			out := filepath.Join(os.Args[2], fmt.Sprintf("%s.go", base))
			if err := os.MkdirAll(filepath.Dir(out), 777); err != nil {
				panic("mkdir failed: " + err.Error())
			}
			f, err := os.Create(out)
			if err != nil {
				panic("file create failed: " + err.Error())
			}
			_, err = extractor.Extract(at, path.Join(parts[1], base), f)
			if err != nil {
				panic("extraction failed: " + err.Error())
			}
			fmt.Printf("created extractions for target %s at %s\n", at, out)
		}
	}
	fmt.Println("done")
}
