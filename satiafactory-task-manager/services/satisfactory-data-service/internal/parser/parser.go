package parser

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/models"
	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/repository"
)

// DocsRoot соответствует структуре вашего файла
type DocsRoot struct {
	GameVersion string     `json:"GameVersion"`
	Classes     []RawClass `json:"Classes"`
}

type RawClass struct {
	NativeClass string                 `json:"NativeClass"`
	ClassName   string                 `json:"Class"` // полное имя класса
	Name        string                 `json:"Name"`  // короткое имя (ID)
	Properties  map[string]interface{} `json:"Properties"`
}

// RunParser читает Docs.json и наполняет БД
func RunParser(db *sql.DB, filePath string) error {
	log.Println("Parsing Docs.json...")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var root DocsRoot
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return err
	}

	itemRepo := repository.NewItemRepository(db)
	buildingRepo := repository.NewBuildingRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	// schematicRepo := repository.NewSchematicRepository(db) // если понадобится

	var pendingItems []*models.Item
	var pendingBuildings []*models.Building
	var pendingRecipes []*models.Recipe

	for _, c := range root.Classes {
		switch c.NativeClass {
		case "FGItemDescriptor", "FGResourceDescriptor", "FGEquipmentDescriptor", "FGConsumableDescriptor":
			if item := parseItem(c); item != nil {
				pendingItems = append(pendingItems, item)
			}
		case "FGBuildable", "FGBuildableFactory", "FGBuildableResourceExtractor", "FGBuildableManufacturer", "FGBuildableGenerator", "FGBuildableConstructor", "FGBuildableMiner", "FGBuildablePipeline", "FGBuildableTrainPlatform", "FGBuildableVehicle":
			if building := parseBuilding(c); building != nil {
				pendingBuildings = append(pendingBuildings, building)
			}
		case "FGRecipe":
			if recipe := parseRecipe(c); recipe != nil {
				pendingRecipes = append(pendingRecipes, recipe)
			}
		}
	}

	itemCount := 0
	for _, item := range pendingItems {
		if err := itemRepo.Insert(item); err != nil {
			log.Printf("Warning: insert item %s: %v", item.ClassName, err)
		} else {
			itemCount++
		}
	}

	buildingCount := 0
	for _, building := range pendingBuildings {
		if err := buildingRepo.Insert(building); err != nil {
			log.Printf("Warning: insert building %s: %v", building.ClassName, err)
		} else {
			buildingCount++
		}
	}

	recipeCount := 0
	for _, recipe := range pendingRecipes {
		if err := recipeRepo.Insert(recipe); err != nil {
			log.Printf("Warning: insert recipe %s: %v", recipe.ClassName, err)
		} else {
			recipeCount++
		}
	}

	log.Printf("Inserted/updated %d items", itemCount)
	log.Printf("Inserted/updated %d buildings", buildingCount)
	log.Printf("Inserted/updated %d recipes", recipeCount)
	log.Println("Import completed successfully")
	return nil
}

// Вспомогательные функции для извлечения строк из Properties
func getStringProp(props map[string]interface{}, key string) string {
	if v, ok := props[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getFloatProp(props map[string]interface{}, key string) float64 {
	if v, ok := props[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		}
	}
	return 0
}

func getIntProp(props map[string]interface{}, key string) int {
	if v, ok := props[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return 0
}

func getStringSliceProp(props map[string]interface{}, key string) []string {
	if v, ok := props[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			res := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					res = append(res, s)
				} else if m, ok := item.(map[string]interface{}); ok {
					// иногда элементы могут быть объектами с полем "ItemClass" и т.п.
					if className, ok := m["Class"].(string); ok {
						res = append(res, extractShortName(className))
					}
				}
			}
			return res
		}
	}
	return nil
}

// extractShortName из полного имени класса "BlueprintGeneratedClass'FactoryGame/.../Desc_IronIngot_C'"
// возвращает "Desc_IronIngot_C"
func extractShortName(full string) string {
	// Ищем последний сегмент после '/'
	idx := strings.LastIndex(full, "/")
	if idx != -1 && idx+1 < len(full) {
		partial := full[idx+1:]
		// Убираем возможную кавычку в конце
		partial = strings.TrimSuffix(partial, "'")
		return partial
	}
	return full
}

