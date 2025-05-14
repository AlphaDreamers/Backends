package model

import (
	"github.com/google/uuid"
	"time"
)

// BaseModel defines common fields for all entities
type BaseModel struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// Badge represents a badge that can be assigned to users
type Badge struct {
	BaseModel
	Label string `gorm:"column:label;type:text;not null"`
	Icon  string `gorm:"column:icon;type:text;not null"`
	Color string `gorm:"column:color;type:text;not null"`
}

func NewBadge(label, icon, color string) Badge {
	return Badge{
		Label: label,
		Icon:  icon,
		Color: color,
	}
}

func (Badge) TableName() string { return "badge" }

// UserBadge represents the association between a User and a Badge
type UserBadge struct {
	BaseModel
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null;index"`
	BadgeID    uuid.UUID `gorm:"column:badge_id;type:uuid;not null;index"`
	Tier       string    `gorm:"column:tier;type:text;not null;default:'BRONZE'"`
	IsFeatured bool      `gorm:"column:is_featured;not null;default:false"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Badge      Badge     `gorm:"foreignKey:BadgeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserBadge) TableName() string { return "user_badge" }

// User represents a user in the marketplace
type User struct {
	BaseModel
	FirstName         string       `gorm:"column:first_name;type:text;not null"`
	LastName          string       `gorm:"column:last_name;type:text;not null"`
	Email             string       `gorm:"column:email;type:text;not null;uniqueIndex"`
	CognitoUsername   string       `gorm:"column:cognito_user_name;type:text;not null;index"`
	Verified          bool         `gorm:"column:verified;not null;default:false"`
	TwoFactorVerified bool         `gorm:"column:two_factor_verified;not null;default:false"`
	Username          string       `gorm:"column:username;type:text;not null;uniqueIndex"`
	Avatar            *string      `gorm:"column:avatar;type:text"`
	Country           string       `gorm:"column:country;type:text;not null"`
	WalletCreated     bool         `gorm:"column:wallet_created;not null;default:false"`
	WalletCreatedTime *time.Time   `gorm:"column:wallet_created_time"`
	UserBadges        []UserBadge  `gorm:"foreignKey:UserID"`
	UserSkills        []UserSkill  `gorm:"foreignKey:UserID"`
	Gigs              []Gig        `gorm:"foreignKey:SellerID"`
	Orders            []Order      `gorm:"foreignKey:BuyerID,SellerID"`
	Biometrics        []Biometrics `gorm:"foreignKey:UserID"`
}

func (User) TableName() string { return "user" }

// Skill represents a skill that users can have
type Skill struct {
	BaseModel
	Label      string      `gorm:"column:label;type:text;not null;uniqueIndex"`
	UserSkills []UserSkill `gorm:"foreignKey:SkillID"`
}

func NewSkill(label string) Skill {
	return Skill{Label: label}
}

func (Skill) TableName() string { return "skill" }

// UserSkill represents the association between a User and a Skill
type UserSkill struct {
	SkillID   uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey;index"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey;index"`
	Level     int       `gorm:"column:level;not null;default:1"`
	Endorsed  bool      `gorm:"column:endorsed;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skill     Skill     `gorm:"foreignKey:SkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserSkill) TableName() string { return "user_skills" }

// Biometrics represents biometric data for a user
type Biometrics struct {
	BaseModel
	Type            string    `gorm:"column:type;type:text;not null"`
	CognitoUsername *string   `gorm:"column:cognito_user_name;type:text;index"`
	Value           string    `gorm:"column:value;type:text;not null"`
	IsVerified      bool      `gorm:"column:is_verified;not null;default:false"`
	UserID          uuid.UUID `gorm:"column:user_id;type:uuid;not null;index"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Biometrics) TableName() string { return "biometrics" }

// GigTag represents a tag for gigs
type GigTag struct {
	BaseModel
	Label string `gorm:"column:label;type:text;not null;uniqueIndex"`
	Gigs  []Gig  `gorm:"many2many:gig_tag_gigs;foreignKey:ID;joinForeignKey:TagID;References:ID;joinReferences:GigID"`
}

func NewGigTag(label string) GigTag {
	return GigTag{Label: label}
}

func (GigTag) TableName() string { return "gig_tag" }

// Gig represents a service offered by a seller
type Gig struct {
	BaseModel
	Title         string       `gorm:"column:title;type:text;not null"`
	Description   string       `gorm:"column:description;type:text;not null"`
	IsActive      bool         `gorm:"column:is_active;not null;default:true"`
	ViewCount     int          `gorm:"column:view_count;not null;default:0"`
	AverageRating float64      `gorm:"column:average_rating;not null;default:0"`
	RatingCount   int          `gorm:"column:rating_count;not null;default:0"`
	CategoryID    uuid.UUID    `gorm:"column:category_id;type:uuid;not null;index"`
	SellerID      uuid.UUID    `gorm:"column:seller_id;type:uuid;not null;index"`
	Category      Category     `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Seller        User         `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags          []GigTag     `gorm:"many2many:gig_tag_gigs;foreignKey:ID;joinForeignKey:GigID;References:ID;joinReferences:TagID"`
	GigImages     []GigImage   `gorm:"foreignKey:GigID"`
	GigPackages   []GigPackage `gorm:"foreignKey:GigID"`
	Orders        []Order      `gorm:"foreignKey:GigID"`
}

func (Gig) TableName() string { return "gig" }

// GigImage represents an image associated with a gig
type GigImage struct {
	BaseModel
	GigID     uuid.UUID `gorm:"column:gig_id;type:uuid;not null;index"`
	ImageURL  string    `gorm:"column:image_url;type:text;not null"`
	IsPrimary bool      `gorm:"column:is_primary;not null;default:false"`
	SortOrder int       `gorm:"column:sort_order;not null;default:0"`
	Gig       Gig       `gorm:"foreignKey:GigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GigImage) TableName() string { return "gig_image" }

// GigPackage represents a pricing package for a gig
type GigPackage struct {
	BaseModel
	GigID       uuid.UUID `gorm:"column:gig_id;type:uuid;not null;index"`
	PackageName string    `gorm:"column:package_name;type:text;not null"`
	Description string    `gorm:"column:description;type:text"`
	Price       float64   `gorm:"column:price;not null"`
	Gig         Gig       `gorm:"foreignKey:GigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GigPackage) TableName() string { return "gig_package" }

// Category represents a category for gigs
type Category struct {
	BaseModel
	Name          string     `gorm:"column:name;type:text;not null"`
	ParentID      *uuid.UUID `gorm:"column:parent_id;type:uuid;index"`
	Parent        *Category  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Subcategories []Category `gorm:"foreignKey:ParentID"`
	Gigs          []Gig      `gorm:"foreignKey:CategoryID"`
}

func NewCategory(name string) Category { return Category{Name: name} }
func (Category) TableName() string     { return "category" }

// Order represents a purchase of a gig
type Order struct {
	BaseModel
	BuyerID     uuid.UUID `gorm:"column:buyer_id;type:uuid;not null;index"`
	SellerID    uuid.UUID `gorm:"column:seller_id;type:uuid;not null;index"`
	GigID       uuid.UUID `gorm:"column:gig_id;type:uuid;not null;index"`
	Price       float64   `gorm:"column:price;not null"`
	Quantity    int       `gorm:"column:quantity;not null"`
	TotalAmount float64   `gorm:"column:total_amount;not null"`
	Buyer       User      `gorm:"foreignKey:BuyerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Seller      User      `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Gig         Gig       `gorm:"foreignKey:GigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Order) TableName() string { return "order" }
