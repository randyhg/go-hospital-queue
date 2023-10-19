package request

type Queue struct {
	QueueNo string `json:"queue_no"`
}

type Patients struct {
	Name             string `json:"name"`
	Birthdate        string `json:"birthdate"`
	NIK              string `json:"nik"`
	Sex              string `json:"sex"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	EmergencyContact string `json:"emergency_contact"`
}

type Doctors struct {
	Name           string `json:"name"`
	EmployeeId     string `json:"employee_id"`
	Specialization string `json:"specialization"`
	Phone          string `json:"phone"`
	WorkDay        string `json:"work_day"`
}
