package updater

import (
	"errors"
	"os"
)

type Config struct {
	IsUpdating bool
	Username   string
	Password   string
	Enabled    bool
	Workers    int
	Empty      bool
	Port       int
}

func (u Config) CheckInfo() error {
	if !u.Enabled {
		return nil
	}

	if u.Username == "{FROM_ENV}" || u.Username == "" {
		u.Username = os.Getenv("fileb0x_username")
	}

	if u.Password == "{FROM_ENV}" || u.Password == "" {
		u.Password = os.Getenv("fileb0x_password")
	}

	// check for empty username and password
	if u.Username == "" {
		return errors.New("fileb0x: You must provide an username in the config file or through an env var: fileb0x_username")

	} else if u.Password == "" {
		return errors.New("fileb0x: You must provide an password in the config file or through an env var: fileb0x_password")
	}

	return nil
}
