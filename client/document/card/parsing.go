package card

import (
	"regexp"
	"strconv"
	"strings"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/montanaflynn/stats"

	"gitlab.com/256/DebateFrame/client/log"
)

// GetCards gets the cards from a document
func GetCards(doc *goquery.Document) []*Card {
	// Get most frequent header
	headerTags := []float64{}
	var header = regexp.MustCompile(`H\d`)
	doc.Find("*").Each(func(_ int, sel *goquery.Selection) {
		tag := strings.ToUpper(goquery.NodeName(sel))
		if header.MatchString(tag) {
			headerTags = append(headerTags, float64(getHeaderLevel(tag)))
		}
	})
	if len(headerTags) == 0 {
		return []*Card{}
	}
	mostFreqList, err := stats.Mode(headerTags)
	if err != nil {
		log.PanicMessage("could not calculate mode of H levels", err)
	}
	mostFreq := uint8(mostFreqList[0])
	log.DebugMessage(fmt.Sprintf("Most frequent H tag: %v", mostFreq))

	sections := getCardSections(mostFreq, doc.Find("*"))

	return getCardsFromSections(sections)
}

func getHeaderLevel(tag string) uint8 {
	numStr := strings.Replace(strings.ToUpper(tag), "H", "", -1)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		panic("failed to convert number from header to int")
	}
	return uint8(num)
}

func getCardSections(hlev uint8, children *goquery.Selection) [][]*goquery.Selection {
	tagStr := fmt.Sprintf("H%v", hlev)

	var sections [][]*goquery.Selection

	var currentSection = -1
	children.Each(func(_ int, sel *goquery.Selection) {
		tag := strings.ToUpper(goquery.NodeName(sel))
		if tag == tagStr {
			currentSection++
			sections = append(sections, []*goquery.Selection{sel})
		} else if currentSection != -1 {
			sections[currentSection] = append(sections[currentSection], sel)
		}
	})
	log.DebugMessage("Found %v sections", len(sections))
	return sections
}

func getCardsFromSections(sections [][]*goquery.Selection) (cards []*Card) {
	for _, section := range sections {
		card := Card{}
		card.Title = section[0].Text()
		var text string
		for _, snippit := range section[1:] {
			text = text + "\n" + snippit.Text()
		}
		if len(section) >= 2 {
			startLine := section[1].Text()
			author, year, _ := getAuthorAndYear(startLine)
			card.Author = author
			card.Year = year
		}
		card.Contents = text
		if strictCards && (len(card.Contents) == 0 || len(card.Author) == 0 || card.Year == 0 || len(card.Title) == 0) {
			log.DebugMessage(fmt.Sprintf("Card with title \"%v\" was blocked by strictCards setting", card.Title))
		} else {
			cards = append(cards, &card)
		}
	}
	log.DebugMessage("Found %v cards", len(cards))
	return
}

var rp = regexp.MustCompile(`([A-Z]+\w+|[A-Z]+\w+\s*&\s*[A-Z]+\w+|[A-Z]+\w+\s+et\s+al|[A-Z]+\w+\s+and\s+[A-Z]+\w+),?\s+‘?'?(\d{4}|\d{1,2}|\d{1,2}[-,/]\d{1,2}[-,/]\d{1,4})[\s,\,]*`)

// Returns the lastname of the author, and the last two digits of the year
func getAuthorAndYear(str string) (string, uint8, error) {
	// Get's the author and year
	res := rp.FindStringSubmatch(str)
	if len(res) >= 3 {
		year := dateToYear(res[2])
		return res[1], year, nil
	}
	return "", 0, fmt.Errorf("could not find author and year")
}

func dateToYear(str string) uint8 {
	rp := regexp.MustCompile(`(?:\d{1,2}[\/,-]\d{1,2}[\/,-])?\d{0,}?(\d{1,2})`)
	matches := rp.FindStringSubmatch(str)
	if len(matches) >= 2 {
		year, err := strconv.Atoi(matches[len(matches)-1])
		if err != nil {
			panic("Failed to convert year to string")
		}
		return uint8(year)
	}
	log.PanicMessage("dateToYear called on non-year!", fmt.Errorf("Bad string: %s", str))
	return 0
}
