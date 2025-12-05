package dtos

import "time"

type UpdateUrlRequest struct {
	NextCheck *time.Time `json:"next_check"`
}