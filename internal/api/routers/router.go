package routers

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	eRouter := ExecsRouter()
	tRouter := TeachersRouter()
	sRouter := StudentsRouter()

	sRouter.Handle("/", eRouter)
	tRouter.Handle("/", sRouter)
	return tRouter
}
