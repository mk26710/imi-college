package permissions

const (
	// allows to view any user
	PermissionViewUser int64 = 1 << 0
	// allows to edit any user
	PermissionEditUser int64 = 1 << 1
	// allows to delete users
	PermissionDeleteUser int64 = 1 << 2

	// full access, absolute power
	PermissionAdmin int64 = 1 << 3
)

func HasPermissions(target int64, required int64) bool {
	// admins have full access
	if (target & PermissionAdmin) == PermissionAdmin {
		return true
	}
	return (target & required) == required
}

func HasViewUser(target int64) bool {
	return HasPermissions(target, PermissionViewUser)
}

func HasEditUser(target int64) bool {
	return HasPermissions(target, PermissionEditUser)
}

func HasDeleteUser(target int64) bool {
	return HasPermissions(target, PermissionDeleteUser)
}

func HasAdmin(target int64) bool {
	return HasPermissions(target, PermissionAdmin)
}

type PermissionTable struct {
	ViewUser   bool `json:"viewUser"`
	EditUser   bool `json:"editUser"`
	DeleteUser bool `json:"deleteUser"`
	Admin      bool `json:"admin"`
}

func NewPermissionTable(permissions int64) PermissionTable {
	return PermissionTable{
		ViewUser:   HasPermissions(permissions, PermissionViewUser),
		EditUser:   HasPermissions(permissions, PermissionEditUser),
		DeleteUser: HasPermissions(permissions, PermissionDeleteUser),
		Admin:      HasPermissions(permissions, PermissionAdmin),
	}
}
