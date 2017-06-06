/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package util

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type transferQueue interface {
	UseClient(*http.Client) error
	Transfer(TransferDefinition) error
	Done()
}

/*
TransferQueue uses a buffered channel as a queue of transfer definitions and
downloads each URL to the specified file
*/
type TransferQueue struct {
	Client *http.Client
	Queue  chan TransferDefinition
	Quit   chan int
}

/*
UseClient allows the user of this function to provide an HTTP client that
will be used for all downloads
*/
func (tq *TransferQueue) UseClient(newClient *http.Client) error {
	if newClient == nil {
		return errors.New("Please specify an HTTP client")
	}

	tq.Client = newClient

	tq.Queue = make(chan TransferDefinition, 15)
	tq.Quit = make(chan int)
	go tq.processQueueItems()

	return nil
}

/*
Done tells the queue that no more transfers will be requested
*/
func (tq *TransferQueue) Done() {
	close(tq.Queue)

	// Wait until we receive the signal that the queue is empty and all processing is complete
	<-tq.Quit
}

/*
Transfer adds a TransferDefinition structure to the queue for future downloading
*/
func (tq *TransferQueue) Transfer(tf TransferDefinition) error {
	if tf.Source == nil {
		return errors.New("Please specify a URL source for the transfer")
	}
	if len(tf.Target) == 0 {
		return errors.New("Please specify a local filesystem path for the transfer")
	}
	if tq.Client == nil {
		return errors.New("Please call the UseClient method before using this method")
	}

	tq.Queue <- tf
	return nil
}

func (tq *TransferQueue) processQueueItems() {
	i := 0
	for tf := range tq.Queue {
		if i > 0 {
			time.Sleep(3 * time.Second)
		}
		err := SaveURLToPath(tf.Source, tf.Target)
		if err != nil {
			fmt.Println(err)
		}
		i++
	}

	tq.Quit <- 0
}
