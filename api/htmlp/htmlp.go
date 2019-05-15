package htmlp

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var bannedMap = map[atom.Atom]bool{
	atom.Svg:    true,
	atom.Img:    true,
	atom.Style:  true,
	atom.Script: true,
}
var whiteSpaces = regexp.MustCompile(`\s+`)

func Parse(htm string) string {
	n, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		log.Fatal(err)
	}

	node := shakeTree(n)

	return renderNode(node)
}

func shakeTree(n *html.Node) *html.Node {
	var banned []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if isBanned(c) {
			banned = append(banned, c)
		} else {
			c.Data = whiteSpaces.ReplaceAllString(strings.TrimSpace(c.Data), " ")
			shakeTree(c)
		}
	}

	for _, b := range banned {
		n.RemoveChild(b)
	}

	return n
}

func isBanned(n *html.Node) bool {
	if _, ok := bannedMap[n.DataAtom]; ok {
		return true
	}

	return false
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	_ = html.Render(w, n)
	return buf.String()

}
