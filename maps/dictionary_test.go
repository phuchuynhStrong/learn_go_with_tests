package maps

import (
	"testing"
)

type Dictionary map[string]string

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

const (
	ErrNotFound         = DictionaryErr("could not find the word you were looking for")
	ErrWordExists       = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (d Dictionary) Add(word, meaning string) error {
	_, err := d.Search(word)

	switch err {
	case ErrNotFound:
		d[word] = meaning
	case nil:
		return ErrWordExists
	default:
		return err
	}

	return nil
}

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}
	return definition, nil
}

func (d Dictionary) Update(word, newDefinition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		d[word] = newDefinition
	default:
		return err
	}

	return nil
}

func (d Dictionary) Delete(word string) {
	delete(d, word)
}

func TestDelete(t *testing.T) {
	word := "test"
	definition := "this is just a text"
	dict := Dictionary{word: definition}

	dict.Delete(word)

	_, err := dict.Search(word)
	if err != ErrNotFound {
		t.Errorf("Expected %q to be deleted", word)
	}
}

func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dict := Dictionary{word: definition}
		newDefinition := "new definition"

		dict.Update(word, newDefinition)

		assertDefiniation(t, dict, word, newDefinition)
	})

	t.Run("new word", func(t *testing.T) {
		word := "test"
		definition := "this is just a text"

		dict := Dictionary{}
		err := dict.Update(word, definition)

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "test"
		definition := "this is a new string"
		err := dictionary.Add(word, definition)

		assertError(t, err, nil)
		assertDefiniation(t, dictionary, word, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is a new string"
		dictionary := Dictionary{word: definition}
		err := dictionary.Add(word, "mew test")

		assertError(t, err, ErrWordExists)
		assertDefiniation(t, dictionary, word, definition)
	})
}

func TestDictionary(t *testing.T) {
	dictionary := Dictionary{"test": "this is a string"}
	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "this is a string"

		assertStrings(t, got, want)
	})

	t.Run("unknowne word", func(t *testing.T) {
		_, err := dictionary.Search("Unknown")
		want := ErrNotFound

		if err == nil {
			t.Fatal("expected to get an error.")
		}

		assertError(t, err, want)
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q but want %q", got, want)
	}
}

func assertDefiniation(t testing.TB, dictionary Dictionary, word, defination string) {
	t.Helper()

	got, err := dictionary.Search(word)

	if err != nil {
		t.Fatal("should  find added word: ", err)
	}

	if defination != got {
		t.Errorf("got %q want %q", got, defination)
	}
}
