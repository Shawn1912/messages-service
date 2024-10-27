package main

import "time"

type Message struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	IsPalindrome bool      `json:"isPalindrome"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
