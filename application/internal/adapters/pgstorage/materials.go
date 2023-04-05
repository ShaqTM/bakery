package pgstorage

import (
	"bakery/application/internal/domain/models"
	"strconv"
	"time"
)

// ReadMaterials читает список материалов
func (s *Storage) ReadMaterials(prices bool) ([]models.Material, error) {
	queryText := ""
	if prices {
		queryText = getMaterialsWithPricesQuery()
	} else {
		queryText = getMaterialsQuery()
	}
	materials := make([]models.Material, 0)
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return nil, err
	}
	coefficient := ""
	price := ""
	for rows.Next() {
		material := models.Material{}
		if err := rows.Scan(&material.Id, &material.Name, &material.Price_unit_id,
			&material.Price_unit_name, &material.Price_unit_short_name,
			&material.Recipe_unit_id, &material.Recipe_unit_name,
			&material.Recipe_unit_short_name,
			&coefficient, &price); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return nil, err
		}
		material.Coefficient = convertNumeric(coefficient)
		material.Price = convertNumeric(price)
		materials = append(materials, material)
	}
	return materials, nil
}

// ReadMaterial читает материал по id
func (s *Storage) ReadMaterial(prices bool, id int) (models.Material, error) {
	queryText := ""
	if prices {
		queryText = getMaterialWithPricesQuery(id)
	} else {
		queryText = getMaterialQuery(id)
	}
	rows, err := s.Pdb.Query(queryText)
	material := models.Material{}
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return models.Material{}, err
	}
	coefficient := ""
	price := ""
	if rows.Next() {

		if err := rows.Scan(&material.Id, &material.Name, &material.Price_unit_id,
			&material.Price_unit_name, &material.Price_unit_short_name,
			&material.Recipe_unit_id, &material.Recipe_unit_name,
			&material.Recipe_unit_short_name,
			&coefficient, &price); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return models.Material{}, err
		}
		material.Coefficient = convertNumeric(coefficient)
		material.Price = convertNumeric(price)
	}
	return material, nil
}

func (s *Storage) WriteMaterial(material models.Material) (int, error) {
	newRecord := material.Id == -1
	lastInsertID := -1
	if newRecord {
		queryText := `INSERT INTO public.materials (name,recipe_unit_id,price_unit_id,coefficient)
									 VALUES($1,$2,$3,$4) RETURNING id;`
		err := s.Pdb.QueryRow(queryText, material.Name, material.Recipe_unit_id, material.Price_unit_id, material.Coefficient).Scan(&lastInsertID)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(material)
			return -1, err
		}
	} else {
		queryText := `UPDATE public.materials SET (name,recipe_unit_id,price_unit_id,coefficient) = ($1,$2,$3,$4) WHERE id = ` + strconv.Itoa(material.Id) + `;`
		_, err := s.Pdb.Exec(queryText, material.Name, material.Recipe_unit_id, material.Price_unit_id, material.Coefficient)
		if err != nil {
			s.Log.Error("Error updating data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(material)
			return -1, err
		}
		lastInsertID = material.Id
	}
	return lastInsertID, nil
}

