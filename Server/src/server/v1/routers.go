package v1

import (
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	json.Marshal("ds")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("User Page")
	w.Write(c)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("All Users Page")
	w.Write(c)
}

func partnerHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("Partner Page")
	w.Write(c)
}

func partnersHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("All Partners Page")
	w.Write(c)
}

func lessonHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("Lesson Page")
	w.Write(c)
}

func lessonsHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("All Lessons Page")
	w.Write(c)
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("Topic Page")
	w.Write(c)
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := json.Marshal("All Topics Page")
	w.Write(c)
}

func InstallHandlers(router *mux.Router) {
	router.HandleFunc("/apis", apiHandler)

	router.HandleFunc("/user", userHandler)
	router.HandleFunc("/users", usersHandler)

	router.HandleFunc("/partner", partnerHandler)
	router.HandleFunc("/partners", partnersHandler)

	router.HandleFunc("/lesson", lessonHandler)
	router.HandleFunc("/lessons", lessonsHandler)

	router.HandleFunc("/topic", topicHandler)
	router.HandleFunc("/topics", topicsHandler)

	// web socket
}

