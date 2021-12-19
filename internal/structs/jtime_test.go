package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJTime_Type(t *testing.T) {
	assert := assert.New(t)

	jtime := JTime{}
	assert.IsType(time.Time{}, time.Time(jtime), "should be type")

}

func TestJTime_MarshalJSON(t *testing.T) {
	assert := assert.New(t)

	timeparsed := time.Date(2021, 12, 31, 12, 0, 0, 0, time.UTC)
	println(timeparsed.String())
	jtime := JTime(timeparsed)

	stamp, _ := jtime.MarshalJSON()

	assert.Equal("\"2021-12-31T12:00:00Z:0\"", string(stamp), "should be equal")
}

func TestJTime_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)

	jtime := JTime{}
	_ = jtime.UnmarshalJSON([]byte("\"2021-12-31T12:00:00Z:0\""))

	assert.Equal("2021-12-31T12:00:00Z:0", jtime.String(), "should be equal")
}
