package main

import (
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/yafeng-Soong/base-shadowsocks2/socks"
	"github.com/yafeng-Soong/base-shadowsocks2/statistic"
)

func tcpLocal(l net.Listener, d Dialer) {
	ch := make(chan *statistic.TrackerInfo, 100)
	statistic.TrackerInfoChan = ch
	go statistic.HandleMetric(ch)
	for {
		c, err := l.Accept()
		if err != nil {
			logf("failed to accept: %v", err)
			continue
		}
		go func() {
			defer c.Close()
			laddr := c.LocalAddr()
			if laddr == nil {
				logf("failed to determine target address")
				return
			}
			rc, err := d.Dial(laddr.Network(), laddr.String())
			if err != nil {
				logf("failed to connect: %v", err)
				return
			}
			defer rc.Close()
			logf("proxy %s <--[%s]--> %s", c.RemoteAddr(), rc.RemoteAddr(), laddr)
			right := statistic.NewTcpTracker(c, ch, laddr.String())
			if err = relay(rc, right); err != nil {
				logf("relay error: %v", err)
			}
			right.WindUp()
		}()
	}
}

// Listen on addr for incoming connections.
func tcpRemote(addr string, shadow func(net.Conn) net.Conn) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logf("failed to listen on %s: %v", addr, err)
		return
	}

	logf("listening TCP on %s", addr)
	for {
		c, err := l.Accept()
		if err != nil {
			logf("failed to accept: %v", err)
			continue
		}

		go func() {
			defer c.Close()
			c = shadow(c)

			tgt, err := socks.ReadAddr(c)
			if err != nil {
				logf("failed to get target address from %v: %v", c.RemoteAddr(), err)
				return
			}

			rc, err := net.Dial("tcp", tgt.String())
			if err != nil {
				logf("failed to connect to target: %v", err)
				return
			}
			defer rc.Close()

			logf("proxy %s <-> %s", c.RemoteAddr(), tgt)
			if err = relay(c, rc); err != nil {
				logf("relay error: %v", err)
			}
		}()
	}
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now()) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now()) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
