package models

import (
    "os"
    "path/filepath"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
)

func Setup() (*gorm.DB, error) {

   //Database retriev
    //dbPath := filepath.Join("/home/azureuser/Amazn/database.db", "database.db")
    dbPath := "database.db"


    //If not create
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        log.Println("Creating the file")
        file, err := os.Create(dbPath)
        if err != nil {
            return nil, err
        }
        file.Close()
    }


    //AutoMigrate
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
    if err = db.AutoMigrate(&OrderItems{}); err != nil {
        log.Println(err)
    }
    if err = db.AutoMigrate(&Order{}); err != nil {
        log.Println(err)
    }
    return db, nil
}

