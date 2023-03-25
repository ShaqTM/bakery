package ports

type Storage interface {
	Start()
	ReadMaterials(prices bool) ([]map[string]interface{}, error)
	ReadMaterial(prices bool, id int) (map[string]interface{}, error)

	ReadOrders() ([]map[string]interface{}, error)
	ReadOrderContent(id int) ([]map[string]interface{}, error)

	ReadRecipe(prices bool, id int) (map[string]interface{}, error)
	ReadRecipes(prices bool) ([]map[string]interface{}, error)
	ReadRecipeContent(prices bool, id int) ([]map[string]interface{}, error)

	ReadUnits() ([]map[string]interface{}, error)
	ReadUnit(id int) (map[string]interface{}, error)
}
