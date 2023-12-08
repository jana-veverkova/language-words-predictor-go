package frequencydictionary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchWord(t *testing.T) {
	require.Equal(t, "abc", matchWord("abc"))
	require.Equal(t, "hello", matchWord("Hello!"))
	require.Equal(t, "hello", matchWord("!Hello"))
	require.Equal(t, "sour-milk", matchWord("sour-milk"))
	require.Equal(t, "milk", matchWord("-milk"))
	require.Equal(t, "hello", matchWord("2hello"))
	require.Equal(t, "", matchWord("b2b"))
	require.Equal(t, "hello", matchWord("HELLO2"))
	require.Equal(t, "hello", matchWord("1.Hello..."))
	require.Equal(t, "hello", matchWord("1.hello..."))
	require.Equal(t, "comment", matchWord("\"comment\""))
	require.Equal(t, "comment", matchWord("(comment)"))
	require.Equal(t, "it’s", matchWord("it’s"))
	require.Equal(t, "", matchWord(""))
	require.Equal(t, "", matchWord("123"))
	require.Equal(t, "publishing", matchWord("publishing'"))
	require.Equal(t, "", matchWord("sey\\x"))
	require.Equal(t, "limitation", matchWord("limitation-"))
	require.Equal(t, "limitation", matchWord("limitation."))
}

func TestOrderKeys(t *testing.T) {
	testData := map[string]int{
		"a": 4,
		"b": 3,
		"c": 1,
		"d": 0,
		"e": 5,
	}

	expected := []string{"e", "a", "b", "c", "d"}

	actual := orderKeys(testData)

	require.Equal(t, expected, actual)
}
