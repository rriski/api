package models

// Create creates a new team <-> namespace relation
func (tn *TeamNamespace) Create(doer *User) (err error) {

	// Check if the rights are valid
	if err = tn.Right.isValid(); err != nil {
		return
	}

	// Check if the team exists
	_, err = GetTeamByID(tn.TeamID)
	if err != nil {
		return
	}

	// Check if the namespace exists
	_, err = GetNamespaceByID(tn.NamespaceID)
	if err != nil {
		return
	}

	// Check if the team already has access to the namespace
	exists, err := x.Where("team_id = ?", tn.TeamID).
		And("namespace_id = ?", tn.NamespaceID).
		Get(&TeamNamespace{})
	if err != nil {
		return
	}
	if exists {
		return ErrTeamAlreadyHasAccess{tn.TeamID, tn.NamespaceID}
	}

	// Insert the new team
	_, err = x.Insert(tn)
	return
}