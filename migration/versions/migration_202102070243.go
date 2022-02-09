package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migration2021020770243() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2021020770243",
		Migrate: func(tx *gorm.DB) error {
			type Example struct {
				UUID        string `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
				Description string
			}
			// Create Table
			if err := tx.AutoMigrate(&Example{}); err != nil {
				return err
			}
			// Add Column
			if !tx.Migrator().HasColumn(&Example{}, "Name") {
				type Example struct {
					Name string
				}
				if err := tx.Migrator().AddColumn(&Example{}, "Name"); err != nil {
					return err
				}
			}
			// Alter Column
			if tx.Migrator().HasColumn(&Example{}, "Name") {
				type Example struct {
					Name int
				}
				if err := tx.Migrator().AlterColumn(&Example{}, "Name"); err != nil {
					return err
				}
			}
			// Rename Column
			if tx.Migrator().HasColumn(&Example{}, "Name") && !tx.Migrator().HasColumn(&Example{}, "Age") {
				type Example struct {
					Name int
					Age  int
				}
				if err := tx.Migrator().RenameColumn(&Example{}, "Name", "Age"); err != nil {
					return err
				}
			}
			// Drop Column
			if err := tx.Migrator().DropColumn(&Example{}, "Description"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Example{}, "Age"); err != nil {
				return err
			}

			// Rename Table
			type ExampleNew struct {
			}
			if tx.Migrator().HasTable(&Alert{}) {
				if err := tx.Migrator().RenameTable(&Example{}, &ExampleNew{}); err != nil {
					return err
				}
			}

			// Drop Table
			if err := tx.Migrator().DropTable(&Example{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&ExampleNew{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
