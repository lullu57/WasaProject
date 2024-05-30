package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) BanUser(bannedBy, bannedUser string) error {

	exists, err := db.BanExists(bannedBy, bannedUser)
	if err != nil {
		return fmt.Errorf("error checking if ban exists: %w", err)
	}
	if exists {
		return fmt.Errorf("user is already banned")
	}
	stmt, err := db.c.Prepare("INSERT INTO new_bans (ban_id,banned_by, banned_user, timestamp) VALUES (?,?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare ban statement: %w", err)
	}
	defer stmt.Close()
	// generate a unique ban id
	banId, err := generateRandomString(10)
	if err != nil {
		return fmt.Errorf("failed to generate ban id: %w", err)
	}
	_, err = stmt.Exec(banId, bannedBy, bannedUser, time.Now())
	if err != nil {
		return fmt.Errorf("failed to execute ban statement: %w", err)
	}

	return nil
}

func (db *appdbimpl) UnbanUser(bannerID, bannedUserID string) error {
	stmt, err := db.c.Prepare("DELETE FROM new_bans WHERE banned_by = ? AND banned_user = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare unban statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(bannerID, bannedUserID)
	if err != nil {
		return fmt.Errorf("failed to execute unban statement: %w", err)
	}

	return nil
}

func (db *appdbimpl) BanExists(bannedBy, bannedUser string) (bool, error) {
	var exists bool
	stmt, err := db.c.Prepare("SELECT EXISTS(SELECT 1 FROM new_bans WHERE banned_by = ? AND banned_user = ?)")
	if err != nil {
		return false, fmt.Errorf("failed to prepare check ban existence statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(bannedBy, bannedUser).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute check ban existence statement: %w", err)
	}

	return exists, nil
}
