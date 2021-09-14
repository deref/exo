package term

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func GetSize() (w, h int) {
	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return 0, 0
	}
	var size unix.Winsize
	defer out.Close()
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)))
	return int(size.Col), int(size.Row)
}

// VisualLength determines the length of a string (taking into account ansi control sequences)
// From <https://github.com/wagoodman/bashful>.
func VisualLength(str string) int {
	inEscapeSeq := false
	length := 0

	for _, r := range str {
		switch {
		case inEscapeSeq:
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscapeSeq = false
			}
		case r == '\x1b':
			inEscapeSeq = true
		default:
			length++
		}
	}

	return length
}

// TrimToVisualLength truncates the given message to the given length (taking into account ansi escape sequences)
// From <https://github.com/wagoodman/bashful>.
func TrimToVisualLength(message string, length int) string {
	for VisualLength(message) > length && len(message) > 1 {
		message = message[:len(message)-1]
	}
	return message
}
