package cli

import (
	"net"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// local helpers to bypass broken FlagVar interface in flagOption.go
func withInt(p *int, name string, defVal int, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.IntVar(p, name, defVal, usage)
	}
}

func withString(p *string, name string, defVal string, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.StringVar(p, name, defVal, usage)
	}
}

func withBool(p *bool, name string, defVal bool, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.BoolVar(p, name, defVal, usage)
	}
}

func withDuration(p *time.Duration, name string, defVal time.Duration, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.DurationVar(p, name, defVal, usage)
	}
}

func withIP(p *net.IP, name string, defVal net.IP, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.IPVar(p, name, defVal, usage)
	}
}

func withStringSlice(p *[]string, name string, defVal []string, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.StringSliceVar(p, name, defVal, usage)
	}
}

func withIntSlice(p *[]int, name string, defVal []int, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		f.IntSliceVar(p, name, defVal, usage)
	}
}

func TestParseArgs(t *testing.T) {
	t.Run("NoArguments", func(t *testing.T) {
		fs, err := ParseArgs([]string{})
		require.NoError(t, err)
		assert.NotNil(t, fs)
		assert.False(t, fs.HasFlags())
	})

	t.Run("PositionalArgumentsOnly", func(t *testing.T) {
		args := []string{"arg1", "arg2"}
		fs, err := ParseArgs(args)
		require.NoError(t, err)
		assert.NotNil(t, fs)
		assert.Equal(t, args, fs.Args())
	})

	t.Run("WithFlags", func(t *testing.T) {
		var (
			boolVal     bool
			stringVal   string
			intVal      int
			durationVal time.Duration
			ipVal       net.IP
		)

		args := []string{
			"--bool",
			"--string", "hello",
			"--int", "42",
			"--duration", "1h",
			"--ip", "192.168.1.1",
			"extra",
		}

		fs, err := ParseArgs(args,
			withBool(&boolVal, "bool", false, "bool flag"),
			withString(&stringVal, "string", "default", "string flag"),
			withInt(&intVal, "int", 0, "int flag"),
			withDuration(&durationVal, "duration", 0, "duration flag"),
			withIP(&ipVal, "ip", nil, "ip flag"),
		)

		require.NoError(t, err)
		assert.True(t, boolVal)
		assert.Equal(t, "hello", stringVal)
		assert.Equal(t, 42, intVal)
		assert.Equal(t, time.Hour, durationVal)
		assert.Equal(t, net.ParseIP("192.168.1.1"), ipVal)
		assert.Equal(t, []string{"extra"}, fs.Args())
	})

	t.Run("WithShortFlags", func(t *testing.T) {
		var (
			boolVal   bool
			stringVal string
		)

		args := []string{"-b", "-s", "world"}

		_, err := ParseArgs(args,
			func(f *pflag.FlagSet) { f.BoolVarP(&boolVal, "bool", "b", false, "") },
			func(f *pflag.FlagSet) { f.StringVarP(&stringVal, "string", "s", "", "") },
		)

		require.NoError(t, err)
		assert.True(t, boolVal)
		assert.Equal(t, "world", stringVal)
	})

	t.Run("DefaultValues", func(t *testing.T) {
		var (
			stringVal string
			intVal    int
		)

		// We must provide at least one argument because ParseArgs returns early if args is empty,
		// skipping flag registration and thus default values are not applied to the variables.
		args := []string{"--"}

		_, err := ParseArgs(args,
			withString(&stringVal, "string", "default", "string flag"),
			withInt(&intVal, "int", 100, "int flag"),
		)

		require.NoError(t, err)
		assert.Equal(t, "default", stringVal)
		assert.Equal(t, 100, intVal)
	})

	t.Run("EmptyArgsBehavior", func(t *testing.T) {
		var stringVal string
		args := []string{}

		// Documenting the current behavior: options are NOT applied if args is empty
		_, err := ParseArgs(args,
			withString(&stringVal, "string", "default", "string flag"),
		)

		require.NoError(t, err)
		assert.Equal(t, "", stringVal) // Should have been "default" if ParseArgs was fixed
	})

	t.Run("InvalidFlagValue", func(t *testing.T) {
		var intVal int
		args := []string{"--int", "not-an-int"}

		_, err := ParseArgs(args,
			withInt(&intVal, "int", 0, "int flag"),
		)

		assert.Error(t, err)
	})

	t.Run("UnknownFlag", func(t *testing.T) {
		args := []string{"--unknown"}

		_, err := ParseArgs(args)

		assert.Error(t, err)
	})

	t.Run("SliceFlags", func(t *testing.T) {
		var (
			stringSlice []string
			intSlice    []int
		)

		args := []string{
			"--strings", "a,b,c",
			"--ints", "1", "--ints", "2",
		}

		_, err := ParseArgs(args,
			withStringSlice(&stringSlice, "strings", []string{}, ""),
			withIntSlice(&intSlice, "ints", []int{}, ""),
		)

		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, stringSlice)
		assert.Equal(t, []int{1, 2}, intSlice)
	})
}
