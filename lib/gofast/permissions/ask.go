package permissions

import (
	"errors"
	"github.com/getynge/gofast/state"
	"math/big"
)

var ParseError = errors.New("failed to parse id")

type Id [32]byte

func ParseId(id string) (Id, error) {
	var out Id
	rid := big.NewInt(0)
	_, good := rid.SetString(id, 10)
	if !good {
		return out, ParseError
	}
	bid := rid.Bytes()

	for i, b := range bid[:32] {
		out[i] = b
	}

	return out, nil
}

// Request requests the given permission, returning whether it was granted
func Request(id Id, perm Permission) bool {
	cid := state.Id(id)
	if state.Permissions[cid] == nil {
		state.Permissions[cid] = make(map[state.Permission][]string)
	}
	state.Permissions[cid][state.Permission(perm.Base)] = perm.Args
	return true
}
