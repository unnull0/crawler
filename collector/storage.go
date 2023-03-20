package collector

type Storage interface {
	Save(data ...*DataCell) error
}

type DataCell struct {
	Data map[string]interface{}
}

func (d *DataCell) GetTaskName() string {
	return d.Data["Task"].(string)
}

func (d *DataCell) GetTableName() string {
	return d.Data["Task"].(string)
}
