package main

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Text   string `json:"text"`
	IsDone bool   `json:"is_done"`
}
