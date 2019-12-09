package wait

import (
	"fmt"
	"math/rand"
	"time"
)

// Until loops until stop channel is closed, running f every period.
//
// Until is syntactic sugar on top of JitterUntil with zero jitter factor and
// with sliding = true (which means the timer for period starts after the f
// completes).
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}

// JitterUntil loops until stop channel is closed, running f every period.
//
// If jitterFactor is positive, the period is jittered before every run of f.
// If jitterFactor is not positive, the period is unchanged and not jittered.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
//
// Close stopCh to stop. f may not be invoked if stop channel is already
// closed. Pass NeverStop to if you don't want it stop.
func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	var t *time.Timer
	var sawTimeout bool
	
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		
		jitteredPeriod := period
		if jitterFactor > 0.0 {
			jitteredPeriod = Jitter(period, jitterFactor)
		}
		
		if !sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}
		
		func() {
			defer func(){
				r := recover()
				if r != nil {
					fmt.Println(r.(string))
				}
			}()
			f()
		}()
		
		if sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}
		
		// NOTE: b/c there is no priority selection in golang
		// it is possible for this to race, meaning we could
		// trigger t.C and stopCh, and t.C select falls through.
		// In order to mitigate we re-check stopCh at the beginning
		// of every loop to prevent extra executions of f().
		select {
		case <-stopCh:
			return
		case <-t.C:
			sawTimeout = true
		}
	}
}

// Jitter returns a time.Duration between duration and duration + maxFactor *
// duration.
//
// This allows clients to avoid converging on periodic behavior. If maxFactor
// is 0.0, a suggested default value will be chosen.
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// resetOrReuseTimer avoids allocating a new timer if one is already in use.
// Not safe for multiple threads.
func resetOrReuseTimer(t *time.Timer, d time.Duration, sawTimeout bool) *time.Timer {
	if t == nil {
		return time.NewTimer(d)
	}
	if !t.Stop() && !sawTimeout {
		<-t.C
	}
	t.Reset(d)
	return t
}
