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
	Pdb **sql.DB
}

var numericCols = []string{"coefficient", "price", "output", "qty", "plan_qty", "plan_cost", "fact_qty", "fact_cost", "materials_cost", "cost"}

func convertIfNumeric(col string, value interface{}) interface{} {
	for i := range numericCols {
		if numericCols[i] == col {
			//			fmt.Println(col)
			res, err := strconv.ParseFloat(string(value.([]uint8)), 64)

			if err != nil {
				fmt.Println("Error parsing float:", value, err)
				return value
			}
			return res

		}
	}
	return value
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
func (mdb MDB) UpdateData(tableName string, data map[string]interface{}) (int, error) {
	db := *mdb.Pdb
	queryText := `SELECT * FROM ` + tableName + ` WHERE false`
	rows, err := db.Query(queryText)
	if err != nil {
		mdb.Log.Error("Error reading columns:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
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
		err := db.QueryRow(queryText, vals...).Scan(&lastInsertID)
		if err != nil {
			mdb.Log.Error("Error inserting data:", err)
			mdb.Log.Error("Query:")
			mdb.Log.Error(queryText)
			mdb.Log.Error(vals...)
			return -1, err
		}
		return lastInsertID, nil
	}
	queryText = `UPDATE public.` + tableName + ` SET (` + namesString + `) = (` + placeholdersString + `) WHERE id = ` + strconv.Itoa(int(data["id"].(float64))) + `;`
	_, err = db.Exec(queryText, vals...)
	if err != nil {
		mdb.Log.Error("Error updating data:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
		mdb.Log.Error(vals...)
		return -1, err
	}
	return int(data["id"].(float64)), nil

}

// UpdateTableData Обновляет данные в табличной части
func (mdb MDB) UpdateTableData(tableName string, data []interface{}, id int) error {
	db := *mdb.Pdb
	queryText := `SELECT * FROM ` + tableName + ` WHERE false`
	rows, err := db.Query(queryText)
	if err != nil {
		mdb.Log.Error("Error reading columns:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
		return err
	}
	cols, _ := rows.Columns()
	queryText = `DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`
	_, err = db.Exec(queryText)
	if err != nil {
		mdb.Log.Error("Error deleting old rows:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
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
	stmt, err := db.Prepare(queryText)
	if err != nil {
		mdb.Log.Error("Error preparing for inserting data:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
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
			mdb.Log.Error("Error inserting data:", err)
			mdb.Log.Error("Query:")
			mdb.Log.Error(queryText)
			mdb.Log.Error(vals...)

			return err
		}
	}
	return nil
}

// DeleteData удaление данных в таблице по id
func (mdb MDB) DeleteData(tableName string, id int) error {
	db := *mdb.Pdb
	queryText := `DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`
	_, err := db.Exec(queryText)
	if err != nil {
		mdb.Log.Error("Error deleting old rows:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
		return err
	}
	return nil
}

// ReadRow Выполняет запрос к SQL и возвращает одну строку в виде map
func (mdb MDB) ReadRow(queryText string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	db := *mdb.Pdb
	rows, err := db.Query(queryText)
	if err != nil {
		mdb.Log.Error("Error reading data:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
		return data, err
	}
	cols, _ := rows.Columns()
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
		mdb.Log.Error("Error scanning rows:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
		return data, err
	}

	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		data[colName] = convertIfNumeric(colName, *val)
	}
	return data, nil
}

// ReadRows Выполняет запрос к SQL и возвращает строки в виде []map
func (mdb MDB) ReadRows(queryText string) ([]map[string]interface{}, error) {

	dataArray := make([]map[string]interface{}, 0)
	db := *mdb.Pdb
	rows, err := db.Query(queryText) // Note: Ignoring errors for brevity
	if err != nil {
		mdb.Log.Error("Error reading data:", err)
		mdb.Log.Error("Query:")
		mdb.Log.Error(queryText)
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
			mdb.Log.Error("Error scanning rows:", err)
			mdb.Log.Error("Query:")
			mdb.Log.Error(queryText)
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
