package constant

const (
	PatientIsDone = "Sudah dilayani"
)

const (
	QueueQuery = `SELECT q.id queue_no, p.name patient_name, d.name doctor_name, q.queue_date date, status 
			FROM t_queue q 
			JOIN t_patients p ON p.nik = q.patient_nik 
			JOIN t_doctors d ON d.employee_id = q.doctor_employee_id 
			WHERE status = 'Belum dilayani' 
			ORDER BY 1;`

	QueueQueryTotal = `SELECT COUNT(*) FROM t_queue WHERE status = 'Belum dilayani';`
)
