package usernamedb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *UsernameDB

func TestMain(m *testing.M) {
	db = &UsernameDB{
		userMap: map[string]string{
			"foo@example.com": "foo",
			"bar@example.com": "bar",
		},
		AddUnresolved: false,
	}

	os.Exit(m.Run())
}

func TestToUsername(t *testing.T) {
	username, ok := db.ToUsername("foo@example.com")
	assert.True(t, ok)
	assert.Equal(t, "@foo", username)
}

func TestToUsernameUnresolved(t *testing.T) {
	_, ok := db.ToUsername("qux@example.com")
	assert.False(t, ok)
}

func TestToUsernames(t *testing.T) {
	usernames := db.ToUsernames([]string{"foo@example.com", "bar@example.com"})
	assert.Equal(t, []string{"@foo", "@bar"}, usernames)
}

func TestToUsernamesWithUnresolved(t *testing.T) {
	db := &UsernameDB{
		userMap: map[string]string{
			"foo@example.com": "foo",
			"bar@example.com": "bar",
		},
		AddUnresolved: true,
	}
	usernames := db.ToUsernames([]string{"foo@example.com", "bar@example.com", "qux@example.com"})
	assert.Equal(t, []string{"@foo", "@bar", "qux@example.com"}, usernames)
}
