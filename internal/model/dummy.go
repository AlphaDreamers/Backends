package model

import "github.com/google/uuid"

type UserDashBoard struct {
	UserProfileInfo UserProfileInfo   `json:"userProfileInfo"`
	Chats           []UserCurrentChat `json:"chats"`
	Orders          []OrderReceived   `json:"orders"`
	UserOfferGigs   []UserOfferGig    `json:"userOfferGigs"`
}
type UserProfileInfo struct {
	UserId      uuid.UUID `json:"userId"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	KycVerified bool      `json:"kycVerified"`
	UserBadges  []string  `json:"userBadges"`
}

type UserCurrentChat struct {
	UserId  uuid.UUID `json:"userId"`
	ChatId  uuid.UUID `json:"chatId"`
	BuyerId uuid.UUID `json:"buyerId"`
}

type OrderReceived struct {
	OrderId     uuid.UUID `json:"orderId"`
	OrderStatus string    `json:"orderStatus"`
	CreatedAt   int64     `json:"createdAt"`
}

type UserOfferGig struct {
	GIgId           uuid.UUID `json:"gigId"`
	GIgDescriptsion string    `json:"gigDescription"`
	GigFeatures     []string  `json:"gigFeatures"`
}
