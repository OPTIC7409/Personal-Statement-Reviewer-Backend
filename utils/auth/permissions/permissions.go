package permissions

const (
	User  = 1 << iota
	Admin = 2 << iota
	Dev   = 3 << iota
)

func HasPermission(userPermissions, permission int) bool {
	return userPermissions&permission == permission
}

func AddPermission(userPermissions, permission int) int {
	return userPermissions | permission
}

func RemovePermission(userPermissions, permission int) int {
	return userPermissions &^ permission
}
