package user

import "flow-users/mysql"

func Get(id uint64) (u User, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, email, password FROM users WHERE id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}

	err = rows.Scan(&u.Name, &u.Email, &u.Password)
	if err != nil {
		return
	}

	u.Id = id
	return
}

func GetByEmail(email string) (u User, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, name, password FROM users WHERE email = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(email)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&u.Id, &u.Name, &u.Password)
	if err != nil {
		return
	}

	u.Email = email
	return
}

func GetWithoutPassword(id uint64) (u UserWithoutPassword, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, email FROM users WHERE id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}

	err = rows.Scan(&u.Name, &u.Email)
	if err != nil {
		return
	}

	u.Id = id
	return
}
