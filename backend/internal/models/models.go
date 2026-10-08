package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username  string    `gorm:"unique;not null;size:50"`
	Email     string    `gorm:"unique;not null;size:255"`
	Password  string    `gorm:"not null"`
	Role      UserRole  `gorm:"default:member"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleModerator UserRole = "moderator"
	RoleMember   UserRole = "member"
)

// Node represents a community/cluster (like a subreddit)
type Node struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	HostID      uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"not null;size:100"`
	Slug        string    `gorm:"unique;not null;size:100"`
	Description string    `gorm:"size:500"`
	Rules       string    `gorm:"size:2000"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Host    User       `gorm:"foreignKey:HostID"`
	Members []User     `gorm:"many2many:user_nodes;"`
	Posts   []Post     `gorm:"foreignKey:NodeID"`
}

// Post represents a submission to a node
type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
	NodeID    uuid.UUID      `gorm:"type:uuid;not null"`
	AuthorID  uuid.UUID      `gorm:"type:uuid;not null"`
	Title     string         `gorm:"not null;size:300"`
	Body      string         `gorm:"type:text"`
	URL       *string        `gorm:"size:500"`
	VoteCount int            `gorm:"default:0"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt *time.Time     `gorm:"index"`

	Node   Node   `gorm:"foreignKey:NodeID"`
	Author User   `gorm:"foreignKey:AuthorID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment represents a nested comment on a post
type Comment struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()"`
	PostID    uuid.UUID  `gorm:"type:uuid;not null"`
	AuthorID  uuid.UUID  `gorm:"type:uuid;not null"`
	ParentID  *uuid.UUID `gorm:"type:uuid"`
	Body      string     `gorm:"type:text;not null"`
	VoteCount int        `gorm:"default:0"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`

	Post   Post     `gorm:"foreignKey:PostID"`
	Author User     `gorm:"foreignKey:AuthorID"`
	Parent *Comment `gorm:"foreignKey:ParentID"`
	Replies []Comment `gorm:"foreignKey:ParentID"`
}

// Vote represents a user's vote on a post or comment
type Vote struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	PostID     *uuid.UUID `gorm:"type:uuid"`
	CommentID  *uuid.UUID `gorm:"type:uuid"`
	VoteType   VoteType  `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	User   User  `gorm:"foreignKey:UserID"`
}

// VoteType represents the type of vote
type VoteType string

const (
	VoteUp   VoteType = "up"
	VoteDown VoteType = "down"
)

// InviteCode represents an invite code for joining the instance
type InviteCode struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Code      string    `gorm:"unique;not null;size:20"`
	HostID    uuid.UUID `gorm:"type:uuid;not null"`
	MaxUses   int       `gorm:"default:1"`
	Uses      int       `gorm:"default:0"`
	ExpiresAt *time.Time `gorm:"index"`
	Role      UserRole  `gorm:"default:member"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	DeletedAt *time.Time `gorm:"index"`

	Host User `gorm:"foreignKey:HostID"`
}

// Tag represents a category/tag for posts (Lobsters-style)
type Tag struct {
	ID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name  string    `gorm:"unique;not null;size:50"`
	Posts []Post    `gorm:"many2many:post_tags;"`
}

// PostTag is the join table for posts and tags
type PostTag struct {
	PostID uuid.UUID `gorm:"primaryKey;type:uuid"`
	TagID  uuid.UUID `gorm:"primaryKey;type:uuid"`
}

// UserNode is the join table for users and nodes (membership)
type UserNode struct {
	UserID uuid.UUID `gorm:"primaryKey;type:uuid"`
	NodeID uuid.UUID `gorm:"primaryKey;type:uuid"`
}

// Moderator represents node-level moderation
type Moderator struct {
	NodeID uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID uuid.UUID `gorm:"primaryKey;type:uuid"`

	Node Node `gorm:"foreignKey:NodeID"`
	User User `gorm:"foreignKey:UserID"`
}
