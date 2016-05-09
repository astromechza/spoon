package logging

import (
    "os"
    "sync"
    "time"

    "github.com/op/go-logging"

    "github.com/AstromechZA/spoon/conf"
)

var log = logging.MustGetLogger("spoon.logging")
var logFormat = logging.MustStringFormatter(
    "%{time:2006-01-02T15:04:05.000} %{module} %{level:.4s} - %{message}",
)

// Initial sets up the initial logging configuration to stdout
func Initial() {
    logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), logFormat))
    log.Debug("Set up basic stdout logger")
}

// TODO
// would be cool to also add a stderr logger for ERROR and CRITICAL messages
// so that they get seen in both the file and stream

// Reconfigure the logging to follow the given logging configuration
func Reconfigure(logcfg *conf.SpoonConfigLog) {
    log.Debugf("Loaded logging configuration: %v", *logcfg)

    if logcfg.Path != "-" {
        log.Debugf("Logging configuration specified path: %s", logcfg.Path)
        rotatingWriter, err := NewRotatingWriter(logcfg.Path, logcfg.RotateSize)
        if err == nil {
            log.Infof("Switching to %s", logcfg.Path)
            logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(rotatingWriter, "", 0), logFormat))
        } else {
            log.Errorf("Error during log reconfiguration: %s", err.Error())
        }
    }
}

type RotatingWriter struct {
    lock sync.Mutex
    filename string
    fp *os.File
    sizeLimit int64
    written int64
}

func NewRotatingWriter(filename string, sizeLimit int64) (*RotatingWriter, error) {
    w := &RotatingWriter{filename: filename, sizeLimit: sizeLimit}
    err := w.Rotate()
    if err != nil {
        return nil, err
    }
    return w, nil
}

func (self *RotatingWriter) Write(logbytes []byte) (int, error) {
    // before writing, check if the logfile is too big
    if self.written > self.sizeLimit {
        err := self.Rotate()
        if err != nil {
            log.Critical("Failed to rotate log file")
            return 0, err
        }
    }

    // then get the lock for the rest of the method
    self.lock.Lock()
    defer self.lock.Unlock()

    // write the content
    w, err := self.fp.Write(logbytes)
    self.written += int64(w)
    return w, err
}

func (self *RotatingWriter) Rotate() error {
    // first grab the lock, we dont want anything to log while we're busy
    // rotating things
    self.lock.Lock()
    defer self.lock.Unlock()

    // now securely close the existing file if we have one and its open
    if self.fp != nil {
        err := self.fp.Close()
        self.fp = nil
        if err != nil {
            return err
        }
    }

    // get the file stat and size
    stat, err := os.Stat(self.filename)
    if err == nil {

        // is rotate required?
        if stat.Size() > self.sizeLimit {
            newFilename := self.filename + "." + time.Now().Format("20060102T150405")
            err = os.Rename(self.filename, newFilename)
            if err != nil { return err }

        // if no rotate is required, just open the same file
        } else {
            self.fp, err = os.OpenFile(self.filename, os.O_RDWR|os.O_APPEND, 0660)
            if err != nil { return err }
            self.written = stat.Size()
            return nil
        }
    }

    self.fp, err = os.Create(self.filename)
    if err != nil { return err }
    self.written = 0
    return nil
}
