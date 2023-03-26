package models

type Unit struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Short_name string `json:"short_name"`
}

type Material struct {
	Id                     int     `json:"id"`
	Name                   string  `json:"name"`
	Recipe_unit_id         int     `json:"recipe_unit_id"`
	Recipe_unit_name       string  `json:"recipe_unit_name"`
	Recipe_unit_short_name string  `json:"recipe_unit_short_name"`
	Price_unit_id          int     `json:"price_unit_id"`
	Price_unit_name        string  `json:"price_unit_name"`
	Price_unit_short_name  string  `json:"price_unit_short_name"`
	Coefficient            float64 `json:"coefficient"`
	Price                  float64 `json:"price"`
}

type Material_price struct {
	Id          int     `json:"id"`
	Material_id int     `json:"material_id"`
	Date        int64   `json:"date"`
	Price       float64 `json:"price"`
}
type Recipe struct {
	Id              int              `json:"id"`
	Name            string           `json:"name"`
	Output          float64          `json:"output"`
	Unit_id         int              `json:"unit_id"`
	Unit_short_name string           `json:"unit_short_name"`
	Price           float64          `json:"price"`
	Content         []Recipe_content `json:"Content"`
}
type Recipe_content struct {
	Id              int     `json:"id"`
	Material_id     int     `json:"material_id"`
	String_order    int     `json:"string_order"`
	Material_name   string  `json:"material_name"`
	Unit_id         int     `json:"unit_id"`
	Unit_name       string  `json:"unit_name"`
	Unit_short_name string  `json:"unit_short_name"`
	Coefficient     float64 `json:"coefficient"`
	Price           float64 `json:"price"`
	Qty             float64 `json:"qty"`
}
type Recipe_price struct {
	Id        int     `json:"id"`
	Recipe_id int     `json:"recipe_id"`
	Date      int64   `json:"date"`
	Price     float64 `json:"price"`
}
type Order struct {
	Id              int            `json:"id"`
	Customer        string         `json:"customer"`
	Recipe_id       int            `json:"recipe_id"`
	Date            int64          `json:"date"`
	Release_date    int64          `json:"release_date"`
	Unit_id         int            `json:"unit_id"`
	Price           float64        `json:"price"`
	Plan_qty        float64        `json:"plan_qty"`
	Plan_cost       float64        `json:"plan_cost"`
	Fact_qty        float64        `json:"fact_qty"`
	Fact_cost       float64        `json:"fact_cost"`
	Materials_cost  float64        `json:"materials_cost"`
	Recipe_name     string         `json:"recipe_name"`
	Unit_short_name string         `json:"unit_short_name"`
	Content         []Order_detail `json:"content"`
}
type Order_detail struct {
	Id              int     `json:"id"`
	Material_id     int     `json:"material_id"`
	Qty             float64 `json:"qty"`
	String_order    int     `json:"string_order"`
	Price           float64 `json:"price"`
	Cost            float64 `json:"cost"`
	By_recipe       bool    `json:"by_recipe"`
	Material_name   string  `json:"material_name"`
	Unit_name       string  `json:"unit_name"`
	Unit_short_name string  `json:"unit_short_name"`
}
