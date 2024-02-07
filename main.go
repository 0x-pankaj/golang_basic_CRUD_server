package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//model for course

type Course struct {
	CouseId     string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fake db
var courses []Course

// middleware , helper
func (c *Course) IsEmpty() bool {
	// return c.CouseId == "" && c.CourseName == ""
	return c.CourseName == ""
}

func main() {
	fmt.Println("api Building")

	r := mux.NewRouter()

	//seeding
	courses = append(courses, Course{CouseId: "2", CourseName: "Golang", CoursePrice: 400, Author: &Author{Fullname: "hitesh", Website: "go.dev"}})
	courses = append(courses, Course{CouseId: "4", CourseName: "Nodejs", CoursePrice: 300, Author: &Author{Fullname: "hitesh", Website: "lco.dev"}})

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourse).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOnecourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteOne).Methods("DELETE")

	//listen to a port
	log.Fatal(http.ListenAndServe(":4000", r))
}

// controllers
// serve home route

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1> Welcome to API by LearnCodeOnline </h1>"))
}

// get all course

func getAllCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all course ")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// get one course

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get one course")
	w.Header().Set("Content-Type", "application/json")

	//grab id
	params := mux.Vars(r)

	// loop through couse and find matching id and return the respone

	for _, course := range courses {
		if course.CouseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("NO Couse found  with given id " + params["id"])
	return

}

func createOnecourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create one course")
	w.Header().Set("Content-Type", "application/json")

	// what if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please Send Some Data")
	}

	// what about - {}
	var course Course
	json.NewDecoder(r.Body).Decode(&course)

	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	//generate unique id, string
	// append courses into courses
	course.CouseId = strconv.Itoa(rand.Intn(100))

	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return

}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updating one course")
	w.Header().Set("Content-Type", "application/json")

	// what if body is empty
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	// what about = {}
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if course.IsEmpty() {
		http.Error(w, "Course data is empty", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)

	// fmt.Println("params: ", params)

	// find and update
	for ind, c := range courses {
		if c.CouseId == params["id"] {
			courses = append(courses[:ind], courses[ind+1:]...)
			json.NewDecoder(r.Body).Decode(&course)
			course.CouseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	http.Error(w, "Course not found", http.StatusNotFound)
}

func deleteOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleting course one ")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	// loop id, remove and delete
	for ind, c := range courses {
		if c.CouseId == params["id"] {
			// courses = append(courses[:ind], courses[ind+1:]...)
			courses = append(courses[:ind], courses[ind+1:]...)
			json.NewEncoder(w).Encode("Successfully Deleted course ")
			return
		}
	}

	http.Error(w, "CourseID not found", http.StatusNotFound)
}
