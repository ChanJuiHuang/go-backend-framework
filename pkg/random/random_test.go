package random

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	replacer := strings.NewReplacer("/", "", "+", "")

	for i := 1; i < 101; i++ {
		i := i
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			t.Parallel()
			str := RandomString(i)
			str = replacer.Replace(str)
			assert.Equal(t, i, len(str))
		})
	}
}
