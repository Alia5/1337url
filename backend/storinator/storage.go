package storinator

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBSettings struct {
	Host         string `short:"h" long:"host" help:"Database host" default:"localhost"`
	User         string `short:"u" long:"user" help:"Database user" default:"root"`
	Password     string `short:"P" long:"password" help:"Database password" default:""`
	Database     string `short:"d" long:"database" help:"Database name" default:"urlshort"`
	DatabasePort int16  `short:"D" long:"database-port" help:"Database port" default:"3306"`
}

type Shorty struct {
	gorm.Model
	ShortText string `gorm:"index;unique"`
	FullLink  string `gorm:"index;unique"`
	Clicks    uint64
}

type ShortTextExistsError struct {
	ShortText string
	FullLink  string
}

func (e *ShortTextExistsError) Error() string {
	return fmt.Sprintf("ShortText %s already exists for FullLink %s", e.ShortText, e.FullLink)
}

type FindBy uint

const (
	ShortText FindBy = iota
	FullLink
)

var db *gorm.DB

func ConnectDB(settings DBSettings) {
	var err error
	db, err = gorm.Open(
		mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				settings.User,
				settings.Password,
				settings.Host,
				settings.DatabasePort,
				settings.Database,
			),
		),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Shorty{})
}

func FindLink(by FindBy, value string, increaseClick bool) (link Shorty, err error) {
	var finder Shorty
	switch by {
	// can I just say "FUCK `gofmt`s switch indentation! srlsy"?
	case ShortText:
		finder = Shorty{ShortText: value}
	case FullLink:
		finder = Shorty{FullLink: value}
	}
	err = db.Where(&finder).First(&link).Error
	if (err == nil) && increaseClick {
		link.Clicks++
		db.Save(&link)
	}
	return link, err
}

func CreateLink(shortText string, fullLink string) (existingShort string, created bool, err error) {
	var link Shorty
	err = db.Where(&Shorty{FullLink: fullLink}).First(&link).Error
	if err == nil {
		return link.ShortText, false, err
	}
	err = db.Where(&Shorty{ShortText: shortText}).First(&link).Error
	if err == nil {
		return "", false, &ShortTextExistsError{ShortText: shortText, FullLink: link.FullLink}
	}
	link = Shorty{FullLink: fullLink, ShortText: shortText}
	err = db.Create(&link).Error
	return link.ShortText, true, err
}
