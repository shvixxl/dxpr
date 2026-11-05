// Package dice provides dice rolling.
package dice

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Resource limits for dice evaluation. They keep adversarial input such as
// 999999999d6 from consuming unbounded memory and time. The parser enforces
// them when resolving constant dice operands and building pools; the
// evaluator enforces them when rolling.
const (
	// MaxCount is the largest number of dice in a single roll or pool.
	MaxCount = 1000
	// MaxSides is the largest number of sides on a single die.
	MaxSides = 10000
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
