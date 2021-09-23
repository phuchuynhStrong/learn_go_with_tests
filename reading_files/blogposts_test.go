package blogposts_test

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"

	blogposts "github.com/phuchuynhStrong/blogposts"
)

type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

func TestNewBlogPost(t *testing.T) {

	t.Run("Test with description", func(t *testing.T) {
		const (
			firstBody = `Title: Post 1
	Description: Description 1
	Tags: tdd, go
	---
	Hello
	World`
			secondBody = `Title: Post 2
	Description: Description 2
	Tags: rust, borrow-checker
	---
	B
	L
	M`
		)

		fs := fstest.MapFS{
			"hello_world.md":  {Data: []byte(firstBody)},
			"hello_world2.md": {Data: []byte(secondBody)},
		}

		posts, _ := blogposts.NewPostsFromFS(fs)

		assertPost(t, posts[0], blogposts.Post{
			Title:       "Post 1",
			Description: "Description 1",
			Tags:        []string{"tdd", "go"},
			Body: `Hello
	World`,
		})
	})
}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
