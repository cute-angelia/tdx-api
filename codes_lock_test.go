package tdx

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestCodesUpdateLockSerializesSameDatabaseKey(t *testing.T) {
	releaseFirst := make(chan struct{})
	firstHasLock := make(chan struct{})
	secondFinished := make(chan struct{})
	var concurrent int32

	go func() {
		unlock := acquireCodesUpdateLock("codes.db")
		close(firstHasLock)
		<-releaseFirst
		unlock()
	}()

	<-firstHasLock

	go func() {
		unlock := acquireCodesUpdateLock("codes.db")
		atomic.StoreInt32(&concurrent, 1)
		unlock()
		close(secondFinished)
	}()

	select {
	case <-secondFinished:
		t.Fatal("second lock acquisition should wait for the first unlock")
	case <-time.After(50 * time.Millisecond):
	}

	close(releaseFirst)

	select {
	case <-secondFinished:
	case <-time.After(time.Second):
		t.Fatal("second lock acquisition did not complete after unlock")
	}

	if atomic.LoadInt32(&concurrent) != 1 {
		t.Fatal("second lock acquisition never entered critical section")
	}
}

func TestCodesUpdateLockDoesNotBlockDifferentDatabaseKeys(t *testing.T) {
	firstHasLock := make(chan struct{})
	secondFinished := make(chan struct{})
	releaseFirst := make(chan struct{})

	go func() {
		unlock := acquireCodesUpdateLock("codes-a.db")
		close(firstHasLock)
		<-releaseFirst
		unlock()
	}()

	<-firstHasLock

	go func() {
		unlock := acquireCodesUpdateLock("codes-b.db")
		unlock()
		close(secondFinished)
	}()

	select {
	case <-secondFinished:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("different database keys should not block each other")
	}

	close(releaseFirst)
}
