package store

//ReadUnits читает список единиц измерения
func (mdb MDB) ReadUnits() ([]map[string]interface{}, error) {
	params := make(map[string]string)
	params["order"] = "name"
	return mdb.ReadRows(GetRowsQuerry("units", params))
}

//ReadUnit читает единицу измерения
func (mdb MDB) ReadUnit(id int) (map[string]interface{}, error) {
	return mdb.ReadRow(GetRowQuerry("units", id))
}
