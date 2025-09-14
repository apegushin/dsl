package printer

import (
	"fmt"
	"time"
)

type Printer struct {
	Message       string
	WithTimestamp bool
	RepeatTimes   int
}

func (p Printer) Print() {
	var curTimeStamp string
	if p.WithTimestamp {
		curTimeStamp = time.Now().Format(time.RFC822)
	}
	for range p.RepeatTimes {
		fmt.Println(curTimeStamp, p.Message)
	}
}
