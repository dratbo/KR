package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations applies *.up.sql files from dir in lexical order.
func RunMigrations(db *sql.DB, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("migrations: read dir %s: %v", dir, err)
		return
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		path := filepath.Join(dir, name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("migrations: read %s: %v", path, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			log.Printf("migrations: %s: %v", name, err)
		} else {
			log.Printf("migrations: applied %s", name)
		}
	}
}
