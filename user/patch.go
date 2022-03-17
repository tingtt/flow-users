package user

import (
	"flow-users/mysql"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type PatchBody struct {
	Name     *string `json:"name" form:"name" validate:"omitempty"`
	Email    *string `json:"email" form:"email" validate:"omitempty,email"`
	Password *string `json:"password" form:"password" validate:"omitempty"`
}

func Patch(id uint64, new PatchBody) (r UserWithOutPassword, invalidEmail bool, usedEmail bool, notFound bool, err error) {
	// Get old
	old, notFound, err := Get(id)
	if err != nil {
		return
	}
	if notFound {
		return
	}
	r.Id = id
	r.Name = old.Name
	r.Email = old.Email

	// Generate query
	queryStr := "UPDATE users SET "
	var queryParams []interface{}
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		r.Name = *new.Name
	}
	if new.Email != nil {
		queryStr += " email = ?,"
		queryParams = append(queryParams, new.Email)
		r.Email = *new.Email
	}
	if new.Password != nil {
		queryStr += " password = ?"
		// Create password hash
		var hashed []byte
		hashed, err = bcrypt.GenerateFromPassword([]byte(*new.Password), 10)
		if err != nil {
			return
		}
		queryParams = append(queryParams, hashed)
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE id = ?"
	queryParams = append(queryParams, id)

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(queryParams...)
	if err != nil {
		return
	}

	return UserWithOutPassword{id, *new.Name, *new.Email}, false, false, false, nil
}
