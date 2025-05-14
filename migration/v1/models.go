package main

//
//import (
//	"github.com/google/uuid"
//	"time"
//)
//
////type UserData struct {
////	UserID        string      `json:"sub" gorm:"column:sub"`
////	Username      string      `json:"cognito:username" gorm:"column:cognito_username"`
////	Email         string      `json:"email" gorm:"column:email"`
////	EmailVerified bool        `json:"email_verified" gorm:"column:email_verified"`
////	FullName      string      `json:"fullName" gorm:"column:full_name"`
////	FirstName     string      `json:"given_name" gorm:"column:first_name"`
////	LastName      string      `json:"family_name" gorm:"column:last_name"`
////	Country       string      `json:"country" gorm:"column:country"`
////	KycVerified   bool        `json:"kyc_verified" gorm:"column:kyc_verified"`
////	BioMetricHash string      `json:"bio_metric_hash" gorm:"column:bio_metric_hash"`
////	AccessToken   AccessToken `json:"access_token" gorm:"column:access_token"`
////	IdTOKEN       IdTOKEN     `json:"id_token" gorm:"column:id_token"`
////}
////
////type AccessToken struct {
////	AccessToken string `json:"access_token" gorm:"column:access_token"`
////}
////
////type IdTOKEN struct {
////	IdToken string `json:"id_token" gorm:"column:id_token"`
////}
////
//// Badge maps to the "badge" table
//// PK: id TEXT
//// Columns: label, icon, color, created_at, updated_at
//// Relations: none
//
//type Badge struct {
//	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	Label     string    `gorm:"column:label;type:text;not null"`
//	Icon      string    `gorm:"column:icon;type:text;not null"`
//	Color     string    `gorm:"column:color;type:text;not null"`
//	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
//}
//
//func NewBadge(label, icon, color string) Badge {
//	return Badge{
//		ID:        uuid.New(),
//		Label:     label,
//		Icon:      icon,
//		Color:     color,
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//	}
//}
//func (Badge) TableName() string { return "badge" }
//
//// UserBadge maps to the "user_badge" table
//// PK: id TEXT
//// FKs: user_id -> user.id, badge_id -> badge.id
//
//type UserBadge struct {
//	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
//	BadgeID    uuid.UUID `gorm:"column:badge_id;type:uuid;not null"`
//	Tier       string    `gorm:"column:tier;type:text;not null;default:'BRONZE'"`
//	IsFeatured bool      `gorm:"column:is_featured;not null;default:false"`
//	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//	Badge Badge `gorm:"foreignKey:BadgeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//}
//
//func (UserBadge) TableName() string { return "user_badge" }
//
//// User maps to the "user" table
//
//type User struct {
//	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	FirstName         string     `gorm:"column:first_name;type:text;not null"`
//	LastName          string     `gorm:"column:last_name;type:text;not null"`
//	Email             string     `gorm:"column:email;type:text;not null;uniqueIndex"`
//	CognitoUsername   string     `gorm:"column:cognito_user_name;type:text;not null"`
//	Verified          bool       `gorm:"column:verified;not null;default:false"`
//	TwoFactorVerified bool       `gorm:"column:two_factor_verified;not null;default:false"`
//	Username          string     `gorm:"column:username;type:text;not null;"`
//	Avatar            *string    `gorm:"column:avatar;type:text"`
//	Country           string     `gorm:"column:country;type:text;not null"`
//	WalletCreated     bool       `gorm:"column:wallet_created;not null;default:false"`
//	WalletCreatedTime *time.Time `gorm:"column:wallet_created_time"`
//	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
//
//	UserBadges     []UserBadge `gorm:"foreignKey:UserID"`
//	Skills         []UserSkill `gorm:"foreignKey:UserID"`
//	Gigs           []Gig       `gorm:"foreignKey:SellerID"`
//	OrdersAsBuyer  []Order     `gorm:"foreignKey:BuyerID"`
//	OrdersAsSeller []Order     `gorm:"foreignKey:SellerID"`
//}
//
//func (User) TableName() string { return "user" }
//
//// Skill maps to the "skill" table
//
//type Skill struct {
//	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	Label     string    `gorm:"column:label;type:text;not null;uniqueIndex"`
//	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	Users []UserSkill `gorm:"foreignKey:SkillID"`
//}
//
//func NewSkill(label string) Skill {
//	return Skill{
//		ID:    uuid.New(),
//		Label: label,
//	}
//}
//
//func (Skill) TableName() string { return "skill" }
//
//// UserSkill maps to the "user_skills" join table
//
//type UserSkill struct {
//	SkillID   uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey"`
//	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
//	Level     int       `gorm:"column:level;not null;default:1"`
//	Endorsed  bool      `gorm:"column:endorsed;not null;default:false"`
//	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//	Skill Skill `gorm:"foreignKey:SkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//}
//
//func (UserSkill) TableName() string { return "user_skills" }
//
//// Biometrics maps to the "biometrics" table
//
//type Biometrics struct {
//	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	Type            string    `gorm:"column:type;type:text;not null"`
//	CognitoUsername *string   `gorm:"column:cognito_user_name;type:uuid;not null"`
//	Value           string    `gorm:"column:value;type:text;not null"`
//	IsVerified      bool      `gorm:"column:is_verified;not null;default:false"`
//	UserID          uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
//	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//}
//
//func (Biometrics) TableName() string { return "biometrics" }
//
//// GigTag maps to the "gig_tag" table
//
//type GigTag struct {
//	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	Label     string    `gorm:"column:label;type:text;not null;uniqueIndex"`
//	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	Gigs []Gig `gorm:"many2many:_gig_to_gig_tag;foreignKey:ID;joinForeignKey:B;References:ID;joinReferences:A"`
//}
//
//func (GigTag) TableName() string { return "gig_tag" }
//
//func NewGigTag(label string) GigTag {
//	return GigTag{
//		ID:    uuid.New(),
//		Label: label,
//	}
//}
//
//// Gig maps to the "gig" table
//
//type Gig struct {
//	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	Title         string    `gorm:"column:title;type:text;not null"`
//	Description   string    `gorm:"column:description;type:text;not null"`
//	IsActive      bool      `gorm:"column:is_active;not null;default:true"`
//	ViewCount     int       `gorm:"column:view_count;not null;default:0"`
//	AverageRating float64   `gorm:"column:average_rating;not null;default:0"`
//	RatingCount   int       `gorm:"column:rating_count;not null;default:0"`
//	CategoryID    uuid.UUID `gorm:"column:category_id;type:uuid;not null"`
//	SellerID      uuid.UUID `gorm:"column:seller_id;type:uuid;not null"`
//	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	Category  Category   `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//	Seller    User       `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
//	Tags      []GigTag   `gorm:"many2many:_gig_to_gig_tag;foreignKey:ID;joinForeignKey:GigID;References:ID;joinReferences:TagID"`
//	GigImages []GigImage `gorm:"foreignKey:GigID"`
//}
//
//func (Gig) TableName() string { return "gig" }
//
//// GigImage maps to the "gig_image" table
//
//type GigImage struct {
//	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	GigID     uuid.UUID `gorm:"column:gig_id;type:uuid;not null"`
//	ImageURL  string    `gorm:"column:image_url;type:text;not null"`
//	IsPrimary bool      `gorm:"column:is_primary;not null;default:false"`
//	SortOrder int       `gorm:"column:sort_order;not null;default:0"`
//	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	Gig Gig `gorm:"foreignKey:GigID"`
//}
//
//func (GigImage) TableName() string { return "gig_image" }
//
//// GigPackage maps to the "gig_package" table
//
//type GigPackage struct {
//	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	GigID       uuid.UUID `gorm:"column:gig_id;type:uuid;not null"`
//	PackageName string    `gorm:"column:package_name;type:text;not null"`
//	Description string    `gorm:"column:description;type:text"`
//	Price       float64   `gorm:"column:price;not null"`
//	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	Gig Gig `gorm:"foreignKey:GigID"`
//}
//
//func (GigPackage) TableName() string { return "gig_package" }
//
//// Category maps to the "category" table
//
//type Category struct {
//	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	Name      string     `gorm:"column:name;type:text;not null"`
//	ParentID  *uuid.UUID `gorm:"column:parent_id;type:uuid"`
//	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
//
//	Parent        *Category  `gorm:"foreignKey:ParentID"`
//	Subcategories []Category `gorm:"foreignKey:ParentID"`
//}
//
//func NewCategory(name string) Category {
//	return Category{
//		Name: name,
//	}
//}
//func (Category) TableName() string { return "category" }
//
//// Order maps to the "order" table
//
//type Order struct {
//	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
//	BuyerID     uuid.UUID `gorm:"column:buyer_id;type:uuid;not null"`
//	SellerID    uuid.UUID `gorm:"column:seller_id;type:uuid;not null"`
//	GigID       uuid.UUID `gorm:"column:gig_id;type:uuid;not null"`
//	Price       float64   `gorm:"column:price;not null"`
//	Quantity    int       `gorm:"column:quantity;not null"`
//	TotalAmount float64   `gorm:"column:total_amount;not null"`
//	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
//	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
//
//	Buyer  User `gorm:"foreignKey:BuyerID"`
//	Seller User `gorm:"foreignKey:SellerID"`
//	Gig    Gig  `gorm:"foreignKey:GigID"`
//}
//
//func (Order) TableName() string { return "order" }
