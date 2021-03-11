package document

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"gitlab.com/256/DebateFrame/client/log"
)

func postProcess(htmlStr string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		panic("failed to convert doc to goquery doc")
	}
	doc.Find("b").Each(func(i int, s *goquery.Selection) {
		h, err := s.Html()
		if err != nil {
			log.PanicMessage("Failed to convert element to html", err)
		}
		s.ReplaceWithHtml("<p class=\"uk-text-bold\">" + h + "</p>")
	})
	doc.Find("span").Each(func(i int, sTemp *goquery.Selection) {
		s := sTemp.Parent()
		if len(s.Nodes) > 0 && s.Nodes[0].Type == html.ElementNode {
			if strings.ToLower(s.Nodes[0].Data) == "p" {
				size := getSize(sTemp)
				tag := fontSizeToTag(size)
				s.ReplaceWithHtml(fmt.Sprintf("<%s class=\"display: inline; font-size: %vpx;\">%s</%s>", tag, size, s.Text(), tag))
			}
		}
	})
	newHTML, err := doc.Html()
	return newHTML
}

var rp = regexp.MustCompile(`font-size:\s(\d+)px;`)

func getSize(sel *goquery.Selection) int {
	style, ok := sel.Attr("style")
	if !ok {
		log.PanicMessage("Span had no style attribute", nil)
	}
	fontSize, err := strconv.Atoi(rp.FindStringSubmatch(style)[0])
	if err != nil {
		log.PanicMessage("Font size was not number!", err)
	}
	return fontSize
}

type possText struct {
	tag      string
	fontSize int
}

var possTextSlice []possText

func init() {
	possTextSlice = []possText{
		{tag: "H1", fontSize: 42},
		{tag: "H2", fontSize: 32},
		{tag: "H3", fontSize: 24},
		{tag: "H4", fontSize: 20},
		{tag: "H6", fontSize: 14},
		{tag: "p", fontSize: 16},
	}
}

func fontSizeToTag(size int) string {
	var bestTag possText
	bestDiff := 999.0
	for i, tag := range possTextSlice {
		if i == 0 {
			bestTag = tag
		} else {
			dif := diff(float64(tag.fontSize), float64(size))
			if dif < bestDiff {
				bestTag = tag
				bestDiff = dif
			}
		}
	}
	return bestTag.tag
}

func diff(a float64, b float64) float64 {
	return math.Abs(a - b)
}
