package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("new logger warning level", func(t *testing.T) {
		testLogger := New("warn", "")
		require.Equal(t, logrus.WarnLevel, testLogger.Level)
	})

	t.Run("new logger info level", func(t *testing.T) {
		testLogger := New("info", "")
		require.Equal(t, logrus.InfoLevel, testLogger.Level)
	})
}
