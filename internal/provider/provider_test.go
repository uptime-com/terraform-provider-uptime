package provider

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		err = godotenv.Load("../../.env")
	}
	os.Exit(m.Run())
}

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"uptime": providerserver.NewProtocol6WithError(VersionFactory("test")()),
	}
}

func testAccAPIClient(t *testing.T) upapi.API {
	t.Helper()

	token := os.Getenv("UPTIME_TOKEN")
	require.NotEmpty(t, token, "UPTIME_TOKEN must be set for acceptance tests")

	api, err := upapi.New(upapi.WithToken(token), upapi.WithRateLimit(0.2))
	require.NoError(t, err, "failed to initialize uptime.com api client")

	return api
}

//go:embed testdata/*
var testdata embed.FS

func testSnippet(t *testing.T, fn string, section int) string {
	t.Helper()

	f, err := testdata.Open("testdata/" + fn)
	if err != nil {
		t.Fatal(err)
	}
	s := bufio.NewScanner(f)
	b := strings.Builder{}
	n := 0
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "// ---") {
			if n == section {
				return b.String()
			}
			b.Reset()
			n++
		}
		b.WriteString(s.Text())
		b.WriteString("\n")
	}
	if n != section {
		t.Fatalf("section not found: file=%s section=%d", fn, section)
	}
	return b.String()
}

func testRenderSnippet(t *testing.T, fn string, section int, data any) string {
	t.Helper()

	b := bytes.NewBuffer(nil)

	tmpl, err := template.New("").
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"petname": func(len reflect.Value, sep reflect.Value) (reflect.Value, error) {
				if len.Kind() != reflect.Int {
					return reflect.ValueOf(nil), fmt.Errorf("petname: length must be an int")
				}
				if sep.Kind() != reflect.String {
					return reflect.ValueOf(nil), fmt.Errorf("petname: separator must be a string")
				}
				return reflect.ValueOf(petname.Generate(int(len.Int()), sep.String())), nil
			},
		}).
		Parse(testSnippet(t, fn, section))
	if err != nil {
		t.Fatal(err)
	}
	err = tmpl.Execute(b, data)
	if err != nil {
		t.Fatal(err)
	}

	return b.String()
}

func TestTestSnippet(t *testing.T) {
	assert.Equal(t, "", testSnippet(t, "test_snippets.txt", 0))
	assert.Contains(t, testSnippet(t, "test_snippets.txt", 1), "Chunk 1")
	assert.Contains(t, testSnippet(t, "test_snippets.txt", 2), "Chunk 2")
}

func TestTestRenderSnippet(t *testing.T) {
	assert.Contains(t, testRenderSnippet(t, "test_snippets.txt", 1, map[string]string{"name": "foo"}), "Chunk 1 foo")
	assert.Contains(t, testRenderSnippet(t, "test_snippets.txt", 1, map[string]string{"name": "bar"}), "Chunk 1 bar")
}
