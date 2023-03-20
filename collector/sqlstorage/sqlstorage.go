package sqlstorage

import (
	"encoding/json"

	"github.com/unnull0/crawler/collector"
	"github.com/unnull0/crawler/grabber/workerengine"
	"github.com/unnull0/crawler/mysqldb"
	"go.uber.org/zap"
)

type SqlStore struct {
	dataDocker  []*collector.DataCell
	columnNames []mysqldb.Field
	db          mysqldb.DBer
	Table       map[string]struct{}
	options
}

func New(opts ...Option) (*SqlStore, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &SqlStore{}
	s.options = options
	s.Table = make(map[string]struct{})

	var err error
	s.db, err = mysqldb.New(
		mysqldb.WithConnUrl(s.sqlUrl),
		mysqldb.WithLogger(s.logger),
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SqlStore) Save(dataCells ...*collector.DataCell) error {
	for _, data := range dataCells {
		name := data.GetTableName()
		if _, ok := s.Table[name]; !ok {
			columnNames := getFields(data)

			err := s.db.CreateTable(mysqldb.TableData{
				TableName:   name,
				ColumnNames: columnNames,
				AutoKey:     true,
			})
			if err != nil {
				s.logger.Error("create table failed", zap.Error(err))
			}
			s.Table[name] = struct{}{}
		}
		if len(s.dataDocker) >= s.BatchCount {
			s.Flush()
		}
		s.dataDocker = append(s.dataDocker, data)
	}
	return nil
}

func getFields(data *collector.DataCell) []mysqldb.Field {
	taskName := data.Data["Task"].(string)
	ruleName := data.Data["Rule"].(string)
	fields := workerengine.GetFields(taskName, ruleName)

	var columnNames = []mysqldb.Field{}
	for _, field := range fields {
		columnNames = append(columnNames, mysqldb.Field{
			Title: field,
			Type:  "MEDIUMTEXT",
		})
	}
	columnNames = append(columnNames,
		mysqldb.Field{Title: "Url", Type: "VARCHAR(255)"},
		mysqldb.Field{Title: "Time", Type: "VARCHAR(255)"},
	)
	return columnNames
}

func (s *SqlStore) Flush() error {
	if len(s.dataDocker) == 0 {
		return nil
	}
	args := make([]interface{}, 0)
	for _, dataCell := range s.dataDocker {
		taskName := dataCell.Data["Task"].(string)
		ruleName := dataCell.Data["Rule"].(string)
		fields := workerengine.GetFields(taskName, ruleName)
		data := dataCell.Data["Data"].(map[string]interface{})
		value := []string{}
		for _, field := range fields {
			v := data[field]
			switch v.(type) {
			case nil:
				value = append(value, "")
			case string:
				value = append(value, v.(string))
			default:
				j, err := json.Marshal(v)
				if err != nil {
					value = append(value, "")
				} else {
					value = append(value, string(j))
				}
			}
		}
		value = append(value, dataCell.Data["Url"].(string), dataCell.Data["Time"].(string))
		for _, v := range value {
			args = append(args, v)
		}
	}

	return s.db.Insert(mysqldb.TableData{
		TableName:   s.dataDocker[0].GetTableName(),
		ColumnNames: getFields(s.dataDocker[0]),
		Args:        args,
		DataCount:   len(s.dataDocker),
	})
}
