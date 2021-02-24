package store

import "strconv"

//ReadOrderList Читает список заказов
func (mdb MDB) ReadOrderList() ([]map[string]interface{}, error) {
	return mdb.ReadRows(OrdersQuery)
}

//ReadOrder Читает заказ по id
func (mdb MDB) ReadOrder(id int) (map[string]interface{}, error) {
	data, err := mdb.ReadRow(GetOrderQuerry(id))
	if err != nil {
		return data, err
	}
	content, err := mdb.ReadOrderContent(id)
	if err != nil {
		return data, err
	}
	data["content"] = content
	return data, nil
}

//ReadOrderContent Читает табличную часть заказа
func (mdb MDB) ReadOrderContent(id int) ([]map[string]interface{}, error) {
	return mdb.ReadRows(GetOrderContentQuerry(id))
}

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
	orders.materials_cost,
	orders.unit_id,
	units.short_name AS unit_short_name
FROM 
	public.orders AS orders 
		LEFT JOIN public.recipes AS recipes 
		ON orders.recipe_id = recipes.id
		LEFT JOIN public.units AS units 
		ON orders.unit_id = units.id

ORDER BY 
	orders.id;`

//GetOrderQuerry Получение заказа по id
func GetOrderQuerry(id int) string {
	return `
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
		orders.materials_cost,
		orders.unit_id,
		units.short_name AS unit_short_name
	FROM 
		public.orders AS orders
			LEFT JOIN public.recipes AS recipes 
			ON orders.recipe_id = recipes.id 
			LEFT JOIN public.units AS units 
			ON orders.unit_id = units.id			
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
		units.short_name AS unit_short_name,
		by_recipe
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
