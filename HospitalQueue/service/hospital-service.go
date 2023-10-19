package service

import (
	common_log "common/log"
	"errors"
	"go-hj-hospital/model/constant"
	"go-hj-hospital/model/dbmodel"
	"go-hj-hospital/model/request"
	"go-hj-hospital/model/response"
	"go-hj-hospital/util"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

var HospitalService = newHospitalService()

type hospitalService struct {
}

func newHospitalService() *hospitalService {
	return &hospitalService{}
}

func (s *hospitalService) InputPatientService(patients *request.Patients, db *gorm.DB, ctx iris.Context) (err error) {
	// validasi apakah pasien sudah ada
	var count int64
	if err = db.Model(dbmodel.Patients{}).Where("nik = ?", patients.NIK).Count(&count).Error; err != nil {
		common_log.Error(err)
		return err
	}

	if count > 0 {
		response.DuplicateMessageResult("Pasien ini sudah ada", ctx)
		return errors.New("Pasien sudah ada")
	}

	if err = db.Model(dbmodel.Patients{}).Create(&patients).Error; err != nil {
		common_log.Error(err)
		return err
	}
	return nil
}

func (s *hospitalService) GetPatientsListService(db *gorm.DB) (list []response.Patients, total int64, err error) {
	if err = db.Model(dbmodel.Patients{}).Limit(50).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	if err = db.Model(dbmodel.Patients{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *hospitalService) InputDoctorService(doctors *request.Doctors, db *gorm.DB, ctx iris.Context) (err error) {
	var count int64
	if err = db.Model(dbmodel.Doctors{}).Where("employee_id = ?", doctors.EmployeeId).Count(&count).Error; err != nil {
		common_log.Error(err)
		return err
	}

	if count > 0 {
		response.DuplicateMessageResult("Dokter ini sudah ada", ctx)
		return errors.New("Dokter ini sudah ada")
	}

	if err = db.Create(&doctors).Error; err != nil {
		common_log.Error(err)
		return err
	}
	return nil
}

func (s *hospitalService) GetDoctorsListService(db *gorm.DB) (list []response.Doctors, total int64, err error) {
	if err = db.Model(dbmodel.Doctors{}).Limit(50).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	if err = db.Model(response.Doctors{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *hospitalService) InputQueueService(queue *dbmodel.Queue, db *gorm.DB, ctx iris.Context) (err error) {
	var count int64
	if err = util.Master().Model(dbmodel.Queue{}).Where("patient_nik = ?", queue.PatientNIK).Count(&count).Error; err != nil {
		common_log.Error(err)
		return err
	}
	if count >= 5 {
		response.QueueMessageResult("Anda sudah mencapai limit antrian per hari", ctx)
		return errors.New("Anda sudah mencapai limit antrian per hari") // buat json response ("Anda sudah mencapai limit antrian per hari")
	}

	if err = util.Master().Create(&queue).Error; err != nil {
		common_log.Error(err)
		return err
	}
	return nil
}

func (s *hospitalService) GetPatientNameByNIK(queue *dbmodel.Queue, db *gorm.DB) (patient_name string, total int64, err error) {
	if err = db.Model(dbmodel.Patients{}).Where("nik = ?", queue.PatientNIK).Count(&total).Error; err != nil {
		common_log.Error(err)
		return "Kueri gagal", total, err
	}

	query := "SELECT name FROM t_patients WHERE nik = " + queue.PatientNIK + " ORDER BY id LIMIT 1"
	if err = db.Raw(query).Scan(&patient_name).Error; err != nil {
		return "Kueri gagal", total, err
	}

	return patient_name, total, nil
}

func (s *hospitalService) GetDoctorNameById(queue *dbmodel.Queue, db *gorm.DB) (doctor_name string, total int64, err error) {
	if err = db.Model(dbmodel.Doctors{}).Where("employee_id = ?", queue.DoctorEmployeeId).Count(&total).Error; err != nil {
		common_log.Error(err)
		return "Kueri gagal", total, err
	}

	query := "SELECT name FROM t_doctors WHERE employee_id = " + queue.DoctorEmployeeId + " ORDER BY id LIMIT 1"
	if err = db.Raw(query).Scan(&doctor_name).Error; err != nil {
		return "Kueri gagal", total, err
	}

	return doctor_name, total, nil
}

func (s *hospitalService) GetCurrentQueueService(db *gorm.DB) (list []response.CurrentQueue, err error) {
	if err = db.Table("t_queue q").
		Joins("JOIN t_patients p ON p.nik = q.patient_nik").
		Joins("JOIN t_doctors d ON d.employee_id = q.doctor_employee_id").
		Where("q.status = 'Belum dilayani'").
		Select("q.id queue_no, p.name patient_name, d.name doctor_name").
		Order("1").
		Limit(1).
		Scan(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (s *hospitalService) GetQueueListService(db *gorm.DB) (list []response.Queue, total int64, err error) {
	if err = db.Raw(constant.QueueQuery).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	if err = db.Raw(constant.QueueQueryTotal).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *hospitalService) ValidateUpdate(queue *request.Queue, db *gorm.DB) (patients []*response.UnhandledPatients, total int64, err error) {
	if err = db.Model(dbmodel.Queue{}).Where("id < ? AND status = 'Belum dilayani'", queue.QueueNo).Count(&total).Error; err != nil {
		common_log.Error(err)
		return
	}
	if err = db.Model(dbmodel.Queue{}).
		Select("t_queue.id queue_no, p.name name, t_queue.status status").
		Where("t_queue.id < ? AND t_queue.status = 'Belum dilayani'", queue.QueueNo).
		Joins("JOIN t_patients p ON p.nik = t_queue.patient_nik").
		Order("1").Scan(&patients).Error; err != nil {
		common_log.Error(err)
		return
	}

	return patients, total, nil
}

func (s *hospitalService) UpdateQueueStatusService(queue *request.Queue, db *gorm.DB) (patient_name string, err error) {
	var patient_nik string
	if err = db.Model(dbmodel.Queue{}).Where("id = ?", queue.QueueNo).Pluck("patient_nik", &patient_nik).Error; err != nil {
		common_log.Error(err)
		return
	}

	if err = db.Model(dbmodel.Patients{}).Where("nik = ?", patient_nik).Pluck("name", &patient_name).Error; err != nil {
		common_log.Error(err)
		return
	}

	if err = db.Model(dbmodel.Queue{}).Where("id = ?", queue.QueueNo).Update("status", constant.PatientIsDone).Error; err != nil {
		common_log.Error(err)
		return
	}

	return patient_name, nil
}

func (s *hospitalService) PatientStatus(queue *request.Queue, db *gorm.DB) (patient_status string, total int64, err error) {
	if err = db.Model(dbmodel.Queue{}).Where("id = ?", queue.QueueNo).Pluck("status", &patient_status).Error; err != nil {
		common_log.Error(err)
		return
	}
	if err = db.Model(dbmodel.Queue{}).Where("id = ?", queue.QueueNo).Count(&total).Error; err != nil {
		common_log.Error(err)
		return
	}
	return patient_status, total, nil
}

func (s *hospitalService) GetQueueTotal(db *gorm.DB) (total int64, err error) {
	if err = db.Raw(constant.QueueQueryTotal).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (s *hospitalService) GetPatientTotal(db *gorm.DB) (total int64, err error) {
	if err = db.Model(dbmodel.Patients{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (s *hospitalService) GetDoctorTotal(db *gorm.DB) (total int64, err error) {
	if err = db.Model(dbmodel.Doctors{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
