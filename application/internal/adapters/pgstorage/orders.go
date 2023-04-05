package pgstorage

import (
	"bakery/application/internal/domain/models"
	"strconv"
)

// ReadOrders Читает список заказов
func (s *Storage) ReadOrders() ([]models.Order, error) {

	queryText := getOrdersQuery()
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return nil, err
	}

	orders := make([]models.Order, 0)
	price := ""
	plan_qty := ""
	plan_cost := ""
	fact_qty := ""
	fact_cost := ""
	materials_cost := ""
	for rows.Next() {
		order := models.Order{}
		if err := rows.Scan(&order.Id, &order.Customer, &order.Recipe_id,
			&order.Recipe_name, &order.Date, &order.Release_date,
			&price, &plan_qty, &plan_cost, &fact_qty, &fact_cost, &materials_cost,
			&order.Unit_id, &order.Unit_short_name); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return nil, err
		}
		order.Price = convertNumeric(price)
		order.Plan_qty = convertNumeric(plan_qty)
		order.Plan_cost = convertNumeric(plan_cost)
		order.Fact_cost = convertNumeric(fact_cost)
		order.Fact_qty = convertNumeric(fact_qty)
		order.Materials_cost = convertNumeric(materials_cost)

		orders = append(orders, order)
	}
	return orders, nil
}

// ReadOrder Читает заказ по id
func (s *Storage) ReadOrder(id int) (models.Order, error) {
	queryText := getOrderQuerry(id)
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return models.Order{}, err
	}
	if rows.Next() {
		order := models.Order{}
		price := ""
		plan_qty := ""
		plan_cost := ""
		fact_qty := ""
		fact_cost := ""
		materials_cost := ""
		if err := rows.Scan(&order.Id, &order.Customer, &order.Recipe_id,
			&order.Recipe_name, &order.Date, &order.Release_date,
			&price, &plan_qty, &plan_cost, &fact_qty, &fact_cost, &materials_cost,
			&order.Unit_id, &order.Unit_short_name); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return models.Order{}, err
		}
		order.Price = convertNumeric(price)
		order.Plan_qty = convertNumeric(plan_qty)
		order.Plan_cost = convertNumeric(plan_cost)
		order.Fact_cost = convertNumeric(fact_cost)
		order.Fact_qty = convertNumeric(fact_qty)
		order.Materials_cost = convertNumeric(materials_cost)

		content, err := s.readOrderContent(id)
		if err != nil {
			return models.Order{}, err
		}
		order.Content = content
		return order, nil
	}
	return models.Order{}, nil
}

// ReadOrderContent Читает табличную часть заказа
func (s *Storage) readOrderContent(id int) ([]models.Order_detail, error) {

	queryText := getOrderContentQuerry(id)
	rows, err := s.Pdb.Query(queryText)
	if err != nil {
		s.Log.Error("Error reading data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return nil, err
	}

	contents := make([]models.Order_detail, 0)

	qty := ""
	price := ""
	cost := ""

	for rows.Next() {
		content := models.Order_detail{}
		if err := rows.Scan(&content.Id, &content.Material_id, &content.Material_name,
			&qty, &content.String_order, &price, &cost,
			&content.Unit_name, &content.Unit_short_name, &content.By_recipe); err != nil {
			s.Log.Error("Error scanning rows:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			return nil, err
		}
		content.Qty = convertNumeric(qty)
		content.Price = convertNumeric(price)
		content.Cost = convertNumeric(cost)
		contents = append(contents, content)
	}
	return contents, nil
}

func (s *Storage) WriteOrder(order models.Order) (int, error) {
	newRecord := order.Id == -1
	lastInsertID := -1

	if newRecord {
		queryText := `INSERT INTO public.orders (id,customer,recipe_id,
			date,release_date,unit_id,price,plan_qty,plan_cost,
			fact_qty,fact_cost,materials_cost)
									 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id;`
		err := s.Pdb.QueryRow(queryText, order.Id, order.Customer, order.Recipe_id, order.Date,
			order.Release_date, order.Unit_id, order.Price, order.Plan_qty, order.Plan_cost,
			order.Fact_qty, order.Fact_cost, order.Materials_cost).Scan(&lastInsertID)
		if err != nil {
			s.Log.Error("Error inserting data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(order)
			return -1, err
		}
	} else {
		queryText := `UPDATE public.orders SET (id,customer,recipe_id,
			date,release_date,unit_id,price,plan_qty,plan_cost,
			fact_qty,fact_cost,materials_cost) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) WHERE id = ` + strconv.Itoa(order.Id) + `;`
		_, err := s.Pdb.Exec(queryText, order.Id, order.Customer, order.Recipe_id, order.Date,
			order.Release_date, order.Unit_id, order.Price, order.Plan_qty, order.Plan_cost,
			order.Fact_qty, order.Fact_cost, order.Materials_cost)
		if err != nil {
			s.Log.Error("Error updating data:", err)
			s.Log.Error("Query:")
			s.Log.Error(queryText)
			s.Log.Error(order)
			return -1, err
		}
		lastInsertID = order.Id
	}

	queryText := `DELETE FROM public.order_details WHERE id = ` + strconv.Itoa(lastInsertID) + `;`
	_, err := s.Pdb.Exec(queryText)
	if err != nil {
		s.Log.Error("Error deleting old rows:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return -1, err
	}
	if len(order.Content) == 0 {
		return lastInsertID, nil
	}

	queryText = `INSERT INTO public.order_details (id,material_id,
		qty,string_order,price,cost,by_recipe) VALUES($1,$2,$3,$4,$5,$6,$7);`
	stmt, err := s.Pdb.Prepare(queryText)
	if err != nil {
		s.Log.Error("Error preparing for inserting data:", err)
		s.Log.Error("Query:")
		s.Log.Error(queryText)
		return -1, nil
	}
	for _, content := range order.Content {
		_, err := stmt.Exec(lastInsertID, content.Material_id, content.Qty, content.String_order, content.Price, content.Cost, content.By_recipe)
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

// OrdersQuery Текст запроса для получения списка заказов
func getOrdersQuery() string {
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

ORDER BY 
	orders.id;`
}

// GetOrderQuerry Получение заказа по id
func getOrderQuerry(id int) string {
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

// GetOrderContentQuerry Получение таблицы заказа по id
func getOrderContentQuerry(id int) string {
	return `
	SELECT 
		order_details.id AS order_id,
		order_details.material_id,
		materials.name AS material_name,
		order_details.qty,
		order_details.string_order,
		order_details.price,
		order_details.cost,
		units.name AS unit_name,
		units.short_name AS unit_short_name,
		order_details.by_recipe
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
