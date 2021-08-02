package maps

import "testing"

type Dictionary map[string]string

func Search(dictionary map[string]string, word string) string {
	return dictionary[word]
}

func (dictionary Dictionary) Search(word string) string {
	return dictionary[word]
}

func TestDictionary(t *testing.T) {
	dictionary := Dictionary{"test": "this is a string"}
	got := dictionary.Search("test")
	want := "this is a string"

	assertStrings(t, got, want)
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q but want %q", got, want)
	}
}
