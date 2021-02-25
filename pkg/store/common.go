package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

//MDB Структура, содержит ссылку на интерфейс к ДБ
type MDB struct {
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

//GetRowQuerry Получение записи по имени таблицы и ID
func GetRowQuerry(tableName string, id int) string {
	return `SELECT * FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`
}

//GetRowsQuerry Получение записей по имени таблицы
func GetRowsQuerry(tableName string, params map[string]string) string {
	orderString := ""
	if params["order"] != "" {
		orderString = ` ORDER BY ` + params["order"]
	}
	return `SELECT * FROM public.` + tableName + orderString + `;`
}

//UpdateData Запись данных в таблицу
func (mdb MDB) UpdateData(tableName string, data map[string]interface{}) (int, error) {
	db := *mdb.Pdb

	rows, err := db.Query(`SELECT * FROM ` + tableName + ` WHERE false`)
	if err != nil {
		fmt.Println("Error reading data: ", err)
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
		err := db.QueryRow(`insert into public.`+tableName+`(`+namesString+`) VALUES(`+placeholdersString+`) RETURNING id;`, vals...).Scan(&lastInsertID)
		if err != nil {
			fmt.Println("Error inserting data: ", err)
			fmt.Println(`insert into public.` + tableName + `(` + namesString + `) VALUES(` + placeholdersString + `) RETURNING id;`)
			fmt.Println(vals)
			return -1, err
		}
		return lastInsertID, nil
	}
	_, err = db.Exec(`UPDATE public.`+tableName+` SET (`+namesString+`) = (`+placeholdersString+`) WHERE id = `+strconv.Itoa(int(data["id"].(float64)))+`;`, vals...)
	if err != nil {
		fmt.Println("Error updating data: ", err)
		return -1, err
	}
	return int(data["id"].(float64)), nil

}

//UpdateTableData Обновляет данные в табличной части
func (mdb MDB) UpdateTableData(tableName string, data []interface{}, id int) error {
	db := *mdb.Pdb
	rows, err := db.Query(`SELECT * FROM ` + tableName + ` WHERE false`)
	if err != nil {
		fmt.Println("Error reading data: ", err)
		return err
	}
	cols, _ := rows.Columns()

	_, err = db.Exec(`DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`)
	if err != nil {
		fmt.Println("Error deleting old rows: ", err)
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
	stmt, err := db.Prepare(`insert into public.` + tableName + `(` + namesString + `) VALUES(` + placeholdersString + `);`)
	if err != nil {
		fmt.Println("Error preparing for inserting data: ", err)
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
			fmt.Println("Error inserting data: ", err)
			return err
		}
	}
	return nil
}

//DeleteData уделение данных в таблице по id
func (mdb MDB) DeleteData(tableName string, id int) error {
	db := *mdb.Pdb
	_, err := db.Exec(`DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`)
	if err != nil {
		fmt.Println("Error deleting old rows: ", err)
		return err
	}
	return nil
}

//ReadRow Выполняет запрос к SQL и возвращает одну строку в виде map
func (mdb MDB) ReadRow(queryText string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	db := *mdb.Pdb
	rows, err := db.Query(queryText)
	if err != nil {
		fmt.Println("Error reading data: ", err)
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
		fmt.Println("Error reading data: ", err)
		return data, err
	}

	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		data[colName] = convertIfNumeric(colName, *val)
	}
	return data, nil
}

//ReadRows Выполняет запрос к SQL и возвращает строки в виде []map
func (mdb MDB) ReadRows(queryText string) ([]map[string]interface{}, error) {

	dataArray := make([]map[string]interface{}, 0)
	db := *mdb.Pdb
	rows, err := db.Query(queryText) // Note: Ignoring errors for brevity
	if err != nil {
		fmt.Println("Error reading data: ", err)
		fmt.Println(queryText)
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
			fmt.Println("Error reading data: ", err)
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
