package user

import (
	"flow-users/mysql"

	"golang.org/x/crypto/bcrypt"
)

type PatchBody struct {
	Name     *string `json:"name" form:"name" validate:"omitempty"`
	Email    *string `json:"email" form:"email" validate:"omitempty,email"`
	Password *string `json:"password" form:"password" validate:"omitempty"`
}

func Patch(id uint64, new PatchBody) (r UserPostResponse, invalidEmail bool, usedEmail bool, notFound bool, err error) {
	// Get old
	u, notFound, err := Get(id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Update values
	if new.Name == nil {
		new.Name = &u.Name
	}
	if new.Email == nil {
		new.Email = &u.Email
	}
	hashed := u.Password
	if new.Password != nil {
		// Create password hash
		hashed, err = bcrypt.GenerateFromPassword([]byte(*new.Password), 10)
		if err != nil {
			return
		}
	}

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(new.Name, new.Email, hashed, id)
	if err != nil {
		return
	}

	return UserPostResponse{id, *new.Name, *new.Email}, false, false, false, nil
}
