package usernamedb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContributorsFileDBLoaderLoad(t *testing.T) {
	loader := &ContributorsFileDBLoader{
		Filename: "fixtures/CONTRIBUTORS",
	}

	users, err := loader.Load()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"foo@bar.com": "foo",
	}, users)
}

func TestSimpleDBLoader(t *testing.T) {
	loader := &SimpleDBLoader{}

	loader.AddEntry("foo@bar.com", "foo")
	assert.Equal(t, map[string]string{
		"foo@bar.com": "foo",
	}, loader.users)

	users, err := loader.Load()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"foo@bar.com": "foo",
	}, users)
}
