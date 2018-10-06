package models

// UserRight defines the rights users can have for lists/namespaces
type UserRight int

// define unknown user right
const (
	UserRightUnknown = -1
)

// Enumerate all the user rights
const (
	// Can read lists in a User
	UserRightRead UserRight = iota
	// Can write tasks in a User like lists and todo tasks. Cannot create new lists.
	UserRightWrite
	// Can manage a list/namespace, can do everything
	UserRightAdmin
)

func (r UserRight) isValid() error {
	if r != UserRightAdmin && r != UserRightRead && r != UserRightWrite {
		return ErrInvalidUserRight{r}
	}

	return nil
}

// CanCreate checks if the user can create a new user <-> list relation
func (lu *ListUser) CanCreate(doer *User) bool {
	// Get the list and check if the user has write access on it
	l := List{ID: lu.ListID}
	if err := l.GetSimpleByID(); err != nil {
		Log.Error("Error occurred during CanCreate for ListUser: %s", err)
		return false
	}
	return l.CanWrite(doer)
}

// CanDelete checks if the user can delete a user <-> list relation
func (lu *ListUser) CanDelete(doer *User) bool {
	// Get the list and check if the user has write access on it
	l := List{ID: lu.ListID}
	if err := l.GetSimpleByID(); err != nil {
		Log.Error("Error occurred during CanDelete for ListUser: %s", err)
		return false
	}
	return l.CanWrite(doer)
}

// CanUpdate checks if the user can update a user <-> list relation
func (lu *ListUser) CanUpdate(doer *User) bool {
	// Get the list and check if the user has write access on it
	l := List{ID: lu.ListID}
	if err := l.GetSimpleByID(); err != nil {
		Log.Error("Error occurred during CanUpdate for ListUser: %s", err)
		return false
	}
	return l.CanWrite(doer)
}
