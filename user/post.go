package user

import (
	"flow-users/mysql"

	"golang.org/x/crypto/bcrypt"
)

type PostBody struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UserPostResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (post *PostBody) PostResponse(id uint64) UserPostResponse {
	return UserPostResponse{
		Id:    id,
		Name:  post.Name,
		Email: post.Email,
	}
}

func Post(post PostBody) (u User, invalidEmail bool, usedEmail bool, err error) {
	_, notFound, err := GetByEmail(post.Email)
	if err != nil {
		return
	}
	if !notFound {
		invalidEmail = true
		return
	}

	// Create password hash
	hashed, err := bcrypt.GenerateFromPassword([]byte(post.Password), 10)
	if err != nil {
		return
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO users (name, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(post.Name, post.Email, hashed)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	u.Id = uint64(id)
	u.Name = post.Name
	u.Email = post.Email
	u.Password = hashed
	return
}
