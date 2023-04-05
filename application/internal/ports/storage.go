package ports

import "bakery/application/internal/domain/models"

type Storage interface {
	Start()
	ReadMaterials(prices bool) ([]models.Material, error)
	ReadMaterial(prices bool, id int) (models.Material, error)
	WriteMaterial(data models.Material) (int, error)
	WriteMaterialPrice(data models.Material_price) (int, error)

	ReadOrders() ([]models.Order, error)
	ReadOrder(id int) (models.Order, error)
	WriteOrder(data models.Order) (int, error)

	ReadRecipe(prices bool, id int) (models.Recipe, error)
	ReadRecipes(prices bool) ([]models.Recipe, error)
	WriteRecipe(data models.Recipe) (int, error)
	WriteRecipePrice(data models.Recipe_price) (int, error)

	ReadUnits() ([]models.Unit, error)
	ReadUnit(id int) (models.Unit, error)
	WriteUnit(data models.Unit) (int, error)
}
