package model

import "github.com/google/uuid"

type GigCreationReq struct {
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	IsActive     bool                `json:"is_active"`
	SellerId     uuid.UUID           `json:"seller_id"`
	CategoryName string              `json:"category_name"`
	GigPkg       []GigPkgCreationReq `json:"gig_pkg"`
}

type GigImageCreationReq struct {
}

type GigPkgCreationReq struct {
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	Price        float32                    `json:"price"`
	DeliveryTime int                        `json:"delivery_time"`
	Revision     int                        `json:"revision"`
	IsActive     bool                       `json:"is_active"`
	Features     []GigPkgFeatureCreationReq `json:"features"`
}

type GigPkgFeatureCreationReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Included    bool   `json:"included"`
}
