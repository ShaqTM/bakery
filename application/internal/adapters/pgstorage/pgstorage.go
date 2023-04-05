package pgstorage

import (
	"bakery/application/internal/config"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Storage Структура, содержит ссылку на интерфейс к ДБ
type Storage struct {
	Log    *logrus.Logger
	Pdb    *sql.DB
	config *config.Config
}

func New(log *logrus.Logger, config *config.Config) *Storage {
	s := &Storage{
		Log:    log,
		config: config,
	}
	return s
}

func convertNumeric(value interface{}) float64 {
	switch value.(type) {
	case string:
		res, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			fmt.Println("Error parsing float:", value, err)
			return 0
		}
		return res
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	case int:
		return float64(value.(int))
	}

	return 0
}

// DeleteData удaление данных в таблице по id
func (s *Storage) DeleteData(tableName string, id int) error {
	queryText := `DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`
	_, err := s.Pdb.Exec(queryText)
	if err != nil {
		s.Log.Error("Error deleting old rows:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return err
	}
	return nil
}