func parseItem(c RawClass) *models.Item {
	props := c.Properties
	displayName := getStringProp(props, "mDisplayName")
	if displayName == "" {
		// некоторые предметы не имеют отображаемого имени — пропускаем
		return nil
	}
	description := getStringProp(props, "mDescription")
	stackSize := getIntProp(props, "mStackSize")
	energyValue := getFloatProp(props, "mEnergyValue")
	form := getStringProp(props, "mForm")
	// иконки обычно лежат в mSmallIcon, mBigIcon, но они могут быть сложными объектами
	smallIcon := ""
	bigIcon := ""
	categories := []string{}

	return &models.Item{
		ClassName:   c.Name, // короткое имя используем как primary key
		DisplayName: displayName,
		Description: description,
		StackSize:   stackSize,
		EnergyValue: energyValue,
		Form:        form,
		SmallIcon:   smallIcon,
		BigIcon:     bigIcon,
		Categories:  categories,
	}
}

func parseBuilding(c RawClass) *models.Building {
	props := c.Properties
	displayName := getStringProp(props, "mDisplayName")
	if displayName == "" {
		return nil
	}
	description := getStringProp(props, "mDescription")
	powerConsumption := getFloatProp(props, "mPowerConsumption")
	powerConsumptionExponent := getFloatProp(props, "mPowerConsumptionExponent")
	return &models.Building{
		ClassName:                c.Name,
		DisplayName:              displayName,
		Description:              description,
		PowerConsumption:         powerConsumption,
		PowerConsumptionExponent: powerConsumptionExponent,
	}
}

func parseRecipe(c RawClass) *models.Recipe {
	props := c.Properties
	displayName := getStringProp(props, "mDisplayName")
	if displayName == "" {
		return nil
	}
	duration := getFloatProp(props, "mManufactoringDuration")
	if duration == 0 {
		duration = getFloatProp(props, "mProductionDuration")
	}
	menuPriority := getIntProp(props, "mManufactoringMenuPriority")

	// Извлекаем mIngredients
	var ingredients []models.Ingredient
	if ingArr, ok := props["mIngredients"].([]interface{}); ok {
		for _, ing := range ingArr {
			if ingMap, ok := ing.(map[string]interface{}); ok {
				if itemClass, ok := ingMap["ItemClass"].(map[string]interface{}); ok {
					if className, ok := itemClass["Name"].(string); ok {
						amount := getFloatProp(ingMap, "Amount")
						ingredients = append(ingredients, models.Ingredient{
							ItemClassName: className,
							Amount:        amount,
						})
					}
				}
			}
		}
	}

	// Извлекаем mProduct (может быть массив или одиночный объект)
	var products []models.Product
	if prodArr, ok := props["mProduct"].([]interface{}); ok {
		for _, prod := range prodArr {
			if prodMap, ok := prod.(map[string]interface{}); ok {
				if itemClass, ok := prodMap["ItemClass"].(map[string]interface{}); ok {
					if className, ok := itemClass["Name"].(string); ok {
						amount := getFloatProp(prodMap, "Amount")
						products = append(products, models.Product{
							ItemClassName: className,
							Amount:        amount,
						})
					}
				}
			}
		}
	} else if prodObj, ok := props["mProduct"].(map[string]interface{}); ok {
		if itemClass, ok := prodObj["ItemClass"].(map[string]interface{}); ok {
			if className, ok := itemClass["Name"].(string); ok {
				amount := getFloatProp(prodObj, "Amount")
				products = append(products, models.Product{
					ItemClassName: className,
					Amount:        amount,
				})
			}
		}
	}

	// Извлекаем mProducedIn (массив)
	var producedIn []string
	if prodInArr, ok := props["mProducedIn"].([]interface{}); ok {
		for _, item := range prodInArr {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if className, ok := itemMap["Name"].(string); ok {
					producedIn = append(producedIn, className)
				}
			}
		}
	}

	return &models.Recipe{
		ClassName:                 c.Name,
		DisplayName:               displayName,
		Ingredients:               ingredients,
		Products:                  products,
		ProducedIn:                producedIn,
		Duration:                  duration,
		ManufactoringMenuPriority: menuPriority,
	}
}
