package broker

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_logger_Info(t *testing.T) {
	type fields struct {
		log logrus.FieldLogger
	}
	type args struct {
		msg           string
		keysAndValues []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				log: tt.fields.log,
			}
			l.Info(tt.args.msg, tt.args.keysAndValues...)
		})
	}
}

type fakeLogger struct{}

func (fakeLogger) WithField(string, interface{}) *logrus.Entry { return nil }
func (fakeLogger) WithFields(logrus.Fields) *logrus.Entry      { return nil }

func (fakeLogger) WithError(error) *logrus.Entry   { return nil }
func (fakeLogger) Debugf(string, ...interface{})   {}
func (fakeLogger) Infof(string, ...interface{})    {}
func (fakeLogger) Printf(string, ...interface{})   {}
func (fakeLogger) Warnf(string, ...interface{})    {}
func (fakeLogger) Warningf(string, ...interface{}) {}
func (fakeLogger) Errorf(string, ...interface{})   {}
func (fakeLogger) Fatalf(string, ...interface{})   {}
func (fakeLogger) Panicf(string, ...interface{})   {}

func (fakeLogger) Debug(...interface{})     {}
func (fakeLogger) Info(...interface{})      {}
func (fakeLogger) Print(...interface{})     {}
func (fakeLogger) Warn(...interface{})      {}
func (fakeLogger) Warning(...interface{})   {}
func (fakeLogger) Error(...interface{})     {}
func (fakeLogger) Fatal(...interface{})     {}
func (fakeLogger) Panic(...interface{})     {}
func (fakeLogger) Debugln(...interface{})   {}
func (fakeLogger) Infoln(...interface{})    {}
func (fakeLogger) Println(...interface{})   {}
func (fakeLogger) Warnln(...interface{})    {}
func (fakeLogger) Warningln(...interface{}) {}
func (fakeLogger) Errorln(...interface{})   {}
func (fakeLogger) Fatalln(...interface{})   {}
func (fakeLogger) Panicln(...interface{})   {}
