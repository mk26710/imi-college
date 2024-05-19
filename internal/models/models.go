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
		&UserIdentity{},
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
		&IdentityDoc{},
		&EducationDoc{},
	)
}

const (
	RoleRegular string = "regular"
	RoleStaff   string = "staff"
	RoleAdmin   string = "admin"
)

type User struct {
	ID         uuid.UUID `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt  time.Time `gorm:"not null;default:now();" json:"createdAt"`
	UserName   string    `gorm:"not null;uniqueIndex;" json:"username"`
	Email      string    `gorm:"not null;uniqueIndex;" json:"email"`
	IsVerified bool      `gorm:"not null;default:false;" json:"isVerified"`
	Role       string    `gorm:"not null;default:'regular';" json:"role"`
	NeedsDorm  bool      `gorm:"not null;default:false;" json:"needsDorm"`
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

type UserIdentity struct {
	ID         uuid.UUID      `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	User       User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID     uuid.UUID      `gorm:"not null;uniqueIndex;" json:"userId"`
	FirstName  string         `gorm:"not null;" json:"firstName"`
	MiddleName string         `gorm:"not null;" json:"middleName"`
	LastName   sql.NullString `json:"lastName"`
	Gender     DictGender     `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GenderID   int            `gorm:"not null;" json:"genderId"`
	Birthday   time.Time      `gorm:"not null;" json:"birthday"`
	Tel        string         `gorm:"not null;" json:"tel"`
	SNILS      sql.NullString `json:"snils"`
}

type UserAddress struct {
	ID         uuid.UUID    `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	Region     DictRegion   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	RegionID   int          `gorm:"not null;" json:"regionId"`
	TownType   DictTownType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TownTypeID int          `gorm:"not null;" json:"townTypeId"`
	Town       string       `gorm:"not null;" json:"town"`
	PostCode   string       `gorm:"not null;" json:"postCode"`
}

type UserFile struct {
	ID           uuid.UUID `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt    time.Time `gorm:"not null;default:now();" json:"createdAt"`
	User         User      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID       uuid.UUID `gorm:"not null;type:uuid;" json:"userId"`
	MimeType     string    `gorm:"not null;" json:"mimeType"`
	AbsolutePath string    `gorm:"not null;" json:"-"`
}

type Application struct {
	ID        uuid.UUID     `gorm:"not null;primaryKey;type:uuid;default:gen_random_uuid();" json:"id"`
	User      User          `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID    uuid.UUID     `json:"userId"`
	Status    DictAppStatus `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StatusID  int           `gorm:"not null;" json:"stateId"`
	CreatedAt time.Time     `gorm:"not null;default:now();" json:"createdAt"`
}

type DictAppStatus struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type DictEduDocType struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type DictIdDocType struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type DictEduLevel struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type DictNationality struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
	SortPriority int            `gorm:"not null;default:0;" json:"sortPriority"`
}

type DictRegion struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	RegionID     int            `gorm:"not null;uniqueIndex;" json:"regionId"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
	SortPriority int            `gorm:"not null;default:0;" json:"sortPriority"`
}

type DictTownType struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type DictGender struct {
	ID           int            `gorm:"not null;primaryKey;" json:"id"`
	Value        string         `gorm:"not null;" json:"value"`
	DisplayValue sql.NullString `json:"displayValue"`
}

type IdentityDoc struct {
	ID            uuid.UUID       `gorm:"not null;type:uuid;default:gen_random_uuid();" json:"id"`
	User          User            `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID        uuid.UUID       `gorm:"not null;type:uuid;" json:"userId"`
	Type          DictIdDocType   `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TypeID        int             `gorm:"not null;" json:"typeId"`
	Series        string          `gorm:"not null;" json:"series"`
	Number        string          `gorm:"not null;" json:"number"`
	Issuer        string          `gorm:"not null;" json:"issuer"`
	IssuedAt      time.Time       `gorm:"not null;" json:"issuedAt"`
	DivisionCode  string          `gorm:"not null;" json:"divisionCode"`
	Nationality   DictNationality `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	NationalityID int             `gorm:"not null;" json:"nationalityId"`
	File          UserFile        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	FileID        uuid.UUID       `gorm:"type:uuid;" json:"fileId"`
}

type EducationDoc struct {
	ID             uuid.UUID      `gorm:"not null;type:uuid;default:gen_random_uuid();" json:"id"`
	CreatedAt      time.Time      `gorm:"not null;default:now();" json:"createdAt"`
	User           User           `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID         uuid.UUID      `gorm:"not null;type:uuid;" json:"userId"`
	Type           DictEduDocType `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TypeID         int            `gorm:"not null;" json:"typeId"`
	Series         string         `gorm:"not null;" json:"series"`
	Number         string         `gorm:"not null;" json:"number"`
	Issuer         string         `gorm:"not null;" json:"issuer"`
	IssuedAt       time.Time      `gorm:"not null;" json:"issuedAt"`
	GradYear       uint8          `gorm:"not null;" json:"gradYear"`
	IssuerRegion   DictRegion     `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	IssuerRegionID int            `gorm:"not null;" json:"issuerRegionId"`
}
