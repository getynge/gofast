//go:generate go run ./tools/extract extracted ./extracted ./lib/gofast/*:github.com/getynge/gofast/lib/gofast
package main

import (
	"crypto/rand"
	"github.com/getynge/gofast/extracted"
	"github.com/getynge/gofast/symbolproc"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"math/big"
	"os"
	"reflect"
)

var program = `package main

import (
	"gofast/permissions"
	"fmt"
	"os"
)

func main() {
	id := permissions.ParseId(os.Args[0])
	perm := permissions.NewPermission(permissions.AccessSystemPath, "/")
	if permissions.Request(id, perm) {
		fmt.Println("file permission granted")
	}
}
`

func main() {
	var opts interp.Options
	max := big.NewInt(0)
	max.Exp(big.NewInt(2), big.NewInt(255), nil)

	extractedSymbols := extracted.Symbols
	stdSymbols := stdlib.Symbols
	symbolproc.Clean(extractedSymbols, "github.com/getynge/gofast/lib/")
	symbolproc.DenyList(stdSymbols,
		"runtime",
		"expvar",
		"path//filepath",
		"syscall",
		"unsafe",
		"os",
		"os//.*",
		"net",
		"net//.*",
		"database",
		"database//.*",
		"debug",
		"debug//.*",
		"embed",
	)
	stdSymbols["os/os"] = map[string]reflect.Value{
		"Args": reflect.ValueOf(&os.Args).Elem(),
	}

	var id [32]byte
	rid, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("failed to set interpreter id " + err.Error())
	}
	sid := rid.String()
	bid := rid.Bytes()
	for i, b := range bid[:32] {
		id[i] = b
	}
	opts.Args = []string{sid}

	interpreter := interp.New(opts)
	if err := interpreter.Use(extractedSymbols); err != nil {
		panic("failed to add exports " + err.Error())
	}
	if err := interpreter.Use(stdSymbols); err != nil {
		panic("failed to add stdlibs " + err.Error())
	}

	prog, err := interpreter.Compile(program)
	if err != nil {
		panic(err)
	}

	interpreter.Execute(prog)
}
