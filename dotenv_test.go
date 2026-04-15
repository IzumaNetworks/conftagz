package conftagz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDotEnvBasic(t *testing.T) {
	input := "KEY1=value1\nKEY2=value2\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, "value1", m["KEY1"])
	assert.Equal(t, "value2", m["KEY2"])
}

func TestParseDotEnvComments(t *testing.T) {
	input := "# this is a comment\nKEY1=value1\n# another comment\nKEY2=value2\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, "value1", m["KEY1"])
	assert.Equal(t, "value2", m["KEY2"])
}

func TestParseDotEnvBlankLines(t *testing.T) {
	input := "\nKEY1=value1\n\n\nKEY2=value2\n   \n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Len(t, m, 2)
}

func TestParseDotEnvExportPrefix(t *testing.T) {
	input := "export KEY1=value1\nexport KEY2=value2\nKEY3=value3\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, "value1", m["KEY1"])
	assert.Equal(t, "value2", m["KEY2"])
	assert.Equal(t, "value3", m["KEY3"])
}

func TestParseDotEnvQuotedValues(t *testing.T) {
	input := "KEY1=\"double quoted\"\nKEY2='single quoted'\nKEY3=unquoted\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, "double quoted", m["KEY1"])
	assert.Equal(t, "single quoted", m["KEY2"])
	assert.Equal(t, "unquoted", m["KEY3"])
}

func TestParseDotEnvValueWithEquals(t *testing.T) {
	input := "KEY1=val=ue\nKEY2=a=b=c\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, "val=ue", m["KEY1"])
	assert.Equal(t, "a=b=c", m["KEY2"])
}

func TestParseDotEnvEmptyValue(t *testing.T) {
	input := "KEY1=\nKEY2=value2\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, "", m["KEY1"])
	assert.Equal(t, "value2", m["KEY2"])
}

func TestParseDotEnvSkipsMalformedLines(t *testing.T) {
	input := "NOEQUALS\nKEY1=value1\n"
	m, err := ParseDotEnv(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Len(t, m, 1)
	assert.Equal(t, "value1", m["KEY1"])
}

func TestLoadDotEnvFileMissing(t *testing.T) {
	m, err := LoadDotEnvFile("/nonexistent/path/.env")
	assert.NoError(t, err)
	assert.Empty(t, m)
}

func TestLoadDotEnvFileValid(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	err := os.WriteFile(envPath, []byte("KEY1=value1\nKEY2=value2\n"), 0644)
	assert.NoError(t, err)

	m, err := LoadDotEnvFile(envPath)
	assert.NoError(t, err)
	assert.Equal(t, "value1", m["KEY1"])
	assert.Equal(t, "value2", m["KEY2"])
}

func TestMergeMaps(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	override := map[string]string{"B": "override", "D": "4"}
	merged := MergeMaps(base, override)
	assert.Equal(t, "1", merged["A"])
	assert.Equal(t, "override", merged["B"])
	assert.Equal(t, "3", merged["C"])
	assert.Equal(t, "4", merged["D"])
}

func TestEnvFieldSubstitutionWithDotEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "Important=FromDotEnv\nVeryImportant=AlsoFromDotEnv\nExtremelyImportant=42\n"
	err := os.WriteFile(envPath, []byte(content), 0644)
	assert.NoError(t, err)

	mystruct := MyStruct{"", "", 0, 0, 0, 0, false}
	opts := &EnvFieldSubstOpts{
		UseDotEnv:      true,
		DotEnvFilePath: envPath,
	}
	_, err = EnvFieldSubstitution(&mystruct, opts)
	assert.NoError(t, err)
	assert.Equal(t, "FromDotEnv", mystruct.Field1)
	assert.Equal(t, "AlsoFromDotEnv", mystruct.Field2)
	assert.Equal(t, 42, mystruct.Field3)
}

func TestEnvFieldSubstitutionDotEnvOverriddenByRealEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "CONFTAGZ_TEST_DOTENV_PRIORITY=from_dotenv\n"
	err := os.WriteFile(envPath, []byte(content), 0644)
	assert.NoError(t, err)

	t.Setenv("CONFTAGZ_TEST_DOTENV_PRIORITY", "from_real_env")

	type PriorityStruct struct {
		Val string `env:"CONFTAGZ_TEST_DOTENV_PRIORITY"`
	}

	s := &PriorityStruct{}
	opts := &EnvFieldSubstOpts{
		UseDotEnv:      true,
		DotEnvFilePath: envPath,
	}
	_, err = EnvFieldSubstitution(s, opts)
	assert.NoError(t, err)
	assert.Equal(t, "from_real_env", s.Val)
}

func TestEnvFieldSubstitutionDotEnvDisabledByDefault(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "CONFTAGZ_TEST_DISABLED=should_not_appear\n"
	err := os.WriteFile(envPath, []byte(content), 0644)
	assert.NoError(t, err)

	type DisabledStruct struct {
		Val string `env:"CONFTAGZ_TEST_DISABLED"`
	}

	s := &DisabledStruct{}
	// Default opts - UseDotEnv is false
	_, err = EnvFieldSubstitution(s, nil)
	assert.NoError(t, err)
	assert.Equal(t, "", s.Val)
}

func TestEnvFieldSubstitutionDotEnvWithSlices(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "ENV_STRINGS=x,y,z\nENV_INTS=10,20,30\n"
	err := os.WriteFile(envPath, []byte(content), 0644)
	assert.NoError(t, err)

	mystruct := MyStructWithEnvSlices{}
	opts := &EnvFieldSubstOpts{
		UseDotEnv:      true,
		DotEnvFilePath: envPath,
	}
	_, err = EnvFieldSubstitution(&mystruct, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"x", "y", "z"}, mystruct.StringSlice)
	assert.Equal(t, []int{10, 20, 30}, mystruct.IntSlice)
}
