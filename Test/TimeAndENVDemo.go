package main

import "fmt"
import (
	"time"
	"os"
	"strings"
)

func main() {
	p := fmt.Println
	now := time.Now()
	p(now)
	then := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	p(then)
	p(then.Year())
	p(then.Month())
	p(then.Day())
	p(then.Hour())
	p(then.Minute())
	p(then.Second())
	p(then.Nanosecond())
	p(then.Location())
	p(then.Weekday())
	p(then.Before(now))
	p(then.After(now))
	p(then.Equal(now))
	diff := now.Sub(then)
	p(diff)
	p(diff.Hours())
	p(diff.Minutes())
	p(diff.Seconds())
	p(diff.Nanoseconds())
	p(then.Add(diff))
	p(then.Add(-diff))

	os.Setenv("FOO", "1")
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))
	p()
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		p(pair[0]+" "+pair[1])
	}
}