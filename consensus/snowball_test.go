package consensus

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewSnowBall(t *testing.T) {
	expect := SnowBall{
		decisionThreshold:    20,
		preference:           Preference(3),
		decided:              false,
		consecutiveSuccesses: map[Preference]int{},
	}
	actual := NewSnowBall(3, 20)
	assert.Equal(t, expect, actual)

}

func TestSnowBall(t *testing.T) {
	var Pizza, BBQ, Chicken Preference = 1, 2, 3
	snowball := NewSnowBall(Pizza, 3)

	assert.Equal(t, Pizza, snowball.Preference())
	assert.Equal(t, false, snowball.Decided())

	snowball.SuccessPool(BBQ)
	assert.Equal(t, BBQ, snowball.Preference())
	assert.Equal(t, false, snowball.Decided())

	snowball.SuccessPool(Chicken)
	assert.Equal(t, Chicken, snowball.Preference())
	assert.Equal(t, false, snowball.Decided())

	snowball.SuccessPool(Pizza)
	assert.Equal(t, Pizza, snowball.Preference())
	assert.Equal(t, false, snowball.Decided())

	snowball.SuccessPool(Pizza)
	assert.Equal(t, Pizza, snowball.Preference())
	assert.Equal(t, false, snowball.Decided())

	snowball.SuccessPool(Pizza)
	assert.Equal(t, Pizza, snowball.Preference())
	assert.Equal(t, false, snowball.Decided())

	snowball.SuccessPool(Pizza)
	assert.Equal(t, Pizza, snowball.Preference())
	assert.Equal(t, true, snowball.Decided())
}
