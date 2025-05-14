package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var (
	Badges          []Badge
	Skills          []Skill
	Categories      []Category
	Users           []User
	Gigs            []Gig
	GigTags         []GigTag
	UserBadges      []UserBadge
	GigImages       []GigImage
	GigPackages     []GigPackage
	Orders          []Order
	UserSkills      []UserSkill
	firstNames      []string
	lastNames       []string
	emails          []string
	cognitoUsername []string
	usernames       []string
)

func PopulateBadges() {
	labels := []string{"Newbie", "ProSeller", "FastDelivery"}
	colors := []string{"red", "green", "blue"}
	icons := []string{"🟡", "star", "dummy"}
	for i := 0; i < len(labels); i++ {
		Badges = append(Badges, Badge{
			Label: labels[i],
			Color: colors[i],
			Icon:  icons[i],
		})
	}
}

func PopulateSkills() {
	labels := []string{"web development", "teaching", "web3"}
	for _, label := range labels {
		Skills = append(Skills, Skill{
			Label: label,
		})
	}
}

func PopulateCategories() {
	labels := []string{"backend", "frontend", "crypto"}
	for _, label := range labels {
		Categories = append(Categories, Category{
			Name: label,
		})
	}
}

func PopulateGigTags() {
	labels := []string{
		"web", "mobile", "design", "writing", "marketing",
		"video", "audio", "programming", "business", "lifestyle",
	}
	for _, label := range labels {
		GigTags = append(GigTags, NewGigTag(label))
	}
}

func PopulateUsers() {
	firstNames = []string{"John", "Jane", "Alex", "Chris", "Sam", "Taylor", "Jordan", "Pat", "Casey", "Morgan"}
	lastNames = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Davis", "Miller", "Wilson", "Moore", "Taylor"}
	emails = []string{
		"john.smith@example.com", "jane.johnson@example.com", "alex.williams@example.com", "chris.brown@example.com",
		"sam.jones@example.com", "taylor.davis@example.com", "jordan.miller@example.com", "pat.wilson@example.com",
		"casey.moore@example.com", "morgan.taylor@example.com",
	}
	cognitoUsername = []string{
		"john_smith123", "jane_johnson456", "alex_williams789", "chris_brown101", "sam_jones202",
		"taylor_davis303", "jordan_miller404", "pat_wilson505", "casey_moore606", "morgan_taylor707",
	}
	usernames = []string{
		"johnnyS", "janeyJ", "alexTheGreat", "chrisB123", "samJ99",
		"taylorT11", "jordanM_27", "patWilson007", "caseyM_45", "morganT_88",
	}
	for i := 0; i < len(firstNames); i++ {
		Users = append(Users, User{
			FirstName:         firstNames[i],
			LastName:          lastNames[i],
			Email:             emails[i],
			CognitoUsername:   cognitoUsername[i],
			Verified:          true,
			TwoFactorVerified: true,
			Username:          usernames[i],
			WalletCreated:     true,
			Country:           "USA",
		})
	}
}

func PopulateUserSkills() {
	for i, user := range Users {
		skillIndex := i % len(Skills)
		UserSkills = append(UserSkills, UserSkill{
			SkillID:  Skills[skillIndex].ID,
			UserID:   user.ID,
			Level:    5,
			Endorsed: true,
		})
	}
}

func PopulateUserBadges() {
	tiers := []string{"BRONZE", "SILVER", "GOLD"}
	for i, user := range Users {
		for j := 0; j < 2; j++ { // 2 badges per user
			badgeIndex := (i + j) % len(Badges)
			tierIndex := j % len(tiers)
			isFeatured := (j == 0)
			UserBadges = append(UserBadges, UserBadge{
				UserID:     user.ID,
				BadgeID:    Badges[badgeIndex].ID,
				Tier:       tiers[tierIndex],
				IsFeatured: isFeatured,
			})
		}
	}
}

func PopulateGigs() {
	for i, user := range Users {
		for j := 0; j < 2; j++ { // 2 gigs per user
			categoryIndex := (i + j) % len(Categories)
			category := Categories[categoryIndex]
			tagIndex1 := (i + j) % len(GigTags)
			tagIndex2 := (i + j + 1) % len(GigTags)
			gig := Gig{
				Title:       fmt.Sprintf("%s Service by %s", category.Name, user.Username),
				Description: fmt.Sprintf("This is a great service for %s needs.", category.Name),
				CategoryID:  category.ID,
				SellerID:    user.ID,
				IsActive:    true,
				Tags:        []GigTag{GigTags[tagIndex1], GigTags[tagIndex2]},
			}
			Gigs = append(Gigs, gig)
		}
	}
}

