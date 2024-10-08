package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

// Struct for Student Details
type Student struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Faculty string `json:"faculty"`
	Gender  string `json:"gender"`
}

var (

	// List of students in the system
	studentList = []Student{}
	// Next available student ID
	idCounter = 1
	// Synchronize access to studentList
	studentLock sync.Mutex // Mutex for synchronizing access to studentList
)

func main() {
	// For creating and retrieving students
	http.HandleFunc("/students", handleAllStudents)
	// For retrieving, updating, and deleting a specific student
	http.HandleFunc("/students/", handleSingleStudent)

	// Start server on port 8090
	http.ListenAndServe(":8090", nil)
}

// Handles student collection-level operations (GET, POST)
func handleAllStudents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listStudents(w)
	case http.MethodPost:
		addNewStudent(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handles student-specific operations (GET, PUT, DELETE) by student ID
func handleSingleStudent(w http.ResponseWriter, r *http.Request) {
	// Extract student ID from URL
	idStr := r.URL.Path[len("/students/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		findStudentByID(w, id)
	case http.MethodPut:
		editStudent(w, r, id)
	case http.MethodDelete:
		removeStudent(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Returns the list of all students in JSON format
func listStudents(w http.ResponseWriter) {
	studentLock.Lock()
	defer studentLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentList)
}

// Adds a new student from the POST request body
func addNewStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	studentLock.Lock()
	student.ID = idCounter
	idCounter++
	studentList = append(studentList, student)
	studentLock.Unlock()

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// Finds and returns a student by their ID
func findStudentByID(w http.ResponseWriter, id int) {
	studentLock.Lock()
	defer studentLock.Unlock()

	for _, student := range studentList {
		if student.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// Updates an existing student by their ID
func editStudent(w http.ResponseWriter, r *http.Request, id int) {
	var updated Student
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	studentLock.Lock()
	defer studentLock.Unlock()

	for i, student := range studentList {
		if student.ID == id {
			studentList[i].Name = updated.Name
			studentList[i].Faculty = updated.Faculty
			studentList[i].Gender = updated.Gender
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(studentList[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// Deletes a student by user ID
func removeStudent(w http.ResponseWriter, id int) {
	studentLock.Lock()
	defer studentLock.Unlock()

	for i, student := range studentList {
		if student.ID == id {
			studentList = append(studentList[:i], studentList[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
