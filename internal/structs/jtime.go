package structs

import (
	"fmt"
	"time"
)

const tformat = "2006-01-02T15:04:05Z07:0"

type JTime time.Time

func (t JTime) MarshalJSON() ([]byte, error) {
	//do your serializing here

	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(tformat))
	return []byte(stamp), nil
}

func (t *JTime) UnmarshalJSON(data []byte) error {
	//do your serializing here
	stamp, err := time.Parse(fmt.Sprintf("\"%s\"", tformat), string(data))
	if err != nil {
		return err
	}
	*t = JTime(stamp)
	return nil
}

func (t JTime) String() string {
	return time.Time(t).Format(tformat)
}
