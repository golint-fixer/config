package config

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupValues() {
	var args []string
	for v, _ := range values {
		args = append(args, v)
	}

	for _, arg := range args {
		delete(values, arg)
	}
}

func TestLoad(t *testing.T) {
	defer cleanupValues()

	// ARRAGE
	flagArgDef := argument{FlagName: "flagArgDef", Default: "foo"}
	envArgDef := argument{EnvName: "envArgDef", Default: "foo"}
	flagEnvArgDef := argument{FlagName: "flagEnvArgDef", EnvName: "flagEnvArg", Default: "foo"}

	config := map[string]argument{
		"flagDef":    flagArgDef,
		"envDef":     envArgDef,
		"flagEnvDef": flagEnvArgDef,
	}

	// ACT
	load(config)

	// ASSERT
	require.Len(t, values, len(config))
	for k, v := range config {
		configuredVal, ok := values[k]
		require.True(t, ok, fmt.Sprintf("%s was not loaded correctly", k))
		assert.Equal(t, v.Default, configuredVal.resolve(), fmt.Sprintf("%s was not loaded correctly", k))
	}
}

func TestLoadOntoStruct(t *testing.T) {
	defer cleanupValues()

	// ARRANGE
	config, err := parseJSON("example/config.json")
	require.NoError(t, err)

	load(config)
	require.Len(t, values, 2)

	// ACT
	var cStruct struct {
		Port     uint   `config:"port"`
		MyBucket string `config:"s3_bucket"`
	}
	err = Load(&cStruct)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, "", cStruct.MyBucket)
	assert.Equal(t, uint(8080), cStruct.Port)
}

func TestLoadFlag(t *testing.T) {
	defer cleanupValues()

	testDescription := "a marvelous flag used for fantastic things"
	for _, test := range []struct {
		Default   interface{}
		Type      string
		PtrType   interface{}
		StringVal string
	}{
		// uints
		{
			Default:   float64(123),
			Type:      "uint",
			PtrType:   new(uint64),
			StringVal: "123",
		},
		{
			Default:   float64(123),
			Type:      "uint8",
			PtrType:   new(uint64),
			StringVal: "123",
		},
		{
			Default:   float64(123),
			Type:      "uint16",
			PtrType:   new(uint64),
			StringVal: "123",
		},
		{
			Default:   float64(123),
			Type:      "uint32",
			PtrType:   new(uint64),
			StringVal: "123",
		},
		{
			Default:   float64(123),
			Type:      "uint64",
			PtrType:   new(uint64),
			StringVal: "123",
		},

		// strings
		{
			Default:   "somestring",
			Type:      "string",
			PtrType:   new(string),
			StringVal: "somestring",
		},
	} {
		flagName := "foobar" + test.Type
		a := argument{FlagName: flagName, Default: test.Default, Type: test.Type, Description: testDescription}
		prt, err := loadFlag(a)
		assert.NoError(t, err, test.Type)
		assert.NotNil(t, prt, test.Type)
		assert.IsType(t, test.PtrType, prt, test.Type)

		f := flag.Lookup(flagName)
		assert.Equal(t, flagName, f.Name, test.Type)
		assert.Equal(t, testDescription, f.Usage, test.Type)
		assert.EqualValues(t, test.StringVal, f.DefValue, test.Type)
	}
}

func TestLoadEnv(t *testing.T) {
	defer cleanupValues()

	for _, test := range []struct {
		ParsedVal interface{}
		Type      string
		EnvVal    string
	}{
		// uints
		{
			ParsedVal: uint64(123),
			Type:      "uint",
			EnvVal:    "123",
		},
		{
			ParsedVal: uint64(123),
			Type:      "uint8",
			EnvVal:    "123",
		},
		{
			ParsedVal: uint64(123),
			Type:      "uint16",
			EnvVal:    "123",
		},
		{
			ParsedVal: uint64(123),
			Type:      "uint32",
			EnvVal:    "123",
		},
		{
			ParsedVal: uint64(123),
			Type:      "uint64",
			EnvVal:    "123",
		},

		// strings
		{
			ParsedVal: "somestring",
			Type:      "string",
			EnvVal:    "somestring",
		},
	} {
		envName := "foobar" + test.Type
		os.Setenv(envName, test.EnvVal)
		a := argument{EnvName: envName, Type: test.Type}
		v, err := loadEnv(a)
		assert.NoError(t, err, test.Type)
		assert.Equal(t, test.ParsedVal, v, test.Type)
	}
}
