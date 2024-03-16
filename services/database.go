package services

import(
	"time"
	"github.com/google/uuid"

)

type User struct {
    ID        uuid.UUID `json:"userId"` // Unique identifier
    Username  string   `json:"username"`
    Followers []string  // IDs of followers
    Following []string // IDs of users being followed
    Photos    []string // IDs of photos uploaded by the user
}

type Comment struct {
	ID      string  `json:"commentId"` // Unique identifier
	UserID  string  `json:"userId"`  // ID of the user who commented
	PhotoID string  `json:"photoId"`  // ID of the photo being commented on
	Content string  `json:"content"`  // The comment itself
	Timestamp    time.Time `json:"timestamp"` // Timestamp of when the comment was made
}

type Like struct {
	UserID  string `json:"userId"` // ID of the user who liked the photo
	PhotoID string `json:"photoId"` // ID of the photo being liked
	Timestamp    time.Time `json:"timestamp"` // Timestamp of when the like was made
}


type Photo struct {
    ID       string `json:"photoId"` // Unique identifier
    UserID   string  `json:"userId"` // ID of the user who uploaded the photo
    URL      string  `json:"url"` // URL of the photo
    UploadTime time.Time `json:"uploadTime"` // Timestamp of when the photo was uploaded
    Likes    []Like   
    Comments []Comment
}

type Ban struct {
    ID        string  `json:"banId"` // Unique identifier
    BannedBy  string  `json:"bannedBy"` // ID of the user who banned the other user
    BannedUser string  `json:"bannedUser"` // ID of the user who was banned'
    Time      time.Time `json:"time"` // Timestamp of when the ban was made
}

