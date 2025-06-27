package main

import (
	"fmt"
	"log"

	"moon/internal/config"
	"moon/internal/database"
	"moon/internal/domain/user"
	"moon/pkg/hash"
)

func main() {
	// Load configuration
	if err := config.LoadConfig("../configs/config.yaml"); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	cfg := config.GetConfig()

	// Connect to database
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := database.GetDB()

	// Auto migrate
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Check if admin already exists
	var existingAdmin user.User
	result := db.Where("email = ?", "admin@moon.com").First(&existingAdmin)
	if result.Error == nil {
		fmt.Println("Admin user already exists!")
		fmt.Printf("Email: %s\n", existingAdmin.Email)
		fmt.Printf("Name: %s\n", existingAdmin.Name)
		fmt.Printf("Role: %s\n", existingAdmin.Role)
		return
	}

	// Hash password
	hashedPassword, err := hash.HashPassword("admin123")
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create admin user
	admin := user.User{
		Email:    "admin@moon.com",
		Password: hashedPassword,
		Name:     "Administrator",
		Phone:    nil,
		Address:  nil,
		Lat:      nil,
		Lng:      nil,
		Role:     "admin",
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Println("📧 Email: admin@moon.com")
	fmt.Println("🔑 Password: admin123")
	fmt.Println("⚠️  Please change the password after first login!")
}
