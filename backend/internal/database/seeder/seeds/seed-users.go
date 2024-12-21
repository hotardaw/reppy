package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TestUser struct {
	Email     string
	Password  string
	Username  string
	FirstName string
	LastName  string
	Gender    string
	DOB       time.Time
	Height    int32
	Weight    int32
}

func GetTestUsers() []TestUser {
	return []TestUser{
		{
			Email:     "aaron@test.com",
			Password:  "password123",
			Username:  "aaronhotard",
			FirstName: "Aaron",
			LastName:  "Hotard",
			Gender:    "Male",
			DOB:       time.Date(1999, 11, 19, 0, 0, 0, 0, time.UTC),
			Height:    74,
			Weight:    185,
		},
		{
			Email:     "aiyana@test.com",
			Password:  "password456",
			Username:  "aiyanathomas",
			FirstName: "Aiyana",
			LastName:  "Thomas",
			Gender:    "Female",
			DOB:       time.Date(2000, 7, 5, 0, 0, 0, 0, time.UTC),
			Height:    63,
			Weight:    130,
		},
		{ // User only, no profile, for /user-profiles testing
			Email:    "hunter@test.com",
			Password: "password789",
			Username: "huntertracy",
		},
	}
}

func SeedUsers(queries *sqlc.Queries) error {
	for _, user := range GetTestUsers() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		createdUser, err := queries.CreateUser(context.Background(), sqlc.CreateUserParams{
			Email:        user.Email,
			PasswordHash: string(hashedPassword),
			Username:     user.Username,
		})
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %v", user.Email, err)
		}

		// Only create profile if the user has profile information
		if user.FirstName != "" || user.LastName != "" || !user.DOB.IsZero() || user.Gender != "" || user.Height != 0 || user.Weight != 0 {

			_, err = queries.CreateUserProfile(context.Background(), sqlc.CreateUserProfileParams{
				UserID: sql.NullInt32{
					Int32: createdUser.UserID,
					Valid: true,
				},
				FirstName: sql.NullString{
					String: user.FirstName,
					Valid:  user.FirstName != "",
				},
				LastName: sql.NullString{
					String: user.LastName,
					Valid:  user.LastName != "",
				},
				DateOfBirth: sql.NullTime{
					Time:  user.DOB,
					Valid: !user.DOB.IsZero(),
				},
				Gender: sql.NullString{
					String: user.Gender,
					Valid:  user.Gender != "",
				},
				HeightInches: sql.NullInt32{
					Int32: user.Height,
					Valid: user.Height != 0,
				},
				WeightPounds: sql.NullInt32{
					Int32: user.Weight,
					Valid: user.Weight != 0,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create profile for user %s: %v", user.Email, err)
			}
		}
	}
	fmt.Println("Successfully seeded USERS & USER_PROFILES table")
	return nil
}
