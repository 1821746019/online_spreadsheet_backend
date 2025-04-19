package logger

import (
	"os"

	"github.com/sztu/mutli-table/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupGlobalLogger 用于初始化全局日志
func SetupGlobalLogger(cfg *settings.LoggerConfig) error {
	logger, err := newLogger(cfg)
	if err != nil {
		return err
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}(logger)

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // 启用彩色等级
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	logger = zap.New(core)
	zap.ReplaceGlobals(logger)
	return nil
}
