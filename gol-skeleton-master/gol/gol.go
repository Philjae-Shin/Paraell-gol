package gol

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {
	// Initialize the DistributorChannels
	c := NewDistributorChannels()

	// Start the Distributor
	go Distributor(p, c, keyPresses)

	// Process events from DistributorChannels and forward them to the events channel
	go func() {
		for {
			c.eventMu.Lock()
			for len(c.events) == 0 {
				c.eventCond.Wait()
			}
			eventsToProcess := c.events
			c.events = nil
			c.eventMu.Unlock()

			for _, event := range eventsToProcess {
				events <- event
				// If the event is FinalTurnComplete, we can close the events channel
				if _, ok := event.(FinalTurnComplete); ok {
					close(events)
					return
				}
			}
		}
	}()
}
