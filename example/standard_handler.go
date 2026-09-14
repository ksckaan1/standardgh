package main

import (
	"context"
	"net/http"
)

type standardRequest struct {
	ID     int    `uri:"id" json:"-"`
	Email  string `json:"email"`
	UserID int    `header:"User-ID" json:"-"`
	Limit  int    `query:"limit" json:"-"`
}

type standardResponse struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
	Limit  int    `json:"limit"`
}

func standardHandler(ctx context.Context, req *standardRequest) (*standardResponse, int, error) {
	return &standardResponse{
		ID:     req.ID,
		Email:  req.Email,
		UserID: req.UserID,
		Limit:  req.Limit,
	}, http.StatusOK, nil
}