func (s *Storage) WriteMaterialPrice(material_price models.Material_price) (int, error) {
	lastInsertID := -1
	queryText := `INSERT INTO public.material_prices (material_id,date,price)
	VALUES($1,$2,$3) RETURNING id;`
	err := s.Pdb.QueryRow(queryText, material_price.Material_id, time.Now(), material_price.Price).Scan(&lastInsertID)
	if err != nil {
		s.Log.Error("Error inserting data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		s.Log.Error(material_price)
		return -1, err
	}
	return lastInsertID, nil

}

// GetMaterialWithPricesQuery Получение материала с ценой id
func getMaterialWithPricesQuery(id int) string {
	return `
	CREATE TEMP TABLE price_periods 
	ON COMMIT DROP
		AS
	SELECT
		material_prices.material_id,
		MAX(material_prices.date) AS date
	FROM 
		public.material_prices AS material_prices
	WHERE material_prices.material_id = ` + strconv.Itoa(id) + `
	GROUP BY material_prices.material_id;
	CREATE TEMP TABLE prices 
	ON COMMIT DROP
		AS
	SELECT
		material_prices.material_id,
		material_prices.price
	FROM 
		public.material_prices AS material_prices
			INNER JOIN price_periods 
			ON price_periods.material_id = material_prices.material_id
			AND price_periods.date = material_prices.date;
	
	SELECT 
		materials.id,
		materials.name,
		materials.price_unit_id,
		COALESCE(units1.name,'') AS price_unit_name,
		COALESCE(units1.short_name,'') AS price_unit_short_name,
		materials.recipe_unit_id,
		COALESCE(units2.name,'') AS recipe_unit_name,
		COALESCE(units2.short_name,'') AS recipe_unit_short_name,
		materials.coefficient,
		COALESCE(prices.price,0) AS price
	FROM 
		public.materials AS materials 
			LEFT JOIN public.units AS units1
			ON units1.id  = materials.price_unit_id
				LEFT JOIN public.units AS units2
				ON units2.id  = materials.recipe_unit_id
					LEFT JOIN prices 
					ON 	materials.id = prices.material_id
	WHERE 
		materials.id = ` + strconv.Itoa(id) + `
	ORDER BY 
		materials.id
	;`

}

// GetMaterialsWithPricesQuery Получение списка материалов с ценой
func getMaterialsWithPricesQuery() string {
	return `
	CREATE TEMP TABLE price_periods 
	ON COMMIT DROP
		AS
	SELECT
		material_prices.material_id,
		MAX(material_prices.date) AS date
	FROM 
		public.material_prices AS material_prices
	GROUP BY material_prices.material_id;
	CREATE TEMP TABLE prices 
	ON COMMIT DROP
		AS
	SELECT
		material_prices.material_id,
		material_prices.price
	FROM 
		public.material_prices AS material_prices
			INNER JOIN price_periods 
			ON price_periods.material_id = material_prices.material_id
			AND price_periods.date = material_prices.date;
	
	SELECT 
		materials.id,
		materials.name,
		materials.price_unit_id,
		COALESCE(units1.name,'') AS price_unit_name,
		COALESCE(units1.short_name,'') AS price_unit_short_name,
		materials.recipe_unit_id,
		COALESCE(units2.name,'') AS recipe_unit_name,
		COALESCE(units2.short_name,'') AS recipe_unit_short_name,
		materials.coefficient,
		COALESCE(prices.price,0) AS price
	FROM 
		public.materials AS materials 
			LEFT JOIN public.units AS units1
			ON units1.id  = materials.price_unit_id
				LEFT JOIN public.units AS units2
				ON units2.id  = materials.recipe_unit_id
					LEFT JOIN prices 
					ON 	materials.id = prices.material_id
	ORDER BY 
		materials.id
	;`

}

// GetMaterialQuery Получение материала по id
func getMaterialQuery(id int) string {
	return `
	SELECT 
		materials.id,
		materials.name,
		materials.price_unit_id,
		COALESCE(units1.name,'') AS price_unit_name,
		COALESCE(units1.short_name,'') AS price_unit_short_name,
		materials.recipe_unit_id,
		COALESCE(units2.name,'') AS recipe_unit_name,
		COALESCE(units2.short_name,'') AS recipe_unit_short_name,
		materials.coefficient,
		0 AS price
	FROM 
		public.materials AS materials 
			LEFT JOIN public.units AS units1
			ON units1.id  = materials.price_unit_id
				LEFT JOIN public.units AS units2
				ON units2.id  = materials.recipe_unit_id
	WHERE 
		materials.id = ` + strconv.Itoa(id) + `
	ORDER BY 
		materials.id
	;`

}

// GetMaterialsQuery Получение списка материалов
func getMaterialsQuery() string {
	return `
	SELECT 
		materials.id,
		materials.name,
		materials.price_unit_id,
		COALESCE(units1.name,'') AS price_unit_name,
		COALESCE(units1.short_name,'') AS price_unit_short_name,
		materials.recipe_unit_id,
		COALESCE(units2.name,'') AS recipe_unit_name,
		COALESCE(units2.short_name,'') AS recipe_unit_short_name,
		materials.coefficient,
		0 AS price
	FROM 
		public.materials AS materials 
			LEFT JOIN public.units AS units1
			ON units1.id  = materials.price_unit_id
				LEFT JOIN public.units AS units2
				ON units2.id  = materials.recipe_unit_id
	ORDER BY 
		materials.id
	;`

}
