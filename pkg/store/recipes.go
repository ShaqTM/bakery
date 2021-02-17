package store

import "strconv"

//ReadRecipesList Читает список рецептов
func (mdb MDB) ReadRecipesList() ([]map[string]interface{}, error) {
	params := make(map[string]string)
	params["order"] = "name"
	return mdb.ReadRows(GetRowsQuerry("recipes", params))
}

//ReadRecipe Читает рецепт по id
func (mdb MDB) ReadRecipe(id int) (map[string]interface{}, error) {
	data, err := mdb.ReadRow(GetRowQuerry("recipes", id))
	if err != nil {
		return data, err
	}
	content, err := mdb.ReadRecipeContent(id)
	if err != nil {
		return data, err
	}
	data["content"] = content
	return data, nil
}

//ReadRecipeContent Читает табличную часть рецепта по id
func (mdb MDB) ReadRecipeContent(id int) ([]map[string]interface{}, error) {
	return mdb.ReadRows(GetRecipeContentQuery(id))
}

//ReadRecipeContentWithPrice Читает табличную часть рецепта с ценами по id
func (mdb MDB) ReadRecipeContentWithPrice(id int) ([]map[string]interface{}, error) {
	return mdb.ReadRows(GetRecipeContentWithPricesQuery(id))
}

//RecipesQuery Текст запроса для получения списка рецептов
const RecipesQuery = `
	SELECT 
		recipes.id,
		recipes.name,
		recipes.output,
		recipes.unit_id,
		units.short_name AS unit_short_name
	FROM 
		public.recipes AS recipes 
			LEFT JOIN public.units AS units 
			ON recipes.unit_id = units.id
	
	ORDER BY 
	recipes.id;`

//GetRecipeQuerry Получение рецепта по id
func GetRecipeQuerry(id int) string {
	return `
	SELECT 
		recipes.id,
		recipes.name,
		recipes.output,
		recipes.unit_id,
		units.short_name AS unit_short_name
	FROM 
		public.recipes AS recipes 
			LEFT JOIN public.units AS units 
			ON recipes.unit_id = units.id
		
	WHERE 
		recipes.id = ` + strconv.Itoa(id) + `;`
}

//GetRecipeContentQuery Получение таблицы рецепта по id
func GetRecipeContentQuery(id int) string {
	return `
	SELECT 
		recipes_content.id,
		recipes_content.material_id,
		materials.name AS material_name,
		recipes_content.qty,
		recipes_content.string_order,
		units.name AS unit_name,
		units.short_name AS unit_short_name
	FROM 
		public.recipes_content AS recipes_content 
			LEFT JOIN public.materials AS materials 
			ON recipes_content.material_id = materials.id 
				 LEFT JOIN public.units AS units 
				 ON units.id  = materials.recipe_unit_id
	WHERE 
		recipes_content.id = ` + strconv.Itoa(id) + `
	ORDER BY 
		recipes_content.string_order
	;`
}

//GetRecipeContentWithPricesQuery Получение таблицы рецепта по id
func GetRecipeContentWithPricesQuery(id int) string {
	return `
	CREATE TEMP TABLE temp_materials 
	ON COMMIT DROP
		AS
	SELECT 
		recipes_content.material_id
	FROM  
		public.recipes_content AS recipes_content
	WHERE recipes_content.id = ` + strconv.Itoa(id) + `;
	CREATE TEMP TABLE price_periods AS
	SELECT
		material_prices.material_id,
		MAX(material_prices.date) AS date
	FROM 
		public.material_prices AS material_prices
			INNER JOIN temp_materials 
			ON temp_materials.material_id = material_prices.material_id
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
		recipes_content.id,
		recipes_content.material_id,
		materials.name AS material_name,
		recipes_content.qty,
		recipes_content.string_order,
		units.name AS unit_name,
		units.short_name AS unit_short_name,
		COALESCE(prices.price,0) AS price,
		materials.coefficient AS coefficient
	FROM 
		public.recipes_content AS recipes_content 
			LEFT JOIN public.materials AS materials 
			ON recipes_content.material_id = materials.id 
				LEFT JOIN public.units AS units 
				ON units.id  = materials.recipe_unit_id
					LEFT JOIN prices 
					ON 	recipes_content.material_id = prices.material_id
	WHERE 
		recipes_content.id = ` + strconv.Itoa(id) + `
	ORDER BY 
		recipes_content.string_order
	;`

}
