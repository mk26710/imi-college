package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Password{},
		&UserToken{},
		&UserDetails{},
		&UserAddress{},
		&UserFile{},
		&Application{},
		&DictAppStatus{},
		&DictEduDocType{},
		&DictIdDocType{},
		&DictEduLevel{},
		&DictNationality{},
		&DictRegion{},
		&DictTownType{},
		&DictGender{},
		&CollegeMajor{},
		&DocStatus{},
		&IdentityDoc{},
		&EducationDoc{},
	)
}

type User struct {
	ID          uuid.UUID    `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt   time.Time    `gorm:"not null;default:now();" json:"createdAt"`
	UserName    string       `gorm:"not null;uniqueIndex;" json:"username"`
	Email       string       `gorm:"not null;uniqueIndex;" json:"email"`
	IsVerified  bool         `gorm:"not null;default:false;" json:"isVerified"`
	Permissions int64        `gorm:"not null;default:0;" json:"permissions,string"`
	Details     *UserDetails `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"details"`
	Address     *UserAddress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type Password struct {
	ID     uuid.UUID `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();"`
	User   User      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uuid.UUID `gorm:"not null;uniqueIndex;"`
	Hash   string    `gorm:"not null;"`
}

type UserToken struct {
	ID        uuid.UUID `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID    uuid.UUID `gorm:"not null;" json:"userId"`
	CreatedAt time.Time `gorm:"not null;default:now();" json:"createdAt"`
	ExpiresAt time.Time `gorm:"not null;default:now() + interval '2 days';" json:"expiresAt"`
	Token     string    `gorm:"not null;uniqueIndex;" json:"token"`
}

type UserDetails struct {
	ID         uuid.UUID  `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	UserID     uuid.UUID  `gorm:"not null;uniqueIndex;" json:"userId"`
	FirstName  string     `gorm:"not null;" json:"firstName"`
	MiddleName string     `gorm:"not null;" json:"middleName"`
	LastName   *string    `json:"lastName"`
	Gender     DictGender `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GenderID   int        `gorm:"not null;" json:"genderId"`
	Birthday   time.Time  `gorm:"not null;type:date;" json:"birthday"`
	Tel        string     `gorm:"not null;" json:"tel"`
	SNILS      *string    `json:"snils"`
	NeedsDorm  bool       `gorm:"not null;default:false;" json:"needsDorm"`
}

type UserAddress struct {
	ID         uuid.UUID    `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	UserID     uuid.UUID    `gorm:"not null;uniqueIndex;type:uuid;" json:"userId"`
	Region     DictRegion   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"region"`
	RegionID   int          `gorm:"not null;" json:"regionId"`
	TownType   DictTownType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"townType"`
	TownTypeID int          `gorm:"not null;" json:"townTypeId"`
	Town       string       `gorm:"not null;" json:"town"`
	Address    string       `gorm:"not null;" json:"address"`
	PostCode   string       `gorm:"not null;" json:"postCode"`
}

type UserFile struct {
	ID           uuid.UUID `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt    time.Time `gorm:"not null;default:now();" json:"createdAt"`
	SHA256       string    `gorm:"not null,uniqueIndex;" json:"sha256"`
	User         User      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID       uuid.UUID `gorm:"not null;type:uuid;" json:"userId"`
	MimeType     string    `gorm:"not null;" json:"mimeType"`
	AbsolutePath string    `gorm:"not null;" json:"-"`
}

type Application struct {
	ID         uuid.UUID     `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt  time.Time     `gorm:"not null;default:now();" json:"createdAt"`
	UserID     uuid.UUID     `gorm:"not null;type:uuid;" json:"userId"`
	User       User          `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	MajorID    uuid.UUID     `gorm:"not null;type:uuid;" json:"majorId"`
	Major      CollegeMajor  `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	EduLevelID int           `gorm:"not null;" json:"eduLevelId"`
	EduLevel   DictEduLevel  `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StatusID   int           `gorm:"not null;" json:"statusId"`
	Status     DictAppStatus `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"status,omitempty"`
	Priority   uint8         `gorm:"not null;default:1;" json:"priority"`
}

type DictAppStatus struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	IsDefault    bool    `gorm:"not null;default:false;" json:"isDefault"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
}

type DictEduDocType struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
}

type DictIdDocType struct {
	ID           int            `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type DictEduLevel struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
}

type DictNationality struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
	SortPriority int     `gorm:"not null;default:0;" json:"sortPriority"`
}

type DictRegion struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	RegionID     int     `gorm:"not null;uniqueIndex;" json:"regionId"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
	SortPriority int     `gorm:"not null;default:0;" json:"sortPriority"`
}

type DictTownType struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
}

type DictGender struct {
	ID           int     `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	Value        string  `gorm:"not null;" json:"value"`
	DisplayValue *string `json:"displayValue"`
}

type CollegeMajor struct {
	ID           uuid.UUID `gorm:"not null;type:uuid;default:gen_random_uuid();" json:"id"`
	Name         string    `gorm:"not null;" json:"name"`
	Prefix       string    `gorm:"not null;" json:"prefix"`
	Base         string    `gorm:"not null;" json:"base"`
	NameOfficial string    `gorm:"not null;" json:"nameOfficial"`
	Budget       bool      `gorm:"not null;default:false;" json:"budget"` // TODO: there might be a better name for this field
	Code         string    `gorm:"not null;" json:"code"`
}

type DocStatus struct {
	ID           int    `gorm:"not null;primaryKey;autoIncrement:false;" json:"id"`
	IsDefault    bool   `gorm:"not null;default:false;" json:"isDefault"`
	Value        string `gorm:"not null;" json:"value"`
	DisplayValue string `json:"displayValue"`
}

type IdentityDoc struct {
	ID            uuid.UUID       `gorm:"not null;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt     time.Time       `gorm:"not null;default:now();" json:"createdAt"`
	UserID        uuid.UUID       `gorm:"not null;type:uuid;" json:"userId"`
	User          User            `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StatusID      int             `json:"statusId"`
	Status        DocStatus       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"status"`
	TypeID        int             `gorm:"not null;" json:"typeId"`
	Type          DictIdDocType   `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"type"`
	Series        string          `gorm:"not null;" json:"series"`
	Number        string          `gorm:"not null;" json:"number"`
	Issuer        string          `gorm:"not null;" json:"issuer"`
	IssuedAt      time.Time       `gorm:"not null;type:date;" json:"issuedAt"`
	DivisionCode  string          `gorm:"not null;" json:"divisionCode"`
	NationalityID int             `gorm:"not null;" json:"nationalityId"`
	Nationality   DictNationality `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type EducationDoc struct {
	ID             uuid.UUID      `gorm:"not null;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt      time.Time      `gorm:"not null;default:now();" json:"createdAt"`
	UserID         uuid.UUID      `gorm:"not null;type:uuid;" json:"userId"`
	User           User           `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StatusID       int            `json:"statusId"`
	Status         DocStatus      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"status"`
	TypeID         int            `gorm:"not null;" json:"typeId"`
	Type           DictEduDocType `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Series         string         `gorm:"not null;" json:"series"`
	Number         string         `gorm:"not null;" json:"number"`
	Issuer         string         `gorm:"not null;" json:"issuer"`
	IssuedAt       time.Time      `gorm:"not null;type:date;" json:"issuedAt"`
	GradYear       int16          `gorm:"not null;" json:"gradYear"`
	IssuerRegionID int            `gorm:"not null;" json:"issuerRegionId"`
	IssuerRegion   DictRegion     `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
