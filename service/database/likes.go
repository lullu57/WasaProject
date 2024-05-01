package database

import (
	"fmt"
)

func (db *appdbimpl) LikePhoto(userID string, photoID string) error {
	// Check if the like already exists to avoid duplicates
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND photo_id = ?)", userID, photoID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	if exists {
		return fmt.Errorf("like already exists")
	}

	// Insert the like into the database
	_, err = db.c.Exec("INSERT INTO likes (user_id, photo_id, timestamp) VALUES (?, ?, CURRENT_TIMESTAMP)", userID, photoID)
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}
	return nil
}

func (db *appdbimpl) UnlikePhoto(userID string, photoID string) error {
	// Delete the like from the database
	_, err := db.c.Exec("DELETE FROM likes WHERE user_id = ? AND photo_id = ?", userID, photoID)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}
	return nil
}

func (db *appdbimpl) IsLiked(photoID string, userID string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND photo_id = ?)", userID, photoID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query error: %w", err)
	}
	return exists, nil
}
