package fat

import (
	"fmt"
	"github.com/metakeule/fmtdate"
	"time"
)

var (
	zeroTime, _ = time.Parse(time.UnixDate, time.UnixDate)
	timeFormat  = "DD.MM.YYYY hh:mm:ss"
)

type time_ time.Time

func Time(t time.Time) *time_    { t_ := time_(t); return &t_ }
func (t time_) Typ() string      { return "time" }
func (t time_) Get() interface{} { return time.Time(t) }
func (t time_) String() string   { return fmt.Sprintf(fmtdate.Format(timeFormat, time.Time(t))) }

func (øtime *time_) Set(i interface{}) error {
	t, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("can't convert %T to time", i)
	}
	*øtime = time_(t)
	return nil
}

func (øtime *time_) Scan(s string) error {
	t, err := fmtdate.Parse(timeFormat, s)
	if err != nil {
		return err
	}
	*øtime = time_(t)
	return nil
}
