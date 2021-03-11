package card

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/256/DebateFrame/client/log"
)

var mux *sync.Mutex

func init() {
	mux = &sync.Mutex{}
}

// Filter hides cards based on the input string and relevance
func Filter(cards []*Card, query string) {
	mux.Lock()
	for _, card := range cards {
		if card.Element == nil {
			log.PanicMessage("Given card has no associated element! Please associate an element with the card first!", nil)
		} else {
			card.Element.ClassList().Add("nonmatch")
		}
	}
	cardOrder := search(cards, query)
	for i, card := range cardOrder {
		if card.Element == nil {
			log.PanicMessage("Given card has no associated element! Please associate an element with the card first!", nil)
		} else {
			card.Element.ClassList().Remove("nonmatch")
			card.Element.Style().Set("order", i)
		}
	}
	mux.Unlock()
}

// The threshold in points for the card to accepted
const showThreshold = 30

// If the showThreshold isn't met for at least X cards, then pick the next highest point cards until the limit is satisfied
const minCards = 4

// search returns a slice with the correct order of cards that should be shown
func search(cards []*Card, str string) []*Card {
	log.DebugMessage("Searching cards with query %s", str)
	pointLevels := []int{}
	pointsCards := make(map[int][]*Card)
	for _, card := range cards {
		val := getValue(str, card)

		if !contains(pointLevels, val) {
			pointLevels = append(pointLevels, int(val))
		}

		if pointsCards[val] == nil {
			pointsCards[val] = []*Card{}
		}
		pointsCards[val] = append(pointsCards[val], card)
	}
	matchedCards := []*Card{}

	sort.Sort(sort.Reverse(sort.IntSlice(pointLevels)))
	for _, level := range pointLevels {
		subSlice := pointsCards[level]
		if len(subSlice) != 0 {
			for _, card := range subSlice {
				if (len(matchedCards) < minCards) || level >= showThreshold {
					matchedCards = append(matchedCards, card)
				}
			}
		}
	}

	return matchedCards
}

func getValue(query string, card *Card) int {
	var points int
	if hasSharedWord(card.Author, query) {
		points += 15
	}
	if hasSharedWord(strconv.Itoa(int(card.Year)), query) {
		points += 4
	}
	points += sharedWordCount(query, card.Title) * 3
	points += sharedWordCount(query, card.Contents)
	return points
}

func hasSharedWord(one string, two string) bool {
	if sharedWordCount(one, two) >= 1 {
		return true
	}
	return false
}

func sharedWordCount(one string, two string) int {
	var count int
	splitOne := strings.Split(strings.ToLower(one), " ")
	splitTwo := strings.Split(strings.ToLower(two), " ")
	for _, oneWord := range splitOne {
		for _, twoWord := range splitTwo {
			if oneWord == twoWord {
				count++
			}
		}
	}
	return count
}
