/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"time"
)

type Error struct {
	Error string `json:"error" db:"error"`
}

type User struct {
	ID        string   `json:"userId" db:"user_id"` // Unique identifier
	Username  string   `json:"username" db:"username"`
	Followers []string `json:"followers"` // IDs of followers (handled separately in relational mapping)
	Following []string `json:"following"` // IDs of users being followed (handled separately in relational mapping)
	Photos    []string `json:"photos"`    // IDs of photos uploaded by the user (handled separately in relational mapping)
}

// New Struct for handling followers relationship
type Follower struct {
	UserID     string `json:"userId" db:"user_id"`
	FollowerID string `json:"followerId" db:"follower_id"`
}

// New Struct for handling user and photo relationship
type UserPhoto struct {
	UserID  string `json:"userId" db:"user_id"`
	PhotoID string `json:"photoId" db:"photo_id"`
}

type Comment struct {
	ID        string    `json:"commentId" db:"comment_id"` // Unique identifier
	UserID    string    `json:"userId" db:"user_id"`       // ID of the user who commented
	PhotoID   string    `json:"photoId" db:"photo_id"`     // ID of the photo being commented on
	Content   string    `json:"content" db:"content"`      // The comment itself
	Timestamp time.Time `json:"timestamp" db:"timestamp"`  // Timestamp of when the comment was made
}

type Like struct {
	UserID    string    `json:"userId" db:"user_id"`      // ID of the user who liked the photo
	PhotoID   string    `json:"photoId" db:"photo_id"`    // ID of the photo being liked
	Timestamp time.Time `json:"timestamp" db:"timestamp"` // Timestamp of when the like was made
}

type Photo struct {
	ID        string    `json:"photoId" db:"photo_id"`     // Unique identifier
	UserID    string    `json:"userId" db:"user_id"`       // ID of the user who uploaded the photo
	ImageData []byte    `json:"imageData" db:"image_data"` // The photo data itself
	Timestamp time.Time `json:"timestamp" db:"timestamp"`  // Timestamp of when the photo was uploaded
	Likes     []Like    `json:"likes"`                     // Note: This requires a relational mapping and isn't directly mapped to a single column
	Comments  []Comment `json:"comments"`                  // Note: This requires a relational mapping and isn't directly mapped to a single column
}

type PhotoDetail struct {
	PhotoID    string    `json:"photoId"`
	UserID     string    `json:"userId"`
	Username   string    `json:"username"`
	ImageData  []byte    `json:"imageData"`
	Timestamp  time.Time `json:"timestamp"`
	LikesCount int       `json:"likesCount"`
	Comments   []Comment `json:"comments"`
}
type Ban struct {
	ID         string    `json:"banId" db:"ban_id"`           // Unique identifier
	BannedBy   string    `json:"bannedBy" db:"banned_by"`     // ID of the user who banned the other user
	BannedUser string    `json:"bannedUser" db:"banned_user"` // ID of the user who was banned
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`    // Timestamp of when the ban was made
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	AddUser(user *User) error
	Ping() error
	SetUsername(userId, newUsername string) error
	GetUserProfile(username string) (*User, error)
	LikePhoto(userID string, photoID string) error
	UnlikePhoto(userID string, photoID string) error
	FollowUser(followerID string, followedID string) error
	UnfollowUser(followerID string, followedID string) error
	GetUserIDByUsername(username string) (string, error)
	GetUserByUsername(username string) (*User, error)
	GetUser(userID string) (*User, error)
	AddPhoto(photo Photo) error
	GetPhotos() ([]Photo, error)
	BanUser(bannedBy string, bannedUser string) error
	UnbanUser(bannerID, bannedUserID string) error
	GetBans() ([]Ban, error)
	GetAllUsers(currentUserID string) ([]User, error)
	GetMyStream(userID string) ([]string, error)
	DeleteComment(commentID string) error
	AddComment(comment Comment) error
	DeletePhoto(photoID string) error
	GetCommentsByPhotoId(photoId string) ([]Comment, error)
	GetFollowersByUsername(username string) ([]string, error)
	GetUserProfileByID(userID string) (*User, error)
	GetPhoto(photoId string) (*PhotoDetail, error)
	GetUsername(userID string) (string, error)
	IsLiked(photoID string, userID string) (bool, error)
	IsUserFollowed(followerID, followedID string) (bool, error)
	BanExists(bannedBy, bannedUser string) (bool, error)
}
type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Error table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS errors (
        error TEXT
    );`)
	if err != nil {
		return nil, err
	}

	// Users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        user_id TEXT PRIMARY KEY,
        username TEXT UNIQUE NOT NULL
    );`)
	if err != nil {
		return nil, err
	}

	// Followers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS followers (
        user_id TEXT NOT NULL,
        follower_id TEXT NOT NULL,
        PRIMARY KEY (user_id, follower_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (follower_id) REFERENCES users(user_id)
    );`)
	if err != nil {
		return nil, err
	}

	// User Photos table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_photos (
        user_id TEXT NOT NULL,
        photo_id TEXT NOT NULL,
        PRIMARY KEY (user_id, photo_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (photo_id) REFERENCES new_photos(photo_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Comment table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
        comment_id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        photo_id TEXT NOT NULL,
        content TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (photo_id) REFERENCES new_photos(photo_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Like table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS likes (
        user_id TEXT NOT NULL,
        photo_id TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        PRIMARY KEY (user_id, photo_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (photo_id) REFERENCES new_photos(photo_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Photo table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS new_photos (
        photo_id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
		image_data BLOB,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(user_id)
    );`)
	if err != nil {
		return nil, err
	}

	// Ban table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS new_bans (
        ban_id TEXT PRIMARY KEY,
        banned_by TEXT NOT NULL,
        banned_user TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY (banned_by) REFERENCES users(user_id),
        FOREIGN KEY (banned_user) REFERENCES users(user_id)
    );`)
	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
