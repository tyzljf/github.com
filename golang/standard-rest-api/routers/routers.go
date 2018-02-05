package routers

import (
	"net/http"
	"github.com/golang/standard-rest-api/controllers"
)

func CreateRouters(mux *http.ServeMux, uc *controllers.UserController, jc *controllers.JobController) {
	mux.HandleFunc("/register", uc.Register)
	mux.HandleFunc("/login", uc.Login)

	mux.HandleFunc("/job", jc.Create)
	mux.HandleFunc("/job/", jc.Job)
	mux.HandleFunc("/feed", jc.Feed)
}
