package userModel

import (
	"database/sql"
	"fmt"

	"github.com/Blackmamoth/fileforte/db"
	"github.com/Blackmamoth/fileforte/types"
)

func GetUserByEmail(email string) (*types.User, error) {
	rows, err := db.DB.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.Id == 0 {
		return nil, fmt.Errorf("user with email [%s] does not exist", email)
	}

	return u, nil

}

func CreateUser(user types.User) (sql.Result, error) {
	result, err := db.DB.Exec("INSERT INTO users(username, email, password) VALUES(?, ?, ?)", user.UserName, user.Email, user.Password)
	return result, err
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.Id,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
