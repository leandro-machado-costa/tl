package repository

import (
	"github.com/leandro-machado-costa/tl/internal/config/db"
	"github.com/leandro-machado-costa/tl/internal/domain"
)

func DeleteUserByID(id int64) (int64, error) {
	conn := db.GetDB()

	res, err := conn.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func GetUsers() (users []domain.Users, err error) {
	conn := db.GetDB()

	rows, err := conn.Query("SELECT id, username, email, name, resume, picture, role_id,updated_at, created_at FROM users")
	for rows.Next() {
		var u domain.Users
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Name, &u.Resume, &u.Picture, &u.RoleID, &u.Updated_at, &u.Created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	defer rows.Close()

	return users, err
}

func GetUserByID(id int64) (user domain.Users, err error) {
	conn := db.GetDB()

	row := conn.QueryRow("SELECT id, username, email, name, resume, picture, role_id,updated_at, created_at FROM users WHERE id = $1", id)

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Resume, &user.Picture, &user.RoleID, &user.Updated_at, &user.Created_at)

	return user, err
}

func InsertUser(user domain.Users) (id int64, err error) {
	conn := db.GetDB()

	sql := `INSERT INTO users (username, email, name, resume, picture, role_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = conn.QueryRow(sql, user.Username, user.Email, user.Name, user.Resume, user.Picture, user.RoleID).Scan(&id)

	return id, err
}

func UpdateUserByID(id int64, user domain.Users) (int64, error) {
	conn := db.GetDB()

	res, err := conn.Exec(`UPDATE users SET username = $1, email = $2, name = $3, resume = $4, picture = $5, role_id = $6 WHERE id = $7`,
		user.Username, user.Email, user.Name, user.Resume, user.Picture, user.RoleID, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
