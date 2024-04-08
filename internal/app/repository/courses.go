package repository

import (
	"github.com/leandro-machado-costa/tl/internal/config/db"
	"github.com/leandro-machado-costa/tl/internal/domain"
)

func DeleteCourseByID(id int64) (int64, error) {
	conn := db.GetDB()

	res, err := conn.Exec(`DELETE FROM courses WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func GetCourses() (courses []domain.Courses, err error) {
	conn := db.GetDB()

	rows, err := conn.Query("SELECT id, title, descrition, permisson_id, user_id, updated_at, created_at fROM courses")
	for rows.Next() {
		var u domain.Courses
		err = rows.Scan(&u.ID, &u.Title, &u.Descrition, &u.PermissonID, &u.UserID, &u.Updated_at, &u.Created_at)
		if err != nil {
			return nil, err
		}
		courses = append(courses, u)
	}
	defer rows.Close()

	return courses, err
}

func GetCourseByID(id int64) (course domain.Courses, err error) {
	conn := db.GetDB()

	row := conn.QueryRow("SELECT id, title, descrition, permisson_id, user_id, updated_at, created_at FROM courses WHERE id = $1", id)

	err = row.Scan(&course.ID, &course.Title, &course.Descrition, &course.PermissonID, &course.UserID, &course.Updated_at, &course.Created_at)

	return course, err
}

func InsertCourse(course domain.Courses) (id int64, err error) {
	conn := db.GetDB()

	sql := `INSERT INTO courses (title, descrition, permisson_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err = conn.QueryRow(sql, course.Title, course.Descrition, course.PermissonID, course.UserID).Scan(&id)

	return id, err
}

func UpdateCourseByID(id int64, course domain.Courses) (int64, error) {
	conn := db.GetDB()

	res, err := conn.Exec(`UPDATE courses SET title = $1, descrition = $2, permisson_id = $3, user_id = $4 WHERE id = $5`, course.Title, course.Descrition, course.PermissonID, course.UserID, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
