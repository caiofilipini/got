package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	assert.Equal(t, "greet", Greet().Name())
}

func TestHelp(t *testing.T) {
	assert.Equal(t, "greet – shows greetings", Greet().Help())
}

func TestUsage(t *testing.T) {
	usage := Greet().Usage()

	assert.Len(t, usage, 1)
	assert.Equal(t, "greet <nickname>", usage[0])
}

func TestPattern(t *testing.T) {
	p := Greet().Pattern()

	assert.NotRegexp(t, p, "greet")
	assert.NotRegexp(t, p, "greet ")
	assert.NotRegexp(t, p, "greet       ")

	assert.Regexp(t, p, "greet marvin")
	assert.Regexp(t, p, "greet      marvin")
}

func TestPatternCapturesQuery(t *testing.T) {
	p := Greet().Pattern()
	match := p.FindStringSubmatch("greet marvin")

	assert.Len(t, match, 2)
	assert.Equal(t, "marvin", match[1])
}

func TestRun(t *testing.T) {
	result := Greet().Run("marvin")

	assert.Len(t, result, 1)
	assert.Equal(t, "ohai there, marvin!", result[0])
}
