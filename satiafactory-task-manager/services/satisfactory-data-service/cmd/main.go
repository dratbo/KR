package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/config"
	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/database"
	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/handlers"
	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/parser"
	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/repository"
)

func main() {
	importFlag := flag.Bool("import", false, "import data from Docs.json and exit")
	flag.Parse()

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	database.RunMigrations(db, "migrations")

	if *importFlag {
		log.Println("Starting import...")
		if err := parser.RunParser(db, cfg.DataFilePath); err != nil {
			log.Fatal("Import failed:", err)
		}
		log.Println("Import completed successfully")
		return
	}

	recipeRepo := repository.NewRecipeRepository(db)
	count, err := recipeRepo.Count()
	if err != nil {
		log.Fatal("recipe count:", err)
	}
	if count == 0 {
		log.Println("Database empty, importing Docs.json (first run may take a few minutes)...")
		if err := parser.RunParser(db, cfg.DataFilePath); err != nil {
			log.Fatal("Auto-import failed:", err)
		}
	} else {
		log.Printf("Recipes in database: %d", count)
	}

	// Инициализация репозиториев
	itemRepo := repository.NewItemRepository(db)
	buildingRepo := repository.NewBuildingRepository(db)

	// Хендлеры
	itemHandler := handlers.NewItemHandler(itemRepo)
	recipeHandler := handlers.NewRecipeHandler(recipeRepo)
	buildingHandler := handlers.NewBuildingHandler(buildingRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/items", itemHandler.ListItems)
	mux.HandleFunc("GET /api/items/{className}", itemHandler.GetItem)
	mux.HandleFunc("GET /api/recipes", recipeHandler.ListRecipes)
	mux.HandleFunc("GET /api/recipes/search", recipeHandler.SearchRecipes)
	mux.HandleFunc("GET /api/recipes/{className}", recipeHandler.GetRecipe)
	mux.HandleFunc("GET /api/buildings", buildingHandler.ListBuildings)

	log.Printf("Satisfactory Data Service running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
