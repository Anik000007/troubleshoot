package util

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HomeDir(t *testing.T) {
	tests := []struct {
		name string
		env  string
		dir  string
		want string
	}{
		{
			name: "test linux/unix home directory",
			env:  "HOME",
			dir:  "/home/test",
			want: "/home/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.dir)
			if HomeDir() != tt.want {
				assert.Equal(t, tt.want, HomeDir())
			}
			os.Unsetenv(tt.env) // cleanup
		})
	}
}

func Test_IsURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid url with http",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "valid url with https",
			url:  "https://example.com",
			want: true,
		},
		{
			name: "invalid url without scheme",
			url:  "example.com",
			want: false,
		},
		{
			name: "invalid url with spaces",
			url:  "http://example .com",
			want: false,
		},
		{
			name: "empty string",
			url:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.url); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_AppName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "simple name",
			in:   "myapp",
			want: "Myapp",
		},
		{
			name: "hyphenated name",
			in:   "my-app",
			want: "My App",
		},
		{
			name: "name with ai",
			in:   "my-ai-app",
			want: "My AI App",
		},
		{
			name: "name with io",
			in:   "my-app-io",
			want: "My App.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppName(tt.in); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_SplitYAML(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want []string
	}{
		{
			name: "single document",
			doc:  `apiVersion: v1`,
			want: []string{`apiVersion: v1`},
		},
		{
			name: "multiple documents",
			doc: `apiVersion: v1
---
apiVersion: v1`,
			want: []string{`apiVersion: v1`, `apiVersion: v1`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitYAML(tt.doc); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_EstimateNumberOfLines(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "no line",
			text: "",
			want: 0,
		},
		{
			name: "single line without newline",
			text: "Hello, World!",
			want: 1,
		},
		{
			name: "single line with newline",
			text: "Hello, World!\n",
			want: 1,
		},
		{
			name: "multiple lines",
			text: "Hello,\nWorld!",
			want: 2,
		},
		{
			name: "multiple lines ending with newline",
			text: "Hello,\nWorld!\n",
			want: 2,
		},
		{
			name: "multiple lines with extra newlines",
			text: "\nHello,\nWorld!\n\n",
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EstimateNumberOfLines(tt.text); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	tests := []struct {
		name   string
		target []string
		src    []string
		want   []string
	}{
		{
			name:   "empty target",
			target: []string{},
			src:    []string{"a", "b", "c"},
			want:   []string{"a", "b", "c"},
		},
		{
			name:   "empty src",
			target: []string{"a", "b", "c"},
			src:    []string{},
			want:   []string{"a", "b", "c"},
		},
		{
			name:   "non-empty target and src",
			target: []string{"a", "b", "c"},
			src:    []string{"d", "e", "f"},
			want:   []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			name:   "nil target",
			target: nil,
			src:    []string{"a", "b", "c"},
			want:   []string{"a", "b", "c"},
		},
		{
			name:   "nil src",
			target: []string{"a", "b", "c"},
			src:    nil,
			want:   []string{"a", "b", "c"},
		},
		{
			name:   "nil target and empty src",
			target: nil,
			src:    []string{},
			want:   []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Append(tt.target, tt.src)
			assert.Equal(t, tt.want, got, "Append() = %v, want %v", got, tt.want)
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name    string
		tpl     string
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "empty template and data",
			tpl:     "",
			data:    nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "empty template with data",
			tpl:     "",
			data:    map[string]string{"Name": "World"},
			want:    "",
			wantErr: false,
		},
		{
			name:    "empty data with template with no keys",
			tpl:     "Hello, World!",
			data:    nil,
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "simple template",
			tpl:     "Hello, {{ .Name }}!",
			data:    map[string]string{"Name": "World"},
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "template with missing key",
			tpl:     "Hello, {{ .Name }}!",
			data:    map[string]string{"Name2": "World"},
			want:    "Hello, <no value>!",
			wantErr: false,
		},
		{
			name:    "template with sprig function works",
			tpl:     "{{ \"hello \" | upper }}{{ .Name }}",
			data:    map[string]string{"Name": "World"},
			want:    "HELLO World",
			wantErr: false,
		},
		{
			name:    "template with undefined sprig function errors",
			tpl:     "{{ \"hello \" | upp }}{{ .Name }}",
			data:    map[string]string{"Name": "World"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTemplate(tt.tpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "RenderTemplate() = %v, want %v", got, tt.want)
		})
	}
}
