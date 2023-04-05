package pgstorage

import (
	"bakery/application/internal/domain/models"
	"strconv"
)

// ReadUnits читает список единиц измерения
func (s *Storage) ReadUnits() ([]models.Unit, error) {
	units := make([]models.Unit, 0)
	querryText := `
	SELECT 
		id,
		name,
		short_name
	FROM public.units
	ORDER BY name;`
	rows, err := s.Pdb.Query(querryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(querryText)
		return nil, err
	}

	id := 0
	name := ""
	shortName := ""
	for rows.Next() {

		if err := rows.Scan(&id, &name, &shortName); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(querryText)
			return nil, err
		}
		unit := models.Unit{
			Id:         id,
			Name:       name,
			Short_name: shortName,
		}
		units = append(units, unit)
	}
	return units, nil
}

// ReadUnit читает единицу измерения
func (s *Storage) ReadUnit(id int) (models.Unit, error) {
	querryText := `
	SELECT 
		id,
		name,
		short_name
	FROM public.units
	WHERE id=` + strconv.Itoa(id) + `
	ORDER BY name;`
	rows, err := s.Pdb.Query(querryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(querryText)
		return models.Unit{}, err
	}
	unitId := 0
	name := ""
	shortName := ""
	if err := rows.Scan(&id, &name, &shortName); err != nil {
		s.Log.Error("Error scanning rows:", err)
		s.Log.Error("Query:")
		s.Log.Error(querryText)
		return models.Unit{}, err
	}
	unit := models.Unit{
		Id:         unitId,
		Name:       name,
		Short_name: shortName,
	}
	return unit, nil
}

func (s *Storage) WriteUnit(unit models.Unit) (int, error) {
	s.Log.Info("Writing unit", unit)
	newRecord := unit.Id == -1
	if newRecord {
		lastInsertID := -1
		queryText := `INSERT INTO public.units (name, short_name)
					 VALUES($1,$2) RETURNING id;`
		err := s.Pdb.QueryRow(queryText, unit.Name, unit.Short_name).Scan(&lastInsertID)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(unit)
			return -1, err
		}
		return lastInsertID, nil
	} else {
		queryText := `UPDATE public.units SET(name, short_name)=($1,$2) WHERE id = ` + strconv.Itoa(unit.Id) + `;`
		_, err := s.Pdb.Exec(queryText, unit.Name, unit.Short_name)
		if err != nil {
			s.Log.Error("Error updating data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(unit)
			return -1, err
		}
		return unit.Id, nil
	}
}
