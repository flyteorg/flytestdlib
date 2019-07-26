// Defines global context-aware logger.
// The default implementation uses logrus. This package registers "logger" config section on init(). The structure of the
// config section is expected to be un-marshal-able to Config struct.
package logger

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := SetConfig(&Config{
		Level:             InfoLevel,
		IncludeSourceCode: true,
	}); err != nil {
		panic(err)
	}
}

func TestIsLoggable(t *testing.T) {
	type args struct {
		ctx   context.Context
		level Level
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Debug Is not loggable", args{ctx: context.TODO(), level: DebugLevel}, false},
		{"Info Is loggable", args{ctx: context.TODO(), level: InfoLevel}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLoggable(tt.args.ctx, tt.args.level); got != tt.want {
				t.Errorf("IsLoggable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestPrint(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Print(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warn(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestWarning(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warning(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestPanic(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				Panic(tt.args.ctx, tt.args.args...)
			})
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debugf(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestPrintf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Printf(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestInfofNoCtx(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{format: "%v", args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InfofNoCtx(tt.args.format, tt.args.args...)
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warnf(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestWarningf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warningf(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Errorf(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestPanicf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				Panicf(tt.args.ctx, tt.args.format, tt.args.args...)
			})
		})
	}
}

func TestDebugln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debugln(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestPrintln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Println(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestInfoln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infoln(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestWarnln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warnln(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestWarningln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warningln(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestErrorln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Errorln(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestPanicln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{ctx: context.TODO(), args: []interface{}{"arg"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				Panicln(tt.args.ctx, tt.args.args...)
			})
		})
	}
}

func Test_getLogger(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *logrus.Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLogger(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIndent(t *testing.T) {
	type args struct {
		ctx              context.Context
		additionalIndent string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithIndent(tt.args.ctx, tt.args.additionalIndent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIndent(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIndent(tt.args.ctx); got != tt.want {
				t.Errorf("getIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatal(tt.args.ctx, tt.args.args...)
		})
	}
}

func TestFatalf(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatalf(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func TestFatalln(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatalln(tt.args.ctx, tt.args.args...)
		})
	}
}

func Test_onConfigUpdated(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			onConfigUpdated(tt.args.cfg)
		})
	}
}

func Test_getSourceLocation(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSourceLocation(); got != tt.want {
				t.Errorf("getSourceLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogWriter(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *io.PipeWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogWriter(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopLogger_WithField(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
		want *logrus.Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			if got := n.WithField(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoopLogger.WithField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopLogger_WithFields(t *testing.T) {
	type args struct {
		fields logrus.Fields
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
		want *logrus.Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			if got := n.WithFields(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoopLogger.WithFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopLogger_WithError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
		want *logrus.Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			if got := n.WithError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoopLogger.WithError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopLogger_Debugf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Debugf(tt.args.format, tt.args.args...)
		})
	}
}

func TestNoopLogger_Infof(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Infof(tt.args.format, tt.args.args...)
		})
	}
}

func TestNoopLogger_Warnf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Warnf(tt.args.format, tt.args.args...)
		})
	}
}

func TestNoopLogger_Warningf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Warningf(tt.args.format, tt.args.args...)
		})
	}
}

func TestNoopLogger_Errorf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Errorf(tt.args.format, tt.args.args...)
		})
	}
}

func TestNoopLogger_Debug(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Debug(tt.args.args...)
		})
	}
}

func TestNoopLogger_Info(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Info(tt.args.args...)
		})
	}
}

func TestNoopLogger_Warn(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Warn(tt.args.args...)
		})
	}
}

func TestNoopLogger_Warning(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Warning(tt.args.args...)
		})
	}
}

func TestNoopLogger_Error(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Error(tt.args.args...)
		})
	}
}

func TestNoopLogger_Debugln(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Debugln(tt.args.args...)
		})
	}
}

func TestNoopLogger_Infoln(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Infoln(tt.args.args...)
		})
	}
}

func TestNoopLogger_Warnln(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Warnln(tt.args.args...)
		})
	}
}

func TestNoopLogger_Warningln(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Warningln(tt.args.args...)
		})
	}
}

func TestNoopLogger_Errorln(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Errorln(tt.args.args...)
		})
	}
}

func TestNoopLogger_Print(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Print(tt.args.in0...)
		})
	}
}

func TestNoopLogger_Printf(t *testing.T) {
	type args struct {
		in0 string
		in1 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Printf(tt.args.in0, tt.args.in1...)
		})
	}
}

func TestNoopLogger_Println(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Println(tt.args.in0...)
		})
	}
}

func TestNoopLogger_Fatal(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Fatal(tt.args.in0...)
		})
	}
}

func TestNoopLogger_Fatalf(t *testing.T) {
	type args struct {
		in0 string
		in1 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Fatalf(tt.args.in0, tt.args.in1...)
		})
	}
}

func TestNoopLogger_Fatalln(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Fatalln(tt.args.in0...)
		})
	}
}

func TestNoopLogger_Panic(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Panic(tt.args.in0...)
		})
	}
}

func TestNoopLogger_Panicf(t *testing.T) {
	type args struct {
		in0 string
		in1 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Panicf(tt.args.in0, tt.args.in1...)
		})
	}
}

func TestNoopLogger_Panicln(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name string
		n    NoopLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopLogger{}
			n.Panicln(tt.args.in0...)
		})
	}
}
