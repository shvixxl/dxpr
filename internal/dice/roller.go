// Package dice provides dice rolling.
package dice

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Roller rolls dice. It is an interface so evaluation can be tested
// deterministically with a fixed source of rolls.
type Roller interface {
	// Roll returns a random integer in [1, sides].
	Roll(sides int) (int, error)
}

// CryptoRoller is the default Roller, using crypto/rand.
type CryptoRoller struct{}

// NewRoller returns the default dice roller.
func NewRoller() Roller { return CryptoRoller{} }

// Roll implements Roller.
func (CryptoRoller) Roll(sides int) (int, error) {
	if sides < 1 {
		return 0, fmt.Errorf("invalid dice: %d sides", sides)
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(sides)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 1, nil
}
