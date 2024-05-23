package database

import "fmt"

func (db *appdbimpl) AddComment(comment Comment) error {
	stmt, err := db.c.Prepare("INSERT INTO comments (comment_id, user_id, photo_id, content, timestamp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.ID, comment.UserID, comment.PhotoID, comment.Content, comment.Timestamp)
	return err
}

func (db *appdbimpl) DeleteComment(commentID string) error {
	stmt, err := db.c.Prepare("DELETE FROM comments WHERE comment_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(commentID)
	return err
}
func (db *appdbimpl) GetCommentsByPhotoId(photoId string) ([]Comment, error) {
	// SQL query to fetch all comments for a given photo ID
	query := `SELECT comment_id, user_id, photo_id, content, timestamp FROM comments WHERE photo_id = ? ORDER BY timestamp DESC`
	rows, err := db.c.Query(query, photoId)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.PhotoID, &c.Content, &c.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return comments, nil
}
