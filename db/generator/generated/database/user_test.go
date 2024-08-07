package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByName(t *testing.T) {

	// Test case 1: Successful retrieval of user
	InitEcommerceDb()
	users, err := GetUserByName(context.Background(), EcommerceDb.User, GetUserByNameParams{Name: "John Doe"})
	assert.Nil(t, err)
	assert.Equal(t, "John Doe", users[0].Name)
}
