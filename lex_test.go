package crossplane

import (
	"os"
	"path/filepath"
	"testing"
)

type tokenLine struct {
	value string
	line  int
}

type lexFixture struct {
	name   string
	tokens []tokenLine
}

var lexFixtures = []lexFixture{
	lexFixture{"simple", []tokenLine{
		tokenLine{"events", 1},
		tokenLine{"{", 1},
		tokenLine{"worker_connections", 2},
		tokenLine{"1024", 2},
		tokenLine{";", 2},
		tokenLine{"}", 3},
		tokenLine{"http", 5},
		tokenLine{"{", 5},
		tokenLine{"server", 6},
		tokenLine{"{", 6},
		tokenLine{"listen", 7},
		tokenLine{"127.0.0.1:8080", 7},
		tokenLine{";", 7},
		tokenLine{"server_name", 8},
		tokenLine{"default_server", 8},
		tokenLine{";", 8},
		tokenLine{"location", 9},
		tokenLine{"/", 9},
		tokenLine{"{", 9},
		tokenLine{"return", 10},
		tokenLine{"200", 10},
		tokenLine{"foo bar baz", 10},
		tokenLine{";", 10},
		tokenLine{"}", 11},
		tokenLine{"}", 12},
		tokenLine{"}", 13},
	}},
	lexFixture{"with-comments", []tokenLine{
		tokenLine{"events", 1},
		tokenLine{"{", 1},
		tokenLine{"worker_connections", 2},
		tokenLine{"1024", 2},
		tokenLine{";", 2},
		tokenLine{"}", 3},
		tokenLine{"#comment", 4},
		tokenLine{"http", 5},
		tokenLine{"{", 5},
		tokenLine{"server", 6},
		tokenLine{"{", 6},
		tokenLine{"listen", 7},
		tokenLine{"127.0.0.1:8080", 7},
		tokenLine{";", 7},
		tokenLine{"#listen", 7},
		tokenLine{"server_name", 8},
		tokenLine{"default_server", 8},
		tokenLine{";", 8},
		tokenLine{"location", 9},
		tokenLine{"/", 9},
		tokenLine{"{", 9},
		tokenLine{"## this is brace", 9},
		tokenLine{"# location /", 10},
		tokenLine{"return", 11},
		tokenLine{"200", 11},
		tokenLine{"foo bar baz", 11},
		tokenLine{";", 11},
		tokenLine{"}", 12},
		tokenLine{"}", 13},
		tokenLine{"}", 14},
	}},
	lexFixture{"messy", []tokenLine{
		tokenLine{"user", 1},
		tokenLine{"nobody", 1},
		tokenLine{";", 1},
		tokenLine{"# hello\\n\\\\n\\\\\\n worlddd  \\#\\\\#\\\\\\# dfsf\\n \\\\n \\\\\\n ", 2},
		tokenLine{"events", 3},
		tokenLine{"{", 3},
		tokenLine{"worker_connections", 3},
		tokenLine{"2048", 3},
		tokenLine{";", 3},
		tokenLine{"}", 3},
		tokenLine{"http", 5},
		tokenLine{"{", 5},
		tokenLine{"#forteen", 5},
		tokenLine{"# this is a comment", 6},
		tokenLine{"access_log", 7},
		tokenLine{"off", 7},
		tokenLine{";", 7},
		tokenLine{"default_type", 7},
		tokenLine{"text/plain", 7},
		tokenLine{";", 7},
		tokenLine{"error_log", 7},
		tokenLine{"off", 7},
		tokenLine{";", 7},
		tokenLine{"server", 8},
		tokenLine{"{", 8},
		tokenLine{"listen", 9},
		tokenLine{"8083", 9},
		tokenLine{";", 9},
		tokenLine{"return", 10},
		tokenLine{"200", 10},
		tokenLine{"Ser\" ' ' ver\\\\ \\ $server_addr:\\$server_port\\n\\nTime: $time_local\\n\\n", 10},
		tokenLine{";", 10},
		tokenLine{"}", 11},
		tokenLine{"server", 12},
		tokenLine{"{", 12},
		tokenLine{"listen", 12},
		tokenLine{"8080", 12},
		tokenLine{";", 12},
		tokenLine{"root", 13},
		tokenLine{"/usr/share/nginx/html", 13},
		tokenLine{";", 13},
		tokenLine{"location", 14},
		tokenLine{"~", 14},
		tokenLine{"/hello/world;", 14},
		tokenLine{"{", 14},
		tokenLine{"return", 14},
		tokenLine{"301", 14},
		tokenLine{"/status.html", 14},
		tokenLine{";", 14},
		tokenLine{"}", 14},
		tokenLine{"location", 15},
		tokenLine{"/foo", 15},
		tokenLine{"{", 15},
		tokenLine{"}", 15},
		tokenLine{"location", 15},
		tokenLine{"/bar", 15},
		tokenLine{"{", 15},
		tokenLine{"}", 15},
		tokenLine{"location", 16},
		tokenLine{"/\\{\\;\\}\\ #\\ ab", 16},
		tokenLine{"{", 16},
		tokenLine{"}", 16},
		tokenLine{"# hello", 16},
		tokenLine{"if", 17},
		tokenLine{"($request_method", 17},
		tokenLine{"=", 17},
		tokenLine{"P\\{O\\)\\###\\;ST", 17},
		tokenLine{")", 17},
		tokenLine{"{", 17},
		tokenLine{"}", 17},
		tokenLine{"location", 18},
		tokenLine{"/status.html", 18},
		tokenLine{"{", 18},
		tokenLine{"try_files", 19},
		tokenLine{"/abc/${uri} /abc/${uri}.html", 19},
		tokenLine{"=404", 19},
		tokenLine{";", 19},
		tokenLine{"}", 20},
		tokenLine{"location", 21},
		tokenLine{"/sta;\n                    tus", 21},
		tokenLine{"{", 22},
		tokenLine{"return", 22},
		tokenLine{"302", 22},
		tokenLine{"/status.html", 22},
		tokenLine{";", 22},
		tokenLine{"}", 22},
		tokenLine{"location", 23},
		tokenLine{"/upstream_conf", 23},
		tokenLine{"{", 23},
		tokenLine{"return", 23},
		tokenLine{"200", 23},
		tokenLine{"/status.html", 23},
		tokenLine{";", 23},
		tokenLine{"}", 23},
		tokenLine{"}", 23},
		tokenLine{"server", 24},
		tokenLine{"{", 25},
		tokenLine{"}", 25},
		tokenLine{"}", 25},
	}},
	lexFixture{"quote-behavior", []tokenLine{
		tokenLine{"outer-quote", 1},
		tokenLine{"left", 1},
		tokenLine{"-quote", 1},
		tokenLine{"right-\"quote\"", 1},
		tokenLine{"inner\"-\"quote", 1},
		tokenLine{";", 1},
		tokenLine{"", 2},
		tokenLine{"", 2},
		tokenLine{"left-empty", 2},
		tokenLine{"right-empty\"\"", 2},
		tokenLine{"inner\"\"empty", 2},
		tokenLine{"right-empty-single\"", 2},
		tokenLine{";", 2},
	}},
	lexFixture{"quoted-right-brace", []tokenLine{
		tokenLine{"events", 1},
		tokenLine{"{", 1},
		tokenLine{"}", 1},
		tokenLine{"http", 2},
		tokenLine{"{", 2},
		tokenLine{"log_format", 3},
		tokenLine{"main", 3},
		tokenLine{"escape=json", 3},
		tokenLine{"{ \"@timestamp\": \"$time_iso8601\", ", 4},
		tokenLine{"\"server_name\": \"$server_name\", ", 5},
		tokenLine{"\"host\": \"$host\", ", 6},
		tokenLine{"\"status\": \"$status\", ", 7},
		tokenLine{"\"request\": \"$request\", ", 8},
		tokenLine{"\"uri\": \"$uri\", ", 9},
		tokenLine{"\"args\": \"$args\", ", 10},
		tokenLine{"\"https\": \"$https\", ", 11},
		tokenLine{"\"request_method\": \"$request_method\", ", 12},
		tokenLine{"\"referer\": \"$http_referer\", ", 13},
		tokenLine{"\"agent\": \"$http_user_agent\"", 14},
		tokenLine{"}", 15},
		tokenLine{";", 15},
		tokenLine{"}", 16},
	}},
}

func TestLex(t *testing.T) {
	for _, fixture := range lexFixtures {
		t.Run(fixture.name, func(t *testing.T) {
			path := filepath.Join("testdata", fixture.name, "nginx.conf")
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			i := 0
			for token := range lex(file) {
				expected := fixture.tokens[i]
				if token.Value != expected.value || token.Line != expected.line {
					t.Fatalf("expected (%q,%d) but got (%q,%d)", expected.value, expected.line, token.Value, token.Line)
				}
				i++
			}
		})
	}
}
