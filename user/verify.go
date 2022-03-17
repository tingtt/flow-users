package user

import "golang.org/x/crypto/bcrypt"

type VerifyPostBody struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (u *User) Verify(password string) (varify bool, err error) {
	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
