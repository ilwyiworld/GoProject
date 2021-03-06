package main

import (
	"time"
	"text/tabwriter"
	"os"
	"fmt"
	"sort"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	// text/tabwriter包来生成一个列是整齐对齐和隔开的表格
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
	// Flush 方法会格式化整个表格并且将它写向os.Stdout（标准输出）
}

type byArtist []*Track

func (x byArtist) Len() int {
	return len(x)
}
func (x byArtist) Less(i, j int) bool {
	return x[i].Artist < x[j].Artist
}
func (x byArtist) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type byYear []*Track

func (x byYear) Len() int {
	return len(x)
}
func (x byYear) Less(i, j int) bool {
	return x[i].Year < x[j].Year
}
func (x byYear) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func main() {
	sort.Sort(byArtist(tracks))
	printTracks(tracks)
	fmt.Println()

	/*sort 包定义了一个不公开的 struct 类型 reverse，它嵌入了一个 sort.Interface。
	reverse 的 Less 方法调用了内嵌的 sort.Interface 值的 Less 方法，但是通过交换索引的方式使
	排序结果变成逆序。*/
	sort.Sort(sort.Reverse(byArtist(tracks)))
	printTracks(tracks)

	sort.Sort(byYear(tracks))
	printTracks(tracks)
	fmt.Println()
}
