package google

import (
	"flow-users/mysql"
)

type OAuth2 struct {
	AccessToken string
	OwnerId     string
}

func Get(user_id uint64) (o OAuth2, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return OAuth2{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT access_token, owner_id FROM google_oauth2_tokens WHERE user_id = ?")
	if err != nil {
		return OAuth2{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(user_id)
	if err != nil {
		return OAuth2{}, false, err
	}

	var (
		access_token string
		owner_id     string
	)
	if !rows.Next() {
		// Not found
		return OAuth2{}, true, nil
	}
	err = rows.Scan(&access_token, &owner_id)
	if err != nil {
		return OAuth2{}, false, err
	}

	return OAuth2{access_token, owner_id}, false, nil
}

func Insert(o OAuth2, user_id uint64) (OAuth2, error) {
	_, notFound, err := Get(user_id)
	if err != nil {
		return OAuth2{}, err
	}
	if !notFound {
		// Delete old
		_, err := Delete(user_id)
		if err != nil {
			return OAuth2{}, err
		}
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return OAuth2{}, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO google_oauth2_tokens (user_id, access_token, owner_id) VALUES(?, ?, ?)")
	if err != nil {
		return OAuth2{}, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(user_id, o.AccessToken, o.OwnerId)
	if err != nil {
		return OAuth2{}, err
	}

	return OAuth2{o.AccessToken, o.OwnerId}, nil
}

func Delete(user_id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM google_oauth2_tokens WHERE user_id = ?")
	if err != nil {
		return false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(user_id)
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
