// Package lockfile is a Linux tool to lock cooperating processes based on syscall `flock`.
// While a sync.Mutex helps against concurrency issues within a single process, this package
// is designed to help against concurrency issues between cooperating processes.
// This package can be only used in Linux, and cannot be used as a `sync.Mutex` in a single process.
package lockfile

import (
	"fmt"
	"os"
	"syscall"
)

// TemporaryError is a type of error where a retry after a random amount of sleep should help to mitigate it.
type TemporaryError string

func (t TemporaryError) Error() string   { return string(t) }
func (t TemporaryError) Temporary() bool { return true }

var (
	// ErrBusy means that the lock is being acquired by anthor process.
	// If you get this, retry after a short sleep might help.
	ErrBusy = TemporaryError("locked by other process")
)

// Lockfile is a Linux signal file to implement cross-process locks. The file content
// does not contain anything, it will be always overwritten.
//
// The Lockfile should be used in different process, it cannot be used as a mutex lock
// in a single process. If you use it in a single process, the Lock and TryLock will
// always success.
//
// If the process crashed, all its Lockfile will be released automatically, you donot
// need to handle the deadlock situation.
type Lockfile struct {
	// Path represents the file lock path. It should be shared between different process.
	// Each path can represent an independent lock file.
	Path string

	file *os.File
	fd   int
}

// New creates a new Lockfile.
func New(path string) *Lockfile {
	return &Lockfile{Path: path}
}

// Lock blocks until the lock is acquired successfully.
func (lf *Lockfile) Lock() error {
	return lf.acquireLock(true)
}

// TryLock tries to acquire the lock.
// If the current process successful acquire the lock, this will return nil.
// If the lock is busy, this will return ErrBusy. In this case, you should do some retries.
func (lf *Lockfile) TryLock() error {
	return lf.acquireLock(false)
}

func (lf *Lockfile) acquireLock(block bool) error {
	var err error
	if lf.file == nil {
		// The lockfile is only a semaphore for us, we won't write anything in it.
		// So it is safe to call `os.Create` here to overwrite the lockfile.
		lf.file, err = os.Create(lf.Path)
		if err != nil {
			return err
		}
		// `fd` will be used to do system call.
		lf.fd = int(lf.file.Fd())
	}

	// Linux system call: flock, see: https://man7.org/linux/man-pages/man2/flock.2.html
	// This is an advisory lock that won't block file writing and reading.
	// LOCK_EX: Place an exclusive lock.  Only one process may hold an exclusive lock for
	//          a given file at a given time.
	// LOCK_NB: No block when the lock is busy.
	flag := syscall.LOCK_EX
	if !block {
		flag |= syscall.LOCK_NB
	}
	err = syscall.Flock(lf.fd, flag)
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			if errno.Temporary() {
				return ErrBusy
			}
		}
		return fmt.Errorf("flock error: %v", err)
	}

	return nil
}

// Unlock actively releases the lock.
//
// You donot need to call this if your lock works on the whole process, because the lock
// will be automatic released by Linux after process ends.
// This method only needs to be called if you want to release the lock before process ends.
func (lf *Lockfile) Unlock() error {
	if lf.file == nil {
		panic("please call TryLock or Lock before Unlock")
	}
	// Use LOCK_UN flag to release the flock.
	err := syscall.Flock(lf.fd, syscall.LOCK_UN)
	if err != nil {
		return fmt.Errorf("flock error: %v", err)
	}
	return lf.file.Close()
}
