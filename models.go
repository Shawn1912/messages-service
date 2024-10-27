package main

import "time"

type Message struct {
	ID           int64     `json:"id"`
	Content      string    `json:"text"`
	IsPalindrome bool      `json:"isPalindrome"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
