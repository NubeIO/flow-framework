# Gormigrate [For more details](https://pkg.go.dev/github.com/go-gormigrate/gormigrate/v2#section-documentation)

## Usage (Sqlite)

```go
package main

import (
    "github.com/go-gormigrate/gormigrate/v2",
    "gorm.io/driver/sqlite"
)

func main() {
    db, err := gorm.Open(sqlite.Open("example.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Create migrations
    m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
        ID: "", // Migration identifier
        Migrate: func(tx *gorm.DB) error { // Executes migration		
        },
        Rollback: func(tx *gorm.DB) error { // Executes rollback
        },
    },
    // Executes all migrations
    if err = m.Migrate();
}
```

## Example
 
```
m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202101301400",
			Migrate: func(tx *gorm.DB) error {
				type Alert struct {
					Description string
				}
				return tx.Migrator().AddColumn(&Alert{}, "Description")
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "202102301400",
			Migrate: func(tx *gorm.DB) error {
				query := "ALTER TABLE messages ADD COLUMN status numeric;"
				return tx.Exec(query).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})
	m.InitSchema(func(tx *gorm.DB) error {
		type Alert struct {
			ID       uint   `gorm:"AUTO_INCREMENT;primary_key;index"`
			Type     string
			Duration int
		}
		type Message struct {
			ID       uint   `gorm:"AUTO_INCREMENT;primary_key;index"`
			Message  string `gorm:"type:text"`
			Title    string `gorm:"type:text"`
			Priority int
			Extras   []byte
			Date     time.Time
		}
		return tx.AutoMigrate(
			&Alert{},
			&Message{},
		)
	})
```

### Table operation
```bash
# Add
tx.AutoMigrate(&model.Alert{})

# Drop
tx.Migrator().DropTable(&[]model.Alert{})

#Rename
type AlertNew struct {
    Description string
}
if tx.Migrator().HasTable(&Alert{}) {
    if err := tx.Migrator().RenameTable(&modelAlert{}, &AlertNew{}); err != nil {
        return err
    }
}
```

### Column operation
```bash
type Alert struct {
    Description string
}

# Add
tx.Migrator().AddColumn(&Alert{}, "Description")

# Alter
type Alert struct {
    Description int
}
tx.Migrator().AlterColumn(&Alert{}, "Description")

# Rename
type Alert struct {
    Description  string
    Description1 string
}
tx.Migrator().RenameColumn(&Alert{}, "Description", "Description1")

# Drop
tx.Migrator().DropColumn(&Alert{}, "Description1")
```