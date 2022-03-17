package user

import "flow-users/mysql"

func Delete(id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(id)
	if err != nil {
		return
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return
	}
	notFound = affectedRowCount == 0

	return false, nil
}
