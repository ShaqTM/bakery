package dbservice

import "strconv"

//OrdersQuery Текст запроса для получения списка заказов
const OrdersQuery = `
SELECT 
	orders.id,
	orders.customer,
	orders.recipe_id,
	recipes.name AS recipe_name,
	orders.date,
	orders.release_date,
	orders.price,
	orders.plan_qty,
	orders.plan_cost,
	orders.fact_qty,
	orders.fact_cost,
	orders.materials_cost
FROM 
	public.orders AS orders 
		LEFT JOIN public.recipes AS recipes 
		ON orders.recipe_id = recipes.id
ORDER BY 
	orders.id;`

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

//GetOrderQuerry Получение заказа по id
func GetOrderQuerry(id int) string {
	return `
	SELECT
		orders.id,
		orders.customer,
		orders.recipe_id,
		recipes.name,
		orders.date,
		orders.release_date,
		orders.price,
		orders.plan_qty,
		orders.plan_cost,
		orders.fact_qty,
		orders.fact_cost,
		orders.materials_cost
	FROM 
		public.orders AS orders
			LEFT JOIN public.recipes AS recipes 
			ON orders.recipe_id = recipes.id 
	WHERE 
		orders.id = ` + strconv.Itoa(id) + `;`
}

//GetOrderContentQuerry Получение таблицы заказа по id
func GetOrderContentQuerry(id int) string {
	return `
	SELECT 
		order_details.id,
		order_details.material_id,
		materials.name AS material_name,
		order_details.qty,
		order_details.string_order,
		order_details.price,
		order_details.cost,
		units.name AS unit_name,
		units.short_name AS unit_short_name
	FROM  
		public.order_details AS order_details 
	 		LEFT JOIN public.materials AS materials 
	 		ON order_details.material_id = materials.id
				LEFT JOIN public.units AS units 
				ON units.id  = materials.recipe_unit_id 
	WHERE
		order_details.id = ` + strconv.Itoa(id) + `
	ORDER BY 
		order_details.string_order;`
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
		materials.recipe_unit_id,
		COALESCE(units1.name,'') AS recipe_unit_name,
		COALESCE(units1.short_name,'') AS recipe_unit_short_name,
		materials.price_unit_id,
		COALESCE(units2.name,'') AS price_unit_name,
		COALESCE(units2.short_name,'') AS price_unit_short_name,
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
		materials.recipe_unit_id,
		COALESCE(units1.name,'') AS recipe_unit_name,
		COALESCE(units1.short_name,'') AS recipe_unit_short_name,
		materials.price_unit_id,
		COALESCE(units2.name,'') AS price_unit_name,
		COALESCE(units2.short_name,'') AS price_unit_short_name,
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
		materials.recipe_unit_id,
		COALESCE(units1.name,'') AS recipe_unit_name,
		COALESCE(units1.short_name,'') AS recipe_unit_short_name,
		materials.price_unit_id,
		COALESCE(units2.name,'') AS price_unit_name,
		COALESCE(units2.short_name,'') AS price_unit_short_name,
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
		materials.recipe_unit_id,
		COALESCE(units1.name,'') AS recipe_unit_name,
		COALESCE(units1.short_name,'') AS recipe_unit_short_name,
		materials.price_unit_id,
		COALESCE(units2.name,'') AS price_unit_name,
		COALESCE(units2.short_name,'') AS price_unit_short_name,
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
