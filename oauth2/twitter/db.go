package twitter

import (
	"flow-user/mysql"
	"time"
)

type OAuth2 struct {
	AccessToken          string
	ExpireIn             int64
	RefreshToken         string
	RefreshTokenExpireIn int64
	OwnerId              string
}

func Get(user_id uint64) (o OAuth2, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return OAuth2{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT access_token, access_token_expire_in, refresh_token, refresh_token_expire_in, owner_id FROM twitter_oauth2_tokens WHERE user_id = ?")
	if err != nil {
		return OAuth2{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(user_id)
	if err != nil {
		return OAuth2{}, false, err
	}

	var (
		access_token            string
		access_token_expire_in  time.Time
		refresh_token           string
		refresh_token_expire_in time.Time
		owner_id                string
	)
	if !rows.Next() {
		// Not found
		return OAuth2{}, true, nil
	}
	err = rows.Scan(&access_token, &access_token_expire_in, &refresh_token, &refresh_token_expire_in, &owner_id)
	if err != nil {
		return OAuth2{}, false, err
	}

	return OAuth2{access_token, access_token_expire_in.Unix(), refresh_token, refresh_token_expire_in.Unix(), owner_id}, false, nil
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
	stmtIns, err := db.Prepare("INSERT INTO twitter_oauth2_tokens (user_id, access_token, access_token_expire_in, refresh_token, refresh_token_expire_in, owner_id) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return OAuth2{}, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(user_id, o.AccessToken, o.ExpireIn, o.RefreshToken, o.RefreshTokenExpireIn, o.OwnerId)
	if err != nil {
		return OAuth2{}, err
	}

	return OAuth2{o.AccessToken, o.ExpireIn, o.RefreshToken, o.RefreshTokenExpireIn, o.OwnerId}, nil
}

func Delete(user_id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM twitter_oauth2_tokens WHERE user_id = ?")
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
