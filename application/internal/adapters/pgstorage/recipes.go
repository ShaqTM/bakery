package pgstorage

import (
	"bakery/application/internal/domain/models"
	"strconv"
)

// ReadRecipes Читает список рецептов
func (s *Storage) ReadRecipes(prices bool) ([]models.Recipe, error) {
	queryText := ""
	if prices {
		queryText = getRecipesWithPricesQuery()
	} else {
		queryText = getRecipesQuery()
	}
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return nil, err
	}

	recipes := make([]models.Recipe, 0)
	output := ""
	price := ""
	for rows.Next() {
		recipe := models.Recipe{}
		if err := rows.Scan(&recipe.Id, &recipe.Name, &output,
			&recipe.Unit_id, &recipe.Unit_short_name,
			&price); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return nil, err
		}
		recipe.Output = convertNumeric(output)
		recipe.Price = convertNumeric(price)
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

// ReadRecipe Читает рецепт по id
func (s *Storage) ReadRecipe(prices bool, id int) (models.Recipe, error) {
	queryText := ""
	if prices {
		queryText = getRecipeWithPriceQuerry(id)
	} else {
		queryText = getRecipeQuerry(id)
	}
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return models.Recipe{}, err
	}
	recipe := models.Recipe{}
	if rows.Next() {
		output := ""
		price := ""
		if err := rows.Scan(&recipe.Id, &recipe.Name, &output,
			&recipe.Unit_id, &recipe.Unit_short_name,
			&price); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return models.Recipe{}, err
		}
		recipe.Output = convertNumeric(output)
		recipe.Price = convertNumeric(price)

		content, err := s.readRecipeContent(prices, id)
		if err != nil {
			return models.Recipe{}, err
		}
		recipe.Content = content
	}
	return recipe, nil
}

// ReadRecipeContent Читает табличную часть рецепта по id
func (s *Storage) readRecipeContent(prices bool, id int) ([]models.Recipe_content, error) {
	queryText := ""
	if prices {
		queryText = getRecipeContentWithPricesQuery(id)
	} else {
		queryText = getRecipeContentQuery(id)
	}
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return nil, err
	}

	contents := make([]models.Recipe_content, 0)
	coefficient := ""
	price := ""
	qty := ""

	for rows.Next() {
		content := models.Recipe_content{}
		if err := rows.Scan(&content.Id, &content.Material_id, &content.Material_name, &qty,
			&content.String_order, &content.Unit_id, &content.Unit_name, &content.Unit_short_name,
			&coefficient, &price); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return nil, err
		}
		content.Coefficient = convertNumeric(coefficient)
		content.Price = convertNumeric(price)
		content.Qty = convertNumeric(qty)
		contents = append(contents, content)
	}
	return contents, nil
}

func (s *Storage) WriteRecipe(recipe models.Recipe) (int, error) {
	newRecord := recipe.Id == -1
	lastInsertID := -1

	if newRecord {
		queryText := `INSERT INTO public.recipes (name,output,unit_id)
									 VALUES($1,$2,$3) RETURNING id;`
		err := s.Pdb.QueryRow(queryText, recipe.Name, recipe.Output, recipe.Unit_id).Scan(&lastInsertID)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(recipe)
			return -1, err
		}
	} else {
		queryText := `UPDATE public.recipes SET (name,output,unit_id) = ($1,$2,$3) WHERE id = ` + strconv.Itoa(recipe.Id) + `;`
		_, err := s.Pdb.Exec(queryText, recipe.Name, recipe.Output, recipe.Unit_id)
		if err != nil {
			s.Log.Error("Error updating data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(recipe)
			return -1, err
		}
		lastInsertID = recipe.Id
	}

	queryText := `DELETE FROM public.recipes_content WHERE id = ` + strconv.Itoa(lastInsertID) + `;`
	_, err := s.Pdb.Exec(queryText)
	if err != nil {
		s.Log.Error("Error deleting old rows:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return -1, err
	}
	if len(recipe.Content) == 0 {
		return lastInsertID, nil
	}

	queryText = `INSERT INTO public.recipes_content (id,material_id,string_order,qty) VALUES($1,$2,$3,$4);`
	stmt, err := s.Pdb.Prepare(queryText)
	if err != nil {
		s.Log.Error("Error preparing for inserting data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return -1, nil
	}
	for _, content := range recipe.Content {
		_, err := stmt.Exec(lastInsertID, content.Material_id, content.String_order, content.Qty)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(content)
			return -1, err
		}
	}

	return lastInsertID, nil
}

func (s *Storage) WriteRecipePrice(recipe_price models.Recipe_price) (int, error) {
	lastInsertID := -1
	queryText := `INSERT INTO public.recipe_prices (recipe_id,date,price)
	VALUES($1,$2,$3) RETURNING id;`
	err := s.Pdb.QueryRow(queryText, recipe_price.Recipe_id, recipe_price.Date, recipe_price.Price).Scan(&lastInsertID)
	if err != nil {
		s.Log.Error("Error inserting data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		s.Log.Error(recipe_price)
		return -1, err
	}
	return lastInsertID, nil

}

// RecipesQuery Текст запроса для получения списка рецептов
func getRecipesQuery() string {
	return `
	SELECT 
		recipes.id,
		recipes.name,
		recipes.output,
		recipes.unit_id,
		units.short_name AS unit_short_name,
		0 AS price
	FROM 
		public.recipes AS recipes 
			LEFT JOIN public.units AS units 
			ON recipes.unit_id = units.id
	
	ORDER BY 
	recipes.name;`
}

// RecipesWithPricesQuery Текст запроса для получения списка рецептов с ценами
func getRecipesWithPricesQuery() string {
	return `
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
}

// GetRecipeQuerry Получение рецепта по id
func getRecipeQuerry(id int) string {
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

// GetRecipeWithPriceQuerry Получение рецепта c ценой по id
func getRecipeWithPriceQuerry(id int) string {
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

// GetRecipeContentQuery Получение таблицы рецепта по id
func getRecipeContentQuery(id int) string {
	return `
	SELECT 
		recipes_content.id AS recipe_id,
		recipes_content.material_id,
		materials.name AS material_name,
		recipes_content.qty,
		recipes_content.string_order,
		materials.recipe_unit_id AS unit_id,
		units.name AS unit_name,
		units.short_name AS unit_short_name,
		CAST(0 AS NUMERIC(15,4)) AS coefficient,
		CAST(0 AS NUMERIC(15,4)) AS price		
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

// GetRecipeContentWithPricesQuery Получение таблицы рецепта по id
func getRecipeContentWithPricesQuery(id int) string {
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
		ELSE CAST(COALESCE(prices.price,0)/materials.coefficient AS NUMERIC(15,4))
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
