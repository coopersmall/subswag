package card

import (
	"math/rand"
	"time"
)

func ShuffleCards(deck []SerializableCardID) []SerializableCardID {
	// Create a copy of the deck to avoid modifying the original
	shuffled := make([]SerializableCardID, len(deck))
	copy(shuffled, deck)

	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Fisher-Yates shuffle algorithm
	for i := len(shuffled) - 1; i > 0; i-- {
		// Generate random index between 0 and i
		j := r.Intn(i + 1)

		// Swap elements at i and j
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}
