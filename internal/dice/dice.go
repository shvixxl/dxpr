// Package dice provides dice rolling.
package dice

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Dice rolls dice. It is an interface so evaluation can be tested
// deterministically with a fixed source of rolls.
type Dice interface {
	// Roll returns a random integer in [1, sides].
	Roll(sides int) (int, error)
}

// cryptoDice is the default Dice implementation, using crypto/rand.
type cryptoDice struct{}

// NewDice returns the default dice roller. It uses crypto/rand and is safe
// for game randomness.
func NewDice() Dice { return cryptoDice{} }

// Roll implements Dice.
func (cryptoDice) Roll(sides int) (int, error) {
	if sides < 1 {
		return 0, fmt.Errorf("invalid dice: %d sides", sides)
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(sides)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 1, nil
}
