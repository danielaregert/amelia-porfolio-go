package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	windowFails    = 5
	windowDuration = 15 * time.Minute
	blockDuration  = 30 * time.Minute
	cleanupAge     = 1 * time.Hour
)

type entry struct {
	failCount    int
	firstFail    time.Time
	blockedUntil time.Time
}

type LoginLimiter struct {
	mu     sync.Mutex
	byIP   map[string]*entry
	lastGC time.Time
}

func NewLoginLimiter() *LoginLimiter {
	return &LoginLimiter{
		byIP:   make(map[string]*entry),
		lastGC: time.Now(),
	}
}

func (l *LoginLimiter) Allowed(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.gcLocked()
	e, ok := l.byIP[ip]
	if !ok {
		return true, 0
	}
	now := time.Now()
	if now.Before(e.blockedUntil) {
		return false, e.blockedUntil.Sub(now)
	}
	return true, 0
}

func (l *LoginLimiter) RecordFailure(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.gcLocked()
	now := time.Now()
	e, ok := l.byIP[ip]
	if !ok || now.Sub(e.firstFail) > windowDuration {
		l.byIP[ip] = &entry{failCount: 1, firstFail: now}
		return
	}
	e.failCount++
	if e.failCount >= windowFails {
		e.blockedUntil = now.Add(blockDuration)
	}
}

func (l *LoginLimiter) RecordSuccess(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.byIP, ip)
}

func (l *LoginLimiter) gcLocked() {
	now := time.Now()
	if now.Sub(l.lastGC) < time.Minute {
		return
	}
	l.lastGC = now
	for ip, e := range l.byIP {
		if now.After(e.blockedUntil) && now.Sub(e.firstFail) > cleanupAge {
			delete(l.byIP, ip)
		}
	}
}

func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
