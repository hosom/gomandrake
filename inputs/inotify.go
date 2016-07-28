/*

*/

package inputs

import (
	"log"
	"golang.org/x/exp/inotify"
)

type INotify struct {
	MonitoredPath		string
	AnalysisPipeline	chan string
}

func (i INotify) Monitor() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.AddWatch(i.MonitoredPath, inotify.IN_CLOSE_WRITE)

	for {
		select {
		case ev := <- watcher.Event:
			i.AnalysisPipeline <- ev.Name
		case err := <- watcher.Error:
			log.Printf("INotify Error: %s", err)
		}
	}
}