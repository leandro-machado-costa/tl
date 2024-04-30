package repository

import (
	"github.com/leandro-machado-costa/tl/internal/configenv/db"
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

	rows, err := conn.Query("SELECT id, title, descrition, featured_image, permisson_id, user_id, updated_at, created_at FROM courses WHERE permisson_id = 1")
	for rows.Next() {
		var course domain.Courses
		err = rows.Scan(&course.ID, &course.Title, &course.Descrition, &course.FeaturedImage, &course.PermissonID, &course.UserID, &course.Updated_at, &course.Created_at)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	defer rows.Close()

	return courses, err
}

func GetCourseByID(id int64) (course domain.Courses, err error) {
	conn := db.GetDB()

	row := conn.QueryRow("SELECT id, title, descrition,featured_image, permisson_id, user_id, updated_at, created_at FROM courses WHERE id = $1 AND permisson_id <> 2", id)

	err = row.Scan(&course.ID, &course.Title, &course.Descrition, &course.FeaturedImage, &course.PermissonID, &course.UserID, &course.Updated_at, &course.Created_at)

	modules, err := conn.Query("SELECT id, title, descrition, course_id FROM modules WHERE course_id = $1", course.ID)

	if err != nil {
		return domain.Courses{}, err // Retorna uma instância vazia da struct
	}
	defer modules.Close()

	for modules.Next() {
		module := domain.Module{}
		err := modules.Scan(&module.ID, &module.Title, &module.Descrition, &module.CourseID)
		if err != nil {
			return domain.Courses{}, err // Retorna uma instância vazia da struct
		}

		// Consultar as lições para cada módulo
		lessons, err := conn.Query(`SELECT id, "order", title, lesson, module_id FROM lessons WHERE module_id = $1`, module.ID)
		if err != nil {
			return domain.Courses{}, err // Retorna uma instância vazia da struct
		}
		defer lessons.Close()

		for lessons.Next() {
			lesson := domain.Lesson{}
			err := lessons.Scan(&lesson.ID, &lesson.Order, &lesson.Title, &lesson.Lesson, &lesson.ModuleID)
			if err != nil {
				return domain.Courses{}, err // Retorna uma instância vazia da struct
			}
			module.Lessons = append(module.Lessons, lesson)
		}

		course.Modules = append(course.Modules, module)
	}

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
