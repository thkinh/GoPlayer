package audio

import (
	"os"
	"time"

	"fmt"

	"github.com/hajimehoshi/oto/v2"

	"github.com/hajimehoshi/go-mp3"
)

const (
	StorageDir = "assets/songs/"
)

type Song struct {
	Name      string
	Length    int
	IsPlaying bool
}

func (s *Song) Play(stopSignal chan int) {
	defer close(stopSignal)

	errChan := make(chan error, 1)
	go func() {
		errChan <- RunMP3(StorageDir + s.Name)
	}()

	for {
		select {
		case <-stopSignal: //When the song received a stop signal
			s.IsPlaying = false
			return
		case err := <-errChan:
			if err != nil {
				fmt.Println(err)
			}
			s.IsPlaying = false
			return
		}
	}
}

func (s *Song) Pause(pauseSignal chan int) {
	s.IsPlaying = false
}

func RunMP3(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, ready, err := oto.NewContext(d.SampleRate(), 2, 2)
	if err != nil {
		return err
	}
	<-ready

	p := c.NewPlayer(d)
	defer p.Close()
	p.Play()

	//fmt.Sprintf("Length: %d[bytes]\n", d.Length())
	for {
		time.Sleep(time.Second)
		if !p.IsPlaying() {
			break
		}
	}
	return nil
}
