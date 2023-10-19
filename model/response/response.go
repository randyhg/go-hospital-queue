package response

import "github.com/kataras/iris/v12"

type Information struct {
	HospitalInfo      HospitalInformation `json:"hospital_info"`
	QueueTotal        int64               `json:"queue"`
	PatientRegistered int64               `json:"patient_registered"`
	Doctors           int64               `json:"doctors"`
}

type HospitalInformation struct {
	HospitalName string `json:"hospital_name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
}

type CurrentQueue struct {
	QueueNo     string `json:"queue_no"`
	PatientName string `json:"name"`
	DoctorName  string `json:"doctor_name"`
}

type Queue struct {
	QueueNo     string `json:"queue_no"`
	PatientName string `json:"patient_name"`
	DoctorName  string `json:"doctor_name"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}

type Patients struct {
	Name             string `json:"name"`
	Birthdate        string `json:"birthdate"`
	Sex              string `json:"sex"`
	Phone            string `json:"phone"`
	EmergencyContact string `json:"emergency_contact"`
}

type UnhandledPatients struct {
	QueueNo uint   `json:"queue_no"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type Doctors struct {
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	WorkDay        string `json:"work_day"`
}

func PatientAddMessageResult(message interface{}, name string, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"Status":        message,
		"Patient Added": name,
	})
}

func DoctorAddMessageResult(message interface{}, name string, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"Status":       message,
		"Doctor Added": name,
	})
}

func QueueAddMessageResult(message interface{}, name string, doctor_name string, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"Status":       message,
		"Patient name": name,
		"Doctor name":  doctor_name,
	})
}

func QueueMessageResult(message interface{}, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"Status": message,
	})
}

func DuplicateMessageResult(message interface{}, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"Status": message,
	})
}
