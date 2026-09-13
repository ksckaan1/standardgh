package main

import (
	"context"
	"net/http"
)

type standardRequest struct {
	Email  string `json:"email"`
	UserID int    `header:"User-ID,required" json:"-"`
	Limit  int    `query:"limit" json:"-"`
}

type standardResponse struct {
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
	Limit  int    `json:"limit"`
}

func standardHandler(ctx context.Context, req *standardRequest) (*standardResponse, int, error) {
	return &standardResponse{
		Email:  req.Email,
		UserID: req.UserID,
		Limit:  req.Limit,
	}, http.StatusOK, nil
}
