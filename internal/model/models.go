package model

import (
	"github.com/google/uuid"
	"time"
)

type Badge struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string    `gorm:"column:label;type:text;not null"`
	Icon      string    `gorm:"column:icon;type:text;not null"`
	Color     string    `gorm:"column:color;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func NewBadge(label, icon, color string) Badge {
	return Badge{
		ID:    uuid.New(),
		Label: label,
		Icon:  icon,
		Color: color,
	}
}

func (Badge) TableName() string { return "badge" }

type UserBadge struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	BadgeID    uuid.UUID `gorm:"column:badge_id;type:uuid;not null"`
	Tier       string    `gorm:"column:tier;type:text;not null;default:'BRONZE'"`
	IsFeatured bool      `gorm:"column:is_featured;not null;default:false"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Badge      Badge     `gorm:"foreignKey:BadgeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserBadge) TableName() string { return "user_badge" }

type User struct {
	ID                uuid.UUID   `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName         string      `gorm:"column:first_name;type:text;not null"`
	LastName          string      `gorm:"column:last_name;type:text;not null"`
	Email             string      `gorm:"column:email;type:text;not null;uniqueIndex"`
	CognitoUsername   string      `gorm:"column:cognito_user_name;type:text;not null"`
	Verified          bool        `gorm:"column:verified;not null;default:false"`
	TwoFactorVerified bool        `gorm:"column:two_factor_verified;not null;default:false"`
	Username          string      `gorm:"column:username;type:text;not null;"`
	Avatar            *string     `gorm:"column:avatar;type:text"`
	Country           string      `gorm:"column:country;type:text;not null"`
	WalletCreated     bool        `gorm:"column:wallet_created;not null;default:false"`
	WalletCreatedTime *time.Time  `gorm:"column:wallet_created_time"`
	CreatedAt         time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;autoUpdateTime"`
	UserBadges        []UserBadge `gorm:"foreignKey:UserID"`
	Skills            []UserSkill `gorm:"foreignKey:UserID"`
	Gigs              []Gig       `gorm:"foreignKey:SellerID"`
	OrdersAsBuyer     []Order     `gorm:"foreignKey:BuyerID"`
	OrdersAsSeller    []Order     `gorm:"foreignKey:SellerID"`
}

func (User) TableName() string { return "user" }

type Skill struct {
	ID        uuid.UUID   `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string      `gorm:"column:label;type:text;not null;uniqueIndex"`
	CreatedAt time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time   `gorm:"column:updated_at;autoUpdateTime"`
	Users     []UserSkill `gorm:"foreignKey:SkillID"`
}

func NewSkill(label string) Skill {
	return Skill{ID: uuid.New(), Label: label}
}

func (Skill) TableName() string { return "skill" }

type UserSkill struct {
	SkillID   uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	Level     int       `gorm:"column:level;not null;default:1"`
	Endorsed  bool      `gorm:"column:endorsed;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skill     Skill     `gorm:"foreignKey:SkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserSkill) TableName() string { return "user_skills" }

type Biometrics struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Type            string    `gorm:"column:type;type:text;not null"`
	CognitoUsername *string   `gorm:"column:cognito_user_name;type:uuid;not null"`
	Value           string    `gorm:"column:value;type:text;not null"`
	IsVerified      bool      `gorm:"column:is_verified;not null;default:false"`
	UserID          uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Biometrics) TableName() string { return "biometrics" }

type GigTag struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string    `gorm:"column:label;type:text;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Gigs      []Gig     `gorm:"many2many:_gig_to_gig_tag;foreignKey:ID;joinForeignKey:B;References:ID;joinReferences:A"`
}

func (GigTag) TableName() string { return "gig_tag" }

func NewGigTag(label string) GigTag {
	return GigTag{ID: uuid.New(), Label: label}
}

type Gig struct {
	ID            uuid.UUID    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Title         string       `gorm:"column:title;type:text;not null"`
	Description   string       `gorm:"column:description;type:text;not null"`
	IsActive      bool         `gorm:"column:is_active;not null;default:true"`
	ViewCount     int          `gorm:"column:view_count;not null;default:0"`
	AverageRating float64      `gorm:"column:average_rating;not null;default:0"`
	RatingCount   int          `gorm:"column:rating_count;not null;default:0"`
	CategoryID    uuid.UUID    `gorm:"column:category_id;type:uuid;not null"`
	SellerID      uuid.UUID    `gorm:"column:seller_id;type:uuid;not null"`
	CreatedAt     time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	Category      Category     `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Seller        User         `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags          []GigTag     `gorm:"many2many:_gig_to_gig_tag;foreignKey:ID;joinForeignKey:GigID;References:ID;joinReferences:TagID"`
	GigImages     []GigImage   `gorm:"foreignKey:GigID"`
	GigPackages   []GigPackage `gorm:"foreignKey:GigID"`
}

func (Gig) TableName() string { return "gig" }

type GigImage struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	GigID     uuid.UUID `gorm:"column:gig_id;type:uuid;not null"`
	ImageURL  string    `gorm:"column:image_url;type:text;not null"`
	IsPrimary bool      `gorm:"column:is_primary;not null;default:false"`
	SortOrder int       `gorm:"column:sort_order;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Gig       Gig       `gorm:"foreignKey:GigID"`
}

func (GigImage) TableName() string { return "gig_image" }

type GigPackage struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	GigID       uuid.UUID `gorm:"column:gig_id;type:uuid;not null"`
	PackageName string    `gorm:"column:package_name;type:text;not null"`
	Description string    `gorm:"column:description;type:text"`
	Price       float64   `gorm:"column:price;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Gig         Gig       `gorm:"foreignKey:GigID"`
}

func (GigPackage) TableName() string { return "gig_package" }

type Category struct {
	ID            uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string     `gorm:"column:name;type:text;not null"`
	ParentID      *uuid.UUID `gorm:"column:parent_id;type:uuid"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	Parent        *Category  `gorm:"foreignKey:ParentID"`
	Subcategories []Category `gorm:"foreignKey:ParentID"`
}

func NewCategory(name string) Category { return Category{Name: name} }
func (Category) TableName() string     { return "category" }

type Order struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	BuyerID     uuid.UUID `gorm:"column:buyer_id;type:uuid;not null"`
	SellerID    uuid.UUID `gorm:"column:seller_id;type:uuid;not null"`
	GigID       uuid.UUID `gorm:"column:gig_id;type:uuid;not null"`
	Price       float64   `gorm:"column:price;not null"`
	Quantity    int       `gorm:"column:quantity;not null"`
	TotalAmount float64   `gorm:"column:total_amount;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Buyer       User      `gorm:"foreignKey:BuyerID"`
	Seller      User      `gorm:"foreignKey:SellerID"`
	Gig         Gig       `gorm:"foreignKey:GigID"`
}
