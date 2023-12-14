package MaterLogService

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	LatestID int64
	File     *os.File
	mu       sync.Mutex
}

func InitLogger(Filename string) (*Logger, error) {
	File, err := os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{File: File, LatestID: 0}, nil
}

func (l *Logger) Log(FileID string, LogType string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	errChan := make(chan error)
	if LogType == "Upload" {
		go func() {
			defer close(errChan)
			CurrentTime := time.Now()
			CurrentTimeFormat := CurrentTime.Format("2006-01-02 15:04:05")
			CurrentIndex := l.LatestID + 1
			MessageFormat := fmt.Sprintf("%10d %s %s %s\n", CurrentIndex, CurrentTimeFormat, FileID, "UP")
			_, err := l.File.WriteString(MessageFormat)
			if err != nil {
				errChan <- err
				return
			}
			l.LatestID++
		}()
	} else if LogType == "Delete" {
		go func() {
			defer close(errChan)
			CurrentTime := time.Now()
			CurrentTimeFormat := CurrentTime.Format("2006-01-02 15:04:05")
			CurrentIndex := l.LatestID + 1
			MessageFormat := fmt.Sprintf("%10d %s %s %s\n", CurrentIndex, CurrentTimeFormat, FileID, "DE")
			_, err := l.File.WriteString(MessageFormat)
			if err != nil {
				errChan <- err
				return
			}
			l.LatestID++
		}()
	} else {
		go func() {
			defer close(errChan)
			CurrentTime := time.Now()
			CurrentTimeFormat := CurrentTime.Format("2006-01-02 15:04:05")
			CurrentIndex := l.LatestID + 1
			MessageFormat := fmt.Sprintf("%10d %s %s %s\n", CurrentIndex, CurrentTimeFormat, FileID, "SE")
			_, err := l.File.WriteString(MessageFormat)
			if err != nil {
				errChan <- err
				return
			}
			l.LatestID++

		}()
	}
	err := <-errChan
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) ReadLog(ID int64) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	const LineLength int64 = 35
	var Index int64 = LineLength * (ID - 1)

	_, err := l.File.Seek(Index, 0)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, LineLength)
	_, err = l.File.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}
