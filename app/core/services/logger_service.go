package services

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"scrm-sync-data/app/config"
	"time"
)

type LoggerService struct {
	context     context.Context
	cancel      context.CancelFunc
	config      *config.Config //日志配置
	*zap.Logger                //日志实例
	State       bool           //服务状态
}

const RotateByTimestamp = 0 //自定义时间分片
const RotateByDate = 1      //日分片
const RotateByHour = 2      //1小时分片
const RotateByMinute = 3    //1分钟分片

type KVEncoder struct {
	*zapcore.MapObjectEncoder
}

func (enc KVEncoder) Clone() zapcore.Encoder {
	return KVEncoder{
		zapcore.NewMapObjectEncoder(),
	}
}

func (enc KVEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	kvEnc := enc.Clone().(KVEncoder)
	buf := buffer.NewPool().Get()
	buf.AppendString(fmt.Sprintf(`time="%s"`, entry.Time.Format("2006-01-02 15:04:05")))
	buf.AppendString(" ")
	buf.AppendString(fmt.Sprintf(`level="%s"`, entry.Level.String()))
	buf.AppendString(" ")
	buf.AppendString(fmt.Sprintf(`msg="%s"`, entry.Message))
	buf.AppendString(" ")
	fieldBuf := new(buffer.Buffer)
	l := len(fields)
	for i, field := range fields {
		field.AddTo(kvEnc)
		value := kvEnc.MapObjectEncoder.Fields[field.Key]
		fieldBuf.AppendString(fmt.Sprintf(`"%s"`, field.Key))
		fieldBuf.AppendString(":")
		if value == "" {
			fieldBuf.AppendString("''")
		} else {
			switch value.(type) {
			case string:
				fieldBuf.AppendString(fmt.Sprintf(`"%v"`, value))
			default:
				fieldBuf.AppendString(fmt.Sprintf("%v", value))
			}
		}
		if i < l-1 {
			fieldBuf.AppendString(",")
		}
	}
	if len(fields) > 0 {
		buf.AppendString(fmt.Sprintf(`field="{%s}"`, fieldBuf.String()))
		buf.AppendString(" ")
	}
	if entry.Stack != "" {
		buf.AppendString(`stack="`)
		buf.AppendByte('\n')
		buf.AppendString(entry.Stack)
		buf.AppendString(`" `)
	}
	buf.AppendString(fmt.Sprintf(`file="%s"`, entry.Caller.String()))
	buf.AppendString(" ")
	buf.AppendByte('\n')
	return buf, nil
}

func NewKVEncoder() *KVEncoder {
	return &KVEncoder{}
}

func init() {
	once.Do(func() {
		err := zap.RegisterEncoder("kv", func(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return KVEncoder{
				zapcore.NewMapObjectEncoder(),
			}, nil
		})
		if err != nil {
			panic(err)
		}
	})
}

func (loggerService *LoggerService) Init(parentContext context.Context, config *config.Config) {
	loggerService.config = config
	loggerService.State = false
	loggerService.context, loggerService.cancel = context.WithCancel(parentContext)
	loggerService.Start()
}

func (loggerService *LoggerService) Start() {
	fmt.Println("start logger service...", time.Now())
	var err error
	var loggerWriteSyncer *rotatelogs.RotateLogs
	logLevel := zap.NewAtomicLevelAt(zapcore.Level(loggerService.config.LoggerServiceConfig.Level))
	switch loggerService.config.LoggerServiceConfig.RotateTimeLevel {
	case RotateByTimestamp:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(loggerService.config.LoggerServiceConfig.RotateTime)*time.Minute))
	case RotateByDate:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName+".%Y%m%d",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour))
	case RotateByHour:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName+".%Y%m%d%H",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Hour))
	case RotateByMinute:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+loggerService.config.LoggerServiceConfig.FileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Minute))
	}
	if err != nil {
		panic(err)
	}
	stackLevel := zapcore.ErrorLevel
	if loggerService.config.Development {
		stackLevel = zapcore.WarnLevel
	}
	encoder := NewKVEncoder()
	logger := zap.New(zapcore.NewCore(encoder, zapcore.AddSync(loggerWriteSyncer), logLevel), zap.AddCaller(), zap.AddStacktrace(stackLevel))
	loggerService.Logger = logger
	loggerService.State = true
}

func (loggerService *LoggerService) Stop() {
	fmt.Println("stop logger service...", time.Now())
	loggerService.cancel()
	loggerService.State = false
}

func (loggerService *LoggerService) NewLogger(fileName string) *LoggerService {
	newLoggerService := *loggerService
	var err error
	var loggerWriteSyncer *rotatelogs.RotateLogs
	logLevel := zap.NewAtomicLevelAt(zapcore.Level(loggerService.config.LoggerServiceConfig.Level))
	switch loggerService.config.LoggerServiceConfig.RotateTimeLevel {
	case RotateByTimestamp:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+fileName+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+fileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(loggerService.config.LoggerServiceConfig.RotateTime)*time.Minute))
	case RotateByDate:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+fileName+".%Y%m%d",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+fileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour))
	case RotateByHour:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+fileName+".%Y%m%d%H",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+fileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Hour))
	case RotateByMinute:
		loggerWriteSyncer, err = rotatelogs.New(
			loggerService.config.LoggerServiceConfig.OutputPath+fileName+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(loggerService.config.LoggerServiceConfig.OutputPath+fileName),
			rotatelogs.WithMaxAge(time.Duration(loggerService.config.LoggerServiceConfig.MaxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Minute))
	}
	if err != nil {
		panic(err)
	}
	stackLevel := zapcore.ErrorLevel
	if loggerService.config.Development {
		stackLevel = zapcore.WarnLevel
	}
	encoder := NewKVEncoder()
	logger := zap.New(zapcore.NewCore(encoder, zapcore.AddSync(loggerWriteSyncer), logLevel), zap.AddCaller(), zap.AddStacktrace(stackLevel))
	newLoggerService.Logger = logger
	newLoggerService.State = true
	return &newLoggerService
}

func (loggerService *LoggerService) Debugf(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).Debug(fmt.Sprintf(format, args...))
}

func (loggerService *LoggerService) Infof(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprintf(format, args...))
}

func (loggerService *LoggerService) Warnf(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(fmt.Sprintf(format, args...))
}

func (loggerService *LoggerService) Errorf(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).Error(fmt.Sprintf(format, args...))
}

func (loggerService *LoggerService) DPanicf(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).DPanic(fmt.Sprintf(format, args...))
}

func (loggerService *LoggerService) Panicf(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).Panic(fmt.Sprintf(format, args...))
}

func (loggerService *LoggerService) Fatalf(format string, args ...interface{}) {
	loggerService.Logger.WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprintf(format, args...))
}
