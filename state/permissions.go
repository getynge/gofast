package state

type Permission uint32

var Permissions = make(map[Id]map[Permission][]string)

func Check(id Id, permission Permission, args []string) bool {
	perms := Permissions[id]
	if perms == nil {
		return false
	}
	pa := perms[permission]
	if pa == nil {
		return false
	}

	switch permission {
	case AccessSystemPath:
		fallthrough
	case DownloadScript:
		matches := 0
		for _, a := range args {
			for _, p := range pa {
				if a == p {
					matches += 1
				}
			}
		}
		if matches != len(args) {
			return false
		}
	}

	return true
}

const (
	AccessSystemPath Permission = iota
	DownloadScript
)
