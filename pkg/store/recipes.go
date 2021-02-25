package store

import "strconv"

//ReadRecipes Читает список рецептов
func (mdb MDB) ReadRecipes(prices bool) ([]map[string]interface{}, error) {
	if prices {
		return mdb.ReadRows(RecipesWithPricesQuery)
	}
	return mdb.ReadRows(RecipesQuery)
}

//ReadRecipe Читает рецепт по id
func (mdb MDB) ReadRecipe(prices bool, id int) (map[string]interface{}, error) {
	queryText := ""
	if prices {
		queryText = GetRecipeWithPriceQuerry(id)
	} else {
		queryText = GetRecipeQuerry(id)
	}
	data, err := mdb.ReadRow(queryText)
	if err != nil {
		return data, err
	}

	content, err := mdb.ReadRecipeContent(prices, id)
	if err != nil {
		return data, err
	}
	data["content"] = content
	return data, nil
}

//ReadRecipeContent Читает табличную часть рецепта по id
func (mdb MDB) ReadRecipeContent(prices bool, id int) ([]map[string]interface{}, error) {
	if prices {
		return mdb.ReadRows(GetRecipeContentWithPricesQuery(id))
	}
	return mdb.ReadRows(GetRecipeContentQuery(id))
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
	recipes.name;`

//RecipesWithPricesQuery Текст запроса для получения списка рецептов с ценами
const RecipesWithPricesQuery = `
	CREATE TEMP TABLE price_periods 
	ON COMMIT DROP
		AS
	SELECT
		recipe_prices.recipe_id,
		MAX(recipe_prices.date) AS date
	FROM 
		public.recipe_prices AS recipe_prices
	GROUP BY recipe_prices.recipe_id;
	CREATE TEMP TABLE prices 
	ON COMMIT DROP
		AS
	SELECT
		recipe_prices.recipe_id,
		recipe_prices.price
	FROM 
		public.recipe_prices AS recipe_prices
		INNER JOIN price_periods 
			ON price_periods.recipe_id = recipe_prices.recipe_id
			AND price_periods.date = recipe_prices.date;
	SELECT 
		recipes.id,
		recipes.name,
		recipes.output,
		recipes.unit_id,
		units.short_name AS unit_short_name,
		COALESCE(prices.price,0) AS price

	FROM 
		public.recipes AS recipes 
			LEFT JOIN public.units AS units 
			ON recipes.unit_id = units.id
				LEFT JOIN prices 
				ON 	recipes.id = prices.recipe_id
	
	ORDER BY 
		recipes.name;`

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

//GetRecipeWithPriceQuerry Получение рецепта c ценой по id
func GetRecipeWithPriceQuerry(id int) string {
	return `
	CREATE TEMP TABLE price_periods 
	ON COMMIT DROP
		AS
	SELECT
		recipe_prices.recipe_id,
		MAX(recipe_prices.date) AS date
	FROM 
		public.recipe_prices AS recipe_prices
	WHERE 
	recipe_prices.recipe_id = ` + strconv.Itoa(id) + `
	GROUP BY recipe_prices.recipe_id;
	CREATE TEMP TABLE prices 
	ON COMMIT DROP
		AS
	SELECT
		recipe_prices.recipe_id,
		recipe_prices.price
	FROM 
		public.recipe_prices AS recipe_prices
		INNER JOIN price_periods 
			ON price_periods.recipe_id = recipe_prices.recipe_id
			AND price_periods.date = recipe_prices.date;
	SELECT 
		recipes.id,
		recipes.name,
		recipes.output,
		recipes.unit_id,
		units.short_name AS unit_short_name,
		COALESCE(prices.price,0) AS price
	FROM 
		public.recipes AS recipes 
			LEFT JOIN public.units AS units 
			ON recipes.unit_id = units.id
				LEFT JOIN prices 
				ON 	recipes.id = prices.recipe_id
		
	WHERE 
		recipes.id = ` + strconv.Itoa(id) + `;`
}

//GetRecipeContentQuery Получение таблицы рецепта по id
func GetRecipeContentQuery(id int) string {
	return `
	SELECT 
		recipes_content.id AS recipe_id,
		recipes_content.material_id,
		materials.name AS material_name,
		recipes_content.qty,
		recipes_content.string_order,
		materials.recipe_unit_id AS unit_id,
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
	CREATE TEMP TABLE price_periods 
	ON COMMIT DROP
	AS
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
		recipes_content.id AS recipe_id,
		recipes_content.material_id,
		materials.name AS material_name,
		recipes_content.qty,
		recipes_content.string_order,
		materials.recipe_unit_id AS unit_id,
		units.name AS unit_name,
		units.short_name AS unit_short_name,
		materials.coefficient AS coefficient,
		CASE WHEN materials.coefficient = 0 THEN 0 
		ELSE COALESCE(prices.price,0)/materials.coefficient
		END AS price
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
