package controller

import (
	common_log "common/log"
	"encoding/json"
	"go-hj-hospital/HospitalQueue/service"
	"go-hj-hospital/config"
	"go-hj-hospital/model/constant"
	"go-hj-hospital/model/dbmodel"
	"go-hj-hospital/model/request"
	"go-hj-hospital/model/response"
	"go-hj-hospital/util"

	"github.com/kataras/iris/v12"
)

func Index(ctx iris.Context) {
	hospital_info := &response.HospitalInformation{
		HospitalName: config.Instance.HospitalName,
		Address:      config.Instance.HospitalAddress,
		Phone:        config.Instance.Phone,
	}

	queue_total, err := service.HospitalService.GetQueueTotal(util.Master())
	if err != nil {
		common_log.Error(err)
		return
	}

	patient_total, err := service.HospitalService.GetPatientTotal(util.Master())
	if err != nil {
		common_log.Error(err)
		return
	}

	doctor_total, err := service.HospitalService.GetDoctorTotal(util.Master())
	if err != nil {
		common_log.Error(err)
		return
	}

	information := &response.Information{
		HospitalInfo:      *hospital_info,
		QueueTotal:        queue_total,
		PatientRegistered: patient_total,
		Doctors:           doctor_total,
	}

	jsonData, err := json.Marshal(information)
	if err != nil {
		common_log.Error(err)
		return
	}
	ctx.ContentType("application/json")
	ctx.Write(jsonData)
}

func GetCurrentQueue(ctx iris.Context) {
	current_patient, err := service.HospitalService.GetCurrentQueueService(util.Master())
	if err != nil {
		common_log.Error(err)
		return
	}
	ctx.JSON(iris.Map{
		"Current patient": current_patient,
	})
}

func GetQueueList(ctx iris.Context) {
	list, total, err := service.HospitalService.GetQueueListService(util.Master())
	if err != nil {
		common_log.Error(err)
		return
	}
	ctx.JSON(iris.Map{
		"Total": total,
		"Queue": list,
	})
}

func GetDoctorsList(ctx iris.Context) {
	list, total, err := service.HospitalService.GetDoctorsListService(util.Master())
	if err != nil {
		common_log.Error(err)
	}
	ctx.JSON(iris.Map{
		"Total":   total,
		"Doctors": list,
	})
}

func GetPatientsList(ctx iris.Context) {
	list, total, err := service.HospitalService.GetPatientsListService(util.Master())
	if err != nil {
		common_log.Error(err)
	}
	ctx.JSON(iris.Map{
		"Total":    total,
		"Patients": list,
	})
}

func InputPatient(ctx iris.Context) {
	var params *request.Patients
	if err := ctx.ReadJSON(&params); err != nil {
		common_log.Error(err)
	}

	if err := service.HospitalService.InputPatientService(params, util.Master(), ctx); err != nil {
		common_log.Error(err)
		return
	}
	response.PatientAddMessageResult("Patient Successfully Added", params.Name, ctx)
}

func InputDoctor(ctx iris.Context) {
	var params *request.Doctors
	if err := ctx.ReadJSON(&params); err != nil {
		common_log.Error(err)
	}

	if err := service.HospitalService.InputDoctorService(params, util.Master(), ctx); err != nil {
		return
	}
	response.DoctorAddMessageResult("Doctor Successfully Added", params.Name, ctx)
}

func InputQueue(ctx iris.Context) {
	var params *dbmodel.Queue
	if err := ctx.ReadJSON(&params); err != nil {
		common_log.Error(err)
	}

	patient_name, total, err := service.HospitalService.GetPatientNameByNIK(params, util.Master())
	if total == 0 {
		response.QueueMessageResult("Pasien belum mendaftar", ctx)
		return
	}
	if err != nil {
		return
	}
	doctor_name, total, err := service.HospitalService.GetDoctorNameById(params, util.Master())
	if total == 0 {
		response.QueueMessageResult("Dokter tidak ada", ctx)
		return
	}

	if err := service.HospitalService.InputQueueService(params, util.Master(), ctx); err != nil {
		return
	}
	response.QueueAddMessageResult("Successfully Added", patient_name, doctor_name, ctx)
}

func UpdateQueueStatus(ctx iris.Context) {
	var params *request.Queue
	if err := ctx.ReadJSON(&params); err != nil {
		common_log.Error(err)
	}

	if patient_status, total, err := service.HospitalService.PatientStatus(params, util.Master()); patient_status == constant.PatientIsDone || total == 0 || err != nil {
		response.QueueMessageResult("Nomor antrian ini tidak ada", ctx)
		return
	}

	if patients, total, err := service.HospitalService.ValidateUpdate(params, util.Master()); err != nil || total > 0 {
		ctx.JSON(iris.Map{
			"Message":  "Pasien sebelumnya belum dilayani",
			"Patients": patients,
		})
	} else {
		patient_name, err := service.HospitalService.UpdateQueueStatusService(params, util.Master())
		if err != nil {
			return
		}
		ctx.JSON(iris.Map{
			"Message":      "Status pasien berhasil diubah",
			"Patient name": patient_name,
		})
	}
}
