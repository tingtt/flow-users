package user

import "flow-users/mysql"

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
