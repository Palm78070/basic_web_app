package handlers

import (
	"fmt"
	"os"
)

func (app *App) Cleanup(sig_ch chan os.Signal, done chan struct {}) {
	received_sig_ch := <-sig_ch
	fmt.Println("Received signal: ", received_sig_ch)
	close(app.chat.broadcast)
	app.wg.Wait()

	app.mutex.Lock()
	defer app.mutex.Unlock()
	for k, v := range app.chat.clients {
		if v != nil {
			v.Close()
		}
		delete(app.chat.clients, k)
	}
	for k := range app.chat.rooms{
		delete(app.chat.rooms, k)
	}
	app.chat.rooms = nil
	fmt.Println("Finish cleanup")
	close(done)
	close(sig_ch)
	os.Exit(0)
}
