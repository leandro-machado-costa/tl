package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/leandro-machado-costa/tl/internal/app/repository"

	"github.com/leandro-machado-costa/tl/internal/domain"
)

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course domain.Courses

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		log.Printf("Error decoding course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := repository.InsertCourse(course)

	var resp map[string]any

	if err != nil {
		resp = map[string]any{
			"Error":   true,
			"Message": "Error creating course",
		}
	} else {
		resp = map[string]any{
			"Error":   false,
			"Message": "Course created successfully",
			"ID":      id,
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteCourseByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error converting id to int: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rows, err := repository.DeleteCourseByID(int64(id))
	if err != nil {
		log.Printf("Error updating course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("Error Deleting course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"Error":   false,
		"Message": "Course deleting successfully",
		"ID":      id,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	return
}

func GetCourseByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	course, err := repository.GetCourseByID(int64(id))
	if err != nil {
		log.Printf("Error retornar course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := repository.GetCourses()
	if err != nil {

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)

}

func UpdateCoursesByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error converting id to int: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var Course domain.Courses

	err = json.NewDecoder(r.Body).Decode(&Course)
	if err != nil {
		log.Printf("Error updating course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rows, err := repository.UpdateCourseByID(int64(id), Course)
	if err != nil {
		log.Printf("Error updating course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("Error updating course: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"Error":   false,
		"Message": "Course updated successfully",
		"ID":      id,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	return
}
