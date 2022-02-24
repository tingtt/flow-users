package user

import (
	"flow-users/mysql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       uint64
	Name     string
	Email    string
	Password []byte
}

type UserPost struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UserPatch struct {
	Name     *string `json:"name" form:"name" validate:"omitempty"`
	Email    *string `json:"email" form:"email" validate:"omitempty,email"`
	Password *string `json:"password" form:"password" validate:"omitempty"`
}

type UserPostResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserSignInPost struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (post *UserPost) PostResponse(id uint64) UserPostResponse {
	return UserPostResponse{
		Id:    id,
		Name:  post.Name,
		Email: post.Email,
	}
}

func Get(id uint64) (u User, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return User{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, email, password FROM users WHERE id = ?")
	if err != nil {
		return User{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(id)
	if err != nil {
		return User{}, false, err
	}

	var (
		name     string
		email    string
		password []byte
	)
	if !rows.Next() {
		// Not found
		return User{}, true, nil
	}
	err = rows.Scan(&name, &email, &password)
	if err != nil {
		return User{}, false, err
	}

	return User{id, name, email, password}, false, nil
}

func GetByEmail(email string) (u User, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return User{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, name, password FROM users WHERE email = ?")
	if err != nil {
		return User{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(email)
	if err != nil {
		return User{}, false, err
	}

	var (
		id       uint64
		name     string
		password []byte
	)
	if !rows.Next() {
		// Not found
		return User{}, true, nil
	}
	err = rows.Scan(&id, &name, &password)
	if err != nil {
		return User{}, false, err
	}

	return User{id, name, email, password}, false, nil
}

func Insert(post UserPost) (u User, invalidEmail bool, usedEmail bool, err error) {
	_, notFound, err := GetByEmail(post.Email)
	if err != nil {
		return User{}, false, false, err
	}
	if !notFound {
		return User{}, false, true, nil
	}

	// Create password hash
	hashed, err := bcrypt.GenerateFromPassword([]byte(post.Password), 10)
	if err != nil {
		return User{}, false, false, err
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return User{}, false, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO users (name, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return User{}, false, false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(post.Name, post.Email, hashed)
	if err != nil {
		return User{}, false, false, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, false, false, err
	}

	return User{uint64(id), post.Name, post.Email, hashed}, false, false, nil
}

func (u *User) Verify(password string) (varify bool, err error) {
	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func Update(id uint64, new UserPatch) (r UserPostResponse, invalidEmail bool, usedEmail bool, notFound bool, err error) {
	// Get old
	u, notFound, err := Get(id)
	if err != nil {
		return UserPostResponse{}, false, false, false, err
	}
	if notFound {
		return UserPostResponse{}, false, false, true, nil
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
			return UserPostResponse{}, false, false, false, err
		}
	}

	db, err := mysql.Open()
	if err != nil {
		return UserPostResponse{}, false, false, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?")
	if err != nil {
		return UserPostResponse{}, false, false, false, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(new.Name, new.Email, hashed, id)
	if err != nil {
		return UserPostResponse{}, false, false, false, err
	}

	return UserPostResponse{id, *new.Name, *new.Email}, false, false, false, nil
}

func Delete(id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(id)
	if err != nil {
		return false, err
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRowCount == 0 {
		// Not found
		return true, nil
	}

	return false, nil
}
