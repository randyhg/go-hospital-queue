package dbmodel

const (
	patient_done = "Sudah dilayani"
)

type Model struct {
	Id int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

type Patients struct {
	Name             string `gorm:"index:idx_name;not null;" json:"name"`
	Birthdate        string `gorm:"index:idx_birthdate;default null" json:"birthdate"`
	NIK              string `gorm:"index:idx_nik;not null;" json:"nik"`
	Sex              string `gorm:"index:idx_sex;not null;" json:"sex"`
	Address          string `gorm:"index:idx_address;not null;" json:"address"`
	Phone            string `gorm:"index:idx_phone;not null" json:"phone"`
	EmergencyContact string `gorm:"index:idx_emergency_contact;" json:"emergency_contact"`
}

type Doctors struct {
	Name           string `gorm:"index:idx_name;not null;" json:"name"`
	EmployeeId     string `gorm:"index:idx_employee_id;" json:"employee_id"`
	Specialization string `gorm:"index:idx_specialization;default null" json:"specialization"`
	Phone          string `gorm:"index:idx_phone;not null" json:"phone"`
	WorkDay        string `gorm:"index:idx_work_day;" json:"work_day"`
}

type Queue struct {
	PatientNIK       string `gorm:"index:idx_patient_nik;not null" json:"patient_nik"`
	DoctorEmployeeId string `gorm:"index:idx_doctor_employee_id;not null" json:"doctor_employee_id"`
	QueueDate        string `gorm:"index:idx_queue_date" json:"queue_date"`
}
