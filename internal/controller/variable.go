package controller

import (
	"errors"
	"log"
	"strings"
)

type Object interface {
	Name() string
	Close()
}

func NameValid(input string) error {
	// Check for valid characters
	letters := `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	letters += `abcdefghijklmnopqrstuvwxyz`
	validCharacters := letters
	validCharacters += `0123456789`
	validCharacters += `_`
	for index, character := range input {
		characterString := string(character)
		if index == 1 && !strings.Contains(letters, characterString) {
			log.Println("Invalid variable name")
			return errors.New(`"` + characterString + "\" is not a valid 1st character")
		}
		if !strings.Contains(validCharacters, characterString) {
			log.Println("Invalid variable name")
			return errors.New(`"` + characterString + "\" is not a valid character")
		}
	}
	return nil
}
