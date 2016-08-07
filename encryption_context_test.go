package unicreds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	encContext := NewEncryptionContextValue()

	err := encContext.Set("someval:newly")

	assert.Nil(t, err)

	err = encContext.Set("booo")

	assert.Error(t, err)
}

func TestIsCumulative(t *testing.T) {
	encContext := NewEncryptionContextValue()

	assert.Equal(t, true, encContext.IsCumulative())
}
