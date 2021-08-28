package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func walk(x interface{}, fn func(input string)) {
	val := getValue(x)

	walkValue := func(value reflect.Value) {
		walk(value.Interface(), fn)
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	case reflect.Chan:
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			walkValue(v)
		}
	case reflect.Func:
		valFnResult := val.Call(nil)
		for _, res := range valFnResult {
			walkValue(res)
		}
	}

}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},
		{
			"Struct with two string fields",
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},
		{
			"Struct with non string field",
			struct {
				Name string
				Age  int
			}{"Chris", 20},
			[]string{"Chris"},
		},
		{
			"Struct with nested fields",
			Person{
				"Chris",
				Profile{20, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"Pointers to things",
			&Person{
				"Chris",
				Profile{20, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"Slices",
			[]Profile{
				{33, "Birmingham"},
				{20, "London"},
			},
			[]string{"Birmingham", "London"},
		},
		{
			"Arrays",
			[2]Profile{
				{33, "Birmingham"},
				{20, "London"},
			},
			[]string{"Birmingham", "London"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})

		t.Run("with maps", func(t *testing.T) {
			aMap := map[string]string{
				"Foo": "Bar",
				"Baz": "Boz",
			}
			var got []string
			walk(aMap, func(input string) {
				got = append(got, input)
			})

			assertContains(t, got, "Bar")
			assertContains(t, got, "Boz")
		})

		t.Run("with channels", func(t *testing.T) {
			aChannel := make(chan Profile)

			go func() {
				aChannel <- Profile{33, "Berlin"}
				aChannel <- Profile{33, "Munich"}
				close(aChannel)
			}()

			var got []string
			want := []string{"Berlin", "Munich"}

			walk(aChannel, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})

		t.Run("with functions", func(t *testing.T) {
			aFunction := func() (Profile, Profile) {
				return Profile{33, "Berlin"}, Profile{20, "Munich"}
			}

			var got []string
			want := []string{"Berlin", "Munich"}

			walk(aFunction, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func assertContains(t testing.TB, got []string, val string) {
	t.Helper()
	contains := false
	for _, x := range got {
		if x == val {
			contains = true
		}
	}

	if !contains {
		t.Errorf("expected %+v to contains %q but it didn't", got, val)
	}
}
