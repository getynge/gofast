package permissions

type BasePermission uint32
type Permission struct {
	Base BasePermission
	Args []string
}

func NewPermission(base BasePermission, args ...string) Permission {
	return Permission{
		Base: base,
		Args: args,
	}
}

const (
	AccessSystemPath BasePermission = iota
	DownloadScript
)
