package bakery

import (
	"bakery/application/internal/domain/models"
	"bakery/application/internal/ports"

	"github.com/sirupsen/logrus"
)

type Bakery struct {
	Storage ports.Storage
	Log     *logrus.Logger
}

func New(log *logrus.Logger, storage ports.Storage) *Bakery {
	b := &Bakery{
		Storage: storage,
		Log:     log,
	}
	return b

}

func (b *Bakery) ReadMaterials(prices bool) ([]models.Material, error) {
	return b.Storage.ReadMaterials(prices)
}
func (b *Bakery) ReadMaterial(prices bool, id int) (models.Material, error) {
	return b.Storage.ReadMaterial(prices, id)
}
func (b *Bakery) WriteMaterial(data models.Material) (int, error) {
	return b.Storage.WriteMaterial(data)
}

func (b *Bakery) WriteMaterialPrice(data models.Material_price) (int, error) {
	return b.Storage.WriteMaterialPrice(data)
}

func (b *Bakery) ReadOrders() ([]models.Order, error) {
	return b.Storage.ReadOrders()
}
func (b *Bakery) ReadOrder(id int) (models.Order, error) {
	return b.Storage.ReadOrder(id)
}
func (b *Bakery) WriteOrder(data models.Order) (int, error) {
	return b.Storage.WriteOrder(data)
}

func (b *Bakery) ReadRecipe(prices bool, id int) (models.Recipe, error) {
	return b.Storage.ReadRecipe(prices, id)
}
func (b *Bakery) ReadRecipes(prices bool) ([]models.Recipe, error) {
	return b.Storage.ReadRecipes(prices)
}
func (b *Bakery) WriteRecipe(data models.Recipe) (int, error) {
	return b.Storage.WriteRecipe(data)
}

func (b *Bakery) WriteRecipePrice(data models.Recipe_price) (int, error) {
	return b.Storage.WriteRecipePrice(data)
}

func (b *Bakery) ReadUnits() ([]models.Unit, error) {
	return b.Storage.ReadUnits()
}

func (b *Bakery) ReadUnit(id int) (models.Unit, error) {
	return b.Storage.ReadUnit(id)
}
func (b *Bakery) WriteUnit(data models.Unit) (int, error) {
	return b.Storage.WriteUnit(data)
}
