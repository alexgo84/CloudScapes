package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {

	type args struct {
		level  LogLevel
		msg    string
		fields []zap.Field
	}
	tests := []struct {
		name             string
		dev              bool
		args             args
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name: "dev logger show debug",
			dev:  true,
			args: args{
				DEBUG,
				"this is a debug log on dev logger",
				nil,
			},
			shouldContain: []string{
				"this is a debug log on dev logger",
			},
			shouldNotContain: []string{},
		},
		{
			name: "prod logger not show debug",
			dev:  false,
			args: args{
				DEBUG,
				"this is a debug log on prod logger",
				nil,
			},
			shouldContain: []string{},
			shouldNotContain: []string{
				"this is a debug log on prod logger",
			},
		},
		{
			name: "logger should hold key and value of a string field",
			dev:  false,
			args: args{
				INFO,
				"",
				[]zap.Field{Str("string key", "string value")},
			},
			shouldContain:    []string{`{"string key": "string value"}`},
			shouldNotContain: []string{},
		},
		{
			name: "logger should hold key and value of a int64 field",
			dev:  false,
			args: args{
				INFO,
				"",
				[]zap.Field{Int64("int64 key", int64(42))},
			},
			shouldContain:    []string{`{"int64 key": 42}`},
			shouldNotContain: []string{},
		},
		{
			name: "logger should hold key and value of multiple fields",
			dev:  false,
			args: args{
				INFO,
				"",
				[]zap.Field{Int64("int64 key1", int64(42)), Int64("int64 key2", int64(24)), Str("str key", "str value")},
			},
			shouldContain:    []string{`"int64 key1": 42`, `"int64 key2": 24`, `"str key": "str value"`},
			shouldNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitLogger(tt.dev, &tt.name); err != nil {
				t.Fatalf(err.Error())
			}
			Log(tt.args.level, tt.args.msg, tt.args.fields...)

			// open the log file
			path, err := getLogDirPath()
			if err != nil {
				t.Fatalf(err.Error())
			}

			pathToFile := fmt.Sprintf("%s/%s", path, tt.name)

			if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
				t.Fatalf(err.Error())
			}

			file, err := os.Open(pathToFile)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var lastLine string
			for scanner.Scan() {
				lastLine = scanner.Text()
			}

			for _, shouldContain := range tt.shouldContain {
				if !strings.Contains(lastLine, shouldContain) {
					t.Fatalf("expected log file to contain string '%s'", shouldContain)
				}
			}

			for _, shouldNotContain := range tt.shouldNotContain {
				if strings.Contains(lastLine, shouldNotContain) {
					t.Fatalf("expected log file to not contain string '%s'", shouldNotContain)
				}
			}
		})
	}
}
