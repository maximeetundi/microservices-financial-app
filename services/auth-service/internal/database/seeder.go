package database

import (
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
)

// SeedDefaultUser checks if the default test user exists and creates it if not.
func SeedDefaultUser(userRepo *repository.UserRepository) {
	email := "maximeetundi@gmail.com"

	exists, err := userRepo.EmailExists(email)
	if err != nil {
		log.Printf("[Seeder] Error checking for default user: %v", err)
		return
	}

	if exists {
		log.Printf("[Seeder] Default user %s already exists", email)
		return
	}

	log.Printf("[Seeder] creating default user %s...", email)

	// User details provided by request
	// Phone: +237 698915622
	// Password: 6Lj[,]SaNnX3}W26
	req := &models.RegisterRequest{
		Email:       email,
		Phone:       "+237698915622", // Format as E.164 if possible, or let repo handle it (but repo might expect validated input if called directly? No, repo.Create expects raw strings, AuthService.Register does val. I'm calling repo directly for seed, so I should provide valid data).
		Password:    "6Lj[,]SaNnX3}W26",
		FirstName:   "Maxime",
		LastName:    "ETUNDI",
		Country:     "CM",
		Currency:    "XAF",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	user, err := userRepo.Create(req)
	if err != nil {
		log.Printf("[Seeder] Failed to create default user: %v", err)
		return
	}

	log.Printf("[Seeder] Successfully created default user: %s (ID: %s)", user.Email, user.ID)
}
