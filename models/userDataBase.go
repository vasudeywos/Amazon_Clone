package models

import (
    "os"
    "path/filepath"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
)

func Setup() (*gorm.DB, error) {

    // Define the database
    dbPath := filepath.Join("D:/Amazn/", "database.db")

    // Check if the database file exists, and if not, create it
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        log.Println("Creating the file")
        file, err := os.Create(dbPath)
        if err != nil {
            return nil, err
        }
        file.Close()
    }


    // AutoMigrate your database tables
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

    if err != nil {
        log.Fatal(err.Error())
    }
    if err = db.AutoMigrate(&User{}); err != nil {
        log.Println(err)
    }

    if err = db.AutoMigrate(&Product{}); err != nil {
        log.Println(err)
    }
    return db, nil
}

