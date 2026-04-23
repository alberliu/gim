package safe

import (
	"log/slog"
	"runtime/debug"
)

func Go(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", "panic", r, "stack", string(debug.Stack()))
			}
		}()

		f()
	}()
}

// RecoverPanic 恢复panic
func RecoverPanic() {
	if r := recover(); r != nil {
		slog.Error("panic recovered", "panic", r, "stack", string(debug.Stack()))
	}
}
