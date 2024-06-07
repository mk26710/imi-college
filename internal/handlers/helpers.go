package handlers

import (
	"gorm.io/gorm"
)

type HandlersMap struct {
	Dictionaries DictionariesHandler
	Users        UserHandler
	Tokens       TokensHandler
	Address      AddressHandler
	Files        FilesHandler
	Documents    DocumentsHandler
}

func Create(db *gorm.DB) HandlersMap {
	if db == nil {
		panic("database connection cannot be null! never! neeeverrrr!!!")
	}

	return HandlersMap{
		Dictionaries: DictionariesHandler{db},
		Users:        UserHandler{db},
		Tokens:       TokensHandler{db},
		Address:      AddressHandler{db},
		Files:        FilesHandler{db},
		Documents:    DocumentsHandler{db},
	}
}
