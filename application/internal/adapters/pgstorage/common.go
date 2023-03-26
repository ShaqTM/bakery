package pgstorage

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Storage Структура, содержит ссылку на интерфейс к ДБ
type Storage struct {
	Log *logrus.Logger
	Pdb *sql.DB
}

func New(log *logrus.Logger) *Storage {
	s := &Storage{
		Log: log,
	}
	return s
}

var numericCols = []string{"coefficient", "price", "output", "qty", "plan_qty", "plan_cost", "fact_qty", "fact_cost", "materials_cost", "cost"}

func convertNumeric(value interface{}) float64 {
	switch d := value.(type) {
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

// GetRowQuerry Получение записи по имени таблицы и ID
func GetRowQuerry(tableName string, id int) string {
	return `SELECT * FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`
}

// GetRowsQuerry Получение записей по имени таблицы
func GetRowsQuerry(tableName string, params map[string]string) string {
	orderString := ""
	if params["order"] != "" {
		orderString = ` ORDER BY ` + params["order"]
	}
	return `SELECT * FROM public.` + tableName + orderString + `;`
}

// UpdateData Запись данных в таблицу
func (s *Storage) UpdateData(tableName string, data map[string]interface{}) (int, error) {
	queryText := `SELECT * FROM ` + tableName + ` WHERE false`
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading columns:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return -1, err
	}
	cols, _ := rows.Columns()

	vals := []interface{}{}
	namesString := ""
	placeholdersString := ""
	newRecord := false
	valueIndex := 1
	for _, col := range cols {
		if col == "id" && int((data[col]).(float64)) == -1 {
			newRecord = true
			continue
		}
		if data[col] == nil {
			continue
		}
		namesString += (col + ",")
		placeholdersString += ("$" + strconv.Itoa(valueIndex) + ",")
		vals = append(vals, data[col])
		valueIndex++
	}
	namesString = strings.TrimSuffix(namesString, ",")
	placeholdersString = strings.TrimSuffix(placeholdersString, ",")

	if newRecord {
		lastInsertID := -1
		queryText := `INSERT INTO public.` + tableName + `(` + namesString + `) VALUES(` + placeholdersString + `) RETURNING id;`
		err := s.Pdb.QueryRow(queryText, vals...).Scan(&lastInsertID)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(vals...)
			return -1, err
		}
		return lastInsertID, nil
	}
	queryText = `UPDATE public.` + tableName + ` SET (` + namesString + `) = (` + placeholdersString + `) WHERE id = ` + strconv.Itoa(int(data["id"].(float64))) + `;`
	_, err = s.Pdb.Exec(queryText, vals...)
	if err != nil {
		s.Log.Error("Error updating data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		s.Log.Error(vals...)
		return -1, err
	}
	return int(data["id"].(float64)), nil

}

// UpdateTableData Обновляет данные в табличной части
func (s *Storage) UpdateTableData(tableName string, data []interface{}, id int) error {
	queryText := `SELECT * FROM ` + tableName + ` WHERE false`
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading columns:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return err
	}
	cols, _ := rows.Columns()
	queryText = `DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`
	_, err = s.Pdb.Exec(queryText)
	if err != nil {
		s.Log.Error("Error deleting old rows:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return err
	}
	if len(data) == 0 {
		return nil
	}
	namesString := ""
	placeholdersString := ""
	valueIndex := 1
	for _, col := range cols {
		if col != "id" && (data[0].(map[string]interface{}))[col] == nil {
			continue
		}
		namesString += (col + ",")
		placeholdersString += ("$" + strconv.Itoa(valueIndex) + ",")
		valueIndex++
	}
	namesString = strings.TrimSuffix(namesString, ",")
	placeholdersString = strings.TrimSuffix(placeholdersString, ",")
	queryText = `insert into public.` + tableName + `(` + namesString + `) VALUES(` + placeholdersString + `);`
	stmt, err := s.Pdb.Prepare(queryText)
	if err != nil {
		s.Log.Error("Error preparing for inserting data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return err
	}
	for _, item := range data {
		vals := []interface{}{}
		for _, col := range cols {
			if col == "id" {
				vals = append(vals, id)
				continue
			} else if item.(map[string]interface{})[col] == nil {
				continue
			}
			vals = append(vals, item.(map[string]interface{})[col])

		}
		_, err := stmt.Exec(vals...)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(vals...)

			return err
		}
	}
	return nil
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

// ReadRow Выполняет запрос к SQL и возвращает одну строку в виде map
func (s *Storage) ReadRow(queryText string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return data, err
	}
	cols, _ := Query.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	if !rows.Next() {
		return data, nil
	}
	// Scan the result into the column pointers...
	if err := rows.Scan(columnPointers...); err != nil {
		s.Log.Error("Error scanning rows:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return data, err
	}

	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		data[colName] = convertIfNumeric(colName, *val)
	}
	return data, nil
}

// ReadRows Выполняет запрос к SQL и возвращает строки в виде []map
func (s *Storage) ReadRows(queryText string) ([]map[string]interface{}, error) {

	dataArray := make([]map[string]interface{}, 0)
	rows, err := s.Pdb.Query(queryText) // Note: Ignoring errors for brevity
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return dataArray, err
	}
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	for rows.Next() {
		if err := rows.Scan(columnPointers...); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return dataArray, err
		}
		data := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			data[colName] = convertIfNumeric(colName, *val)
		}
		dataArray = append(dataArray, data)
	}
	return dataArray, nil
}