func PopulateGigImages() {
	for _, gig := range Gigs {
		for k := 0; k < 2; k++ { // 2 images per gig
			isPrimary := (k == 0)
			gigImage := GigImage{
				GigID:     gig.ID,
				ImageURL:  fmt.Sprintf("http://example.com/gig%d_image%d.jpg", k+1, k+1),
				IsPrimary: isPrimary,
				SortOrder: k,
			}
			GigImages = append(GigImages, gigImage)
		}
	}
}

func PopulateGigPackages() {
	packageNames := []string{"Basic", "Standard", "Premium"}
	prices := []float64{10.0, 20.0, 30.0}
	for _, gig := range Gigs {
		for m := 0; m < 3; m++ { // 3 packages per gig
			gigPackage := GigPackage{
				GigID:       gig.ID,
				PackageName: packageNames[m],
				Description: fmt.Sprintf("%s package for %s", packageNames[m], gig.Title),
				Price:       prices[m],
			}
			GigPackages = append(GigPackages, gigPackage)
		}
	}
}

func PopulateOrders() {
	for i, user := range Users {
		sellerIndex := (i + 1) % len(Users)
		seller := Users[sellerIndex]
		var sellerGigs []Gig
		for _, gig := range Gigs {
			if gig.SellerID == seller.ID {
				sellerGigs = append(sellerGigs, gig)
				break // Take first gig
			}
		}
		if len(sellerGigs) > 0 {
			gig := sellerGigs[0]
			order := Order{
				BuyerID:     user.ID,
				SellerID:    seller.ID,
				GigID:       gig.ID,
				Price:       10.0,
				Quantity:    1,
				TotalAmount: 10.0,
			}
			Orders = append(Orders, order)
		}
	}
}

func Insertion(db *gorm.DB) {
	log.Info("🏁 Starting database population...")
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := createBatch(tx, &Badges, "Badges"); err != nil {
			return err
		}
		if err := createBatch(tx, &Skills, "Skills"); err != nil {
			return err
		}
		if err := createBatch(tx, &Categories, "Categories"); err != nil {
			return err
		}
		if err := createBatch(tx, &GigTags, "Gig Tags"); err != nil {
			return err
		}
		if err := createBatch(tx, &Users, "Users"); err != nil {
			return err
		}
		PopulateUserSkills()
		if err := createBatch(tx, &UserSkills, "User Skills"); err != nil {
			return err
		}
		PopulateUserBadges()
		if err := createBatch(tx, &UserBadges, "User Badges"); err != nil {
			return err
		}
		PopulateGigs()
		if err := createBatch(tx, &Gigs, "Gigs"); err != nil {
			return err
		}
		PopulateGigImages()
		if err := createBatch(tx, &GigImages, "Gig Images"); err != nil {
			return err
		}
		PopulateGigPackages()
		if err := createBatch(tx, &GigPackages, "Gig Packages"); err != nil {
			return err
		}
		PopulateOrders()
		if err := createBatch(tx, &Orders, "Orders"); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("❌ Transaction failed: %v", err)
	}
	generateReport(db)
	log.Info("✅ Database population completed successfully")
}

func clearExistingData(db *gorm.DB) {
	db.Exec("DELETE FROM order")
	db.Exec("DELETE FROM gig_package")
	db.Exec("DELETE FROM gig_image")
	db.Exec("DELETE FROM _gig_to_gig_tag")
	db.Exec("DELETE FROM gig")
	db.Exec("DELETE FROM user_badge")
	db.Exec("DELETE FROM user_skills")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM gig_tag")
	db.Exec("DELETE FROM category")
	db.Exec("DELETE FROM skill")
	db.Exec("DELETE FROM badge")
	log.Info("🧹 Cleared existing data from tables")
}

func createBatch(tx *gorm.DB, models interface{}, name string) error {
	result := tx.Create(models)
	if result.Error != nil {
		log.Errorf("❌ Error creating %s: %v", name, result.Error)
		return result.Error
	}
	log.Infof("📦 Created %d %s", result.RowsAffected, name)
	return nil
}

