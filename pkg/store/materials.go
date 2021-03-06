package store

import "strconv"

//ReadMaterials читает список материалов
func (mdb MDB) ReadMaterials(prices bool) ([]map[string]interface{}, error) {
	if prices {
		return mdb.ReadRows(GetMaterialsWithPricesQuery())
	}
	return mdb.ReadRows(GetMaterialsQuery())
}

//ReadMaterial читает материал по id
func (mdb MDB) ReadMaterial(prices bool, id int) (map[string]interface{}, error) {
	if prices {
		return mdb.ReadRow(GetMaterialWithPricesQuery(id))
	}
	return mdb.ReadRow(GetMaterialQuery(id))
}

//GetMaterialWithPricesQuery Получение материала с ценой id
func GetMaterialWithPricesQuery(id int) string {
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

//GetMaterialsWithPricesQuery Получение списка материалов с ценой
func GetMaterialsWithPricesQuery() string {
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

//GetMaterialQuery Получение материала по id
func GetMaterialQuery(id int) string {
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
		materials.coefficient
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

//GetMaterialsQuery Получение списка материалов
func GetMaterialsQuery() string {
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
		materials.coefficient
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
