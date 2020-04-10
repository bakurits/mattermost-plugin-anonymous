package test

import "github.com/stretchr/testify/assert"

// CheckErr checks if function returned error correctly
func CheckErr(tassert *assert.Assertions, wantErr bool, err error) {

	if wantErr {
		tassert.Error(err)
	} else {
		tassert.NoError(err)
	}
}
