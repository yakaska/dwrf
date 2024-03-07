package model

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("id not found")
	ErrIdExists = errors.New("id already exists")
)

type Link struct {
	Short       string    `json:"short"`
	Long        string    `json:"long"`
	Visits      uint64    `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	LastVisited time.Time `json:"last_visited"`
}

type LinkInput struct {
	URL       string
	Id        string
	CreatedBy string
}
