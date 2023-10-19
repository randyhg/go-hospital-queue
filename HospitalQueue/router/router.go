package router

import (
	"go-hj-hospital/HospitalQueue/controller"

	"github.com/kataras/iris/v12"
)

func Routes(app *iris.Application) {
	app.Get("/", controller.Index)
	app.Get("/getPatient", controller.GetPatientsList)
	app.Get("/getDoctor", controller.GetDoctorsList)
	app.Get("/getQueue", controller.GetCurrentQueue)
	app.Get("/getQueueList", controller.GetQueueList)
	app.Post("/addPatient", controller.InputPatient)
	app.Post("/addDoctor", controller.InputDoctor)
	app.Post("/addQueue", controller.InputQueue)
	app.Post("/updateQueue", controller.UpdateQueueStatus)
}
