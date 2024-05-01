package database

//All User related methods are defined here

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
)

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}

func (db *appdbimpl) GetUser(userID string) (*User, error) {
	if userID == "" {
		return nil, nil
	}
	var user User

	// SQL query to select the user by user_id
	query := "SELECT user_id, username FROM users WHERE user_id = ?"

	// Execute the query
	err := db.c.QueryRow(query, userID).Scan(&user.ID, &user.Username)
	if err != nil {
		// Other error occurred
		return nil, err
	}
	return &user, nil
}

// checkUserIDExists now returns an error as well
func (db *appdbimpl) checkUserIDExists(userID string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)", userID).Scan(&exists)
	return exists, err // return the error
}

// generateUniqueID now has a receiver and returns errors
func (db *appdbimpl) generateUniqueID() (string, error) {
	for {
		userID, err := generateRandomString(10)
		if err != nil {
			return "", err
		}

		exists, err := db.checkUserIDExists(userID) // use the method of appdbimpl
		if err != nil {
			return "", err // return the error
		}

		if !exists {
			return userID, nil
		}
	}
}

// AddUser uses the generateUniqueID method of appdbimpl
func (db *appdbimpl) AddUser(user *User) error {
	userID, err := db.generateUniqueID()
	if err != nil {
		return fmt.Errorf("failed to generate user ID: %w", err)
	}
	user.ID = userID

	stmt, err := db.c.Prepare("INSERT INTO users (user_id, username) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Username)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("username already exists: %w", err)
		}
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (db *appdbimpl) SetUsername(userId, newUsername string) error {
	stmt, err := db.c.Prepare("UPDATE users SET username = ? WHERE user_id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUsername, userId)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (db *appdbimpl) GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.c.QueryRow("SELECT user_id, username FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		return nil, nil // User not found is not an error here
	}
	return &user, err
}

func (db *appdbimpl) GetUserProfile(username string) (*User, error) {
	var user User

	// Fetch basic user info
	err := db.c.QueryRow("SELECT user_id, username FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Fetch followers
	rows, err := db.c.Query("SELECT follower_id FROM followers WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var followerID string
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("failed to read follower row: %w", err)
		}
		user.Followers = append(user.Followers, followerID)
	}

	// Fetch following
	rows, err = db.c.Query("SELECT user_id FROM followers WHERE follower_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch following: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var followingID string
		if err := rows.Scan(&followingID); err != nil {
			return nil, fmt.Errorf("failed to read following row: %w", err)
		}
		user.Following = append(user.Following, followingID)
	}

	// Fetch photos
	rows, err = db.c.Query("SELECT photo_id FROM new_photos WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch photos: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var photoID string
		if err := rows.Scan(&photoID); err != nil {
			return nil, fmt.Errorf("failed to read photo row: %w", err)
		}
		user.Photos = append(user.Photos, photoID)
	}

	return &user, nil
}

func (db *appdbimpl) GetUserProfileByID(userID string) (*User, error) {
	// Fetch basic user info
	var user User
	err := db.c.QueryRow("SELECT user_id, username FROM users WHERE user_id = ?", userID).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	// Fetch followers
	rows, err := db.c.Query("SELECT follower_id FROM followers WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var followerID string
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("failed to read follower row: %w", err)
		}
		user.Followers = append(user.Followers, followerID)
	}

	// Fetch following
	rows, err = db.c.Query("SELECT user_id FROM followers WHERE follower_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch following: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var followingID string
		if err := rows.Scan(&followingID); err != nil {
			return nil, fmt.Errorf("failed to read following row: %w", err)
		}
		user.Following = append(user.Following, followingID)
	}

	// Fetch photos
	rows, err = db.c.Query("SELECT photo_id FROM new_photos WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch photos: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var photoID string
		if err := rows.Scan(&photoID); err != nil {
			return nil, fmt.Errorf("failed to read photo row: %w", err)
		}
		user.Photos = append(user.Photos, photoID)
	}

	return &user, nil
}

func (db *appdbimpl) FollowUser(followerID, followedID string) error {
	_, err := db.c.Exec(`INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`, followedID, followerID)
	if err != nil {
		return fmt.Errorf("error following user: %w", err)
	}
	return nil
}

func (db *appdbimpl) UnfollowUser(followerID, followedID string) error {
	_, err := db.c.Exec(`DELETE FROM followers WHERE user_id = ? AND follower_id = ?`, followedID, followerID)
	if err != nil {
		return fmt.Errorf("error unfollowing user: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetUserIDByUsername(username string) (string, error) {
	var userID string
	err := db.c.QueryRow("SELECT user_id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("query error: %w", err)
	}
	return userID, nil
}

// getAllUsers
func (db *appdbimpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT user_id, username FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *appdbimpl) GetFollowersByUsername(username string) ([]string, error) {
	var followers []string
	query := `SELECT follower_id FROM followers WHERE user_id = (SELECT user_id FROM users WHERE username = ?)`
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("error querying followers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var followerID string
		if err := rows.Scan(&followerID); err != nil {
			return nil, fmt.Errorf("error scanning follower ID: %w", err)
		}
		followers = append(followers, followerID)
	}
	return followers, nil
}

func (db *appdbimpl) GetUsername(userID string) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM users WHERE user_id = ?", userID).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("error getting username: %w", err)
	}
	return username, nil
}

func (db *appdbimpl) IsUserFollowed(followedID, followerID string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE user_id = ? AND follower_id = ?)", followedID, followerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if user is followed: %w", err)
	}
	return exists, nil
}
