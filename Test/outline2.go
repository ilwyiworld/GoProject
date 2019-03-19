package main

import (
	"golang.org/x/net/html"
	"os"
	"fmt"
)
//通过递归的方式遍历整个 HTML 结点树，并输出树的结构
var depth int
func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	forEachNode(doc,startElement,endElement)
}

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		//%*s中的*会在字符串之前填充一些空格。在例子中,每次输出会先填充depth*2数量的空格，再输出""，最后再输出HTML标签
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}
}
func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
// forEachNode 针对每个结点 x,都会调用 pre(x)和 post(x)。
// pre 和 post 都是可选的。
// 遍历孩子结点之前,pre 被调用
// 遍历孩子结点之后， post 被调用
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