func generateReport(db *gorm.DB) {
	type CountResult struct {
		Table string
		Count int64
	}
	tables := []string{
		"badge", "skill", "category", "gig_tag", "users",
		"user_skills", "user_badge", "gig", "gig_image", "gig_package", "order",
	}
	var results []CountResult
	for _, table := range tables {
		var count int64
		db.Table(table).Count(&count)
		results = append(results, CountResult{Table: table, Count: count})
	}
	log.Info("\n📊 Database Population Report:")
	log.Info("┌───────────────────────┬───────┐")
	log.Info("│        Table          │ Count │")
	log.Info("├───────────────────────┼───────┤")
	for _, res := range results {
		log.Infof("│ %-21s │ %5d │", res.Table, res.Count)
	}
	log.Info("└───────────────────────┴───────┘")
	var firstUser User
	db.Preload("Skills").First(&firstUser)
	log.Infof("\n👤 Sample User: %s %s", firstUser.FirstName, firstUser.LastName)
	log.Infof("📧 Email: %s", firstUser.Email)
	log.Info("🛠️ Associated Skills:")
	for _, skill := range firstUser.Skills {
		log.Infof("  - %s (Level %d)", skill.Skill.Label, skill.Level)
	}
}

// Struct definitions
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

func (Order) TableName() string { return "order" }

// Database initialization
func InitDB() *gorm.DB {
	dsn := "host=localhost user=swan password=swanhtet dbname=appDatabase sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate schema
	err = db.AutoMigrate(
		&Badge{}, &UserBadge{}, &User{}, &Skill{}, &UserSkill{},
		&Biometrics{}, &GigTag{}, &Gig{}, &GigImage{}, &GigPackage{},
		&Category{}, &Order{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate schema: %v", err)
	}

	return db
}

// API handlers
func CreateGig(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type GigInput struct {
			UserID      string   `json:"user_id"`
			CategoryID  string   `json:"category_id"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			TagIDs      []string `json:"tag_ids"`
		}
		var input GigInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid user_id"})
		}
		categoryID, err := uuid.Parse(input.CategoryID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid category_id"})
		}
		var tags []GigTag
		for _, tagID := range input.TagIDs {
			id, err := uuid.Parse(tagID)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid tag_id"})
			}
			var tag GigTag
			if err := db.First(&tag, "id = ?", id).Error; err != nil {
				return c.Status(404).JSON(fiber.Map{"error": "Tag not found"})
			}
			tags = append(tags, tag)
		}
		gig := Gig{
			Title:       input.Title,
			Description: input.Description,
			CategoryID:  categoryID,
			SellerID:    userID,
			IsActive:    true,
			Tags:        tags,
			GigImages: []GigImage{
				{ImageURL: "http://example.com/new_gig_image.jpg", IsPrimary: true, SortOrder: 0},
			},
			GigPackages: []GigPackage{
				{PackageName: "Basic", Description: "Basic package", Price: 15.0},
			},
		}
		if err := db.Create(&gig).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create gig"})
		}
		return c.JSON(fiber.Map{"message": "Gig created", "gig_id": gig.ID})
	}
}

func PlaceOrder(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type OrderInput struct {
			BuyerID  string  `json:"buyer_id"`
			SellerID string  `json:"seller_id"`
			GigID    string  `json:"gig_id"`
			Price    float64 `json:"price"`
		}
		var input OrderInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		buyerID, err := uuid.Parse(input.BuyerID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid buyer_id"})
		}
		sellerID, err := uuid.Parse(input.SellerID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid seller_id"})
		}
		gigID, err := uuid.Parse(input.GigID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid gig_id"})
		}
		order := Order{
			BuyerID:     buyerID,
			SellerID:    sellerID,
			GigID:       gigID,
			Price:       input.Price,
			Quantity:    1,
			TotalAmount: input.Price,
		}
		if err := db.Create(&order).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to place order"})
		}
		return c.JSON(fiber.Map{"message": "Order placed", "order_id": order.ID})
	}
}

func GetUserProfile(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid user_id"})
		}
		var user User
		err = db.Preload("Skills.Skill").Preload("UserBadges.Badge").Preload("Gigs").First(&user, "id = ?", userID).Error
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(user)
	}
}

func main() {
	// Initialize database
	db := InitDB()

	clearExistingData(db)
	PopulateBadges()
	PopulateSkills()
	PopulateCategories()
	PopulateGigTags()
	PopulateUsers()
	Insertion(db)

	// Set up Fiber API
	app := fiber.New()

	// API endpoints
	app.Post("/gigs", CreateGig(db))
	app.Post("/orders", PlaceOrder(db))
	app.Get("/users/:id", GetUserProfile(db))

	// Start server
	log.Fatal(app.Listen(":3000"))
}
