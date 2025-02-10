package tui

import (
    "fmt"
    "os"
    audio "AP/audio"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbletea"
    )

//type Timer struct {
//	TimeSet     int
//	TimeCurrent int
//	resetSignal chan int
//}

type Styles struct {
    BorderColor lipgloss.Color
    header lipgloss.Style
    playlistBox lipgloss.Style
    playingBox  lipgloss.Style
    cursorStyle lipgloss.Style
}

type SongEndedMsg struct {}

type model struct {
  width    int
    height   int
    header   string
    styles   *Styles
    current  *audio.Song
    playlist []audio.Song
    cursor   int
    //timer    Timer
    stopSignal chan int
}

func InitialModel() model {
playlist := make([]audio.Song, 0)

            entries, err := os.ReadDir(audio.StorageDir)
            if err != nil {
              fmt.Println(err)
            }

          for _, entry := range entries {
            if !entry.IsDir() {
song := audio.Song{
Name:   entry.Name(),
          Length: 0,
      }
      playlist = append(playlist, song)
            }
          }

styles := DefaultStyles()

          return model{
styles: styles,
          header:   "\nMy Audio Player\n",
          playlist: playlist,
          cursor:   0,
          }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {

    case tea.WindowSizeMsg:
      m.width = msg.Width
        m.height = msg.Height

    case SongEndedMsg: 

    case tea.KeyMsg:
        switch msg.String() {
          case "ctrl+c", "q":
            return m, tea.Quit
          case "j":
              if m.cursor < len(m.playlist)-1 {
                m.cursor++
              }
          case "k":
              if m.cursor > 0 {
                m.cursor--
              }
          case "enter":

              //When choosing another song while the current is playing
              if m.current != nil && m.current.IsPlaying {
                m.stopSignal <- 0
              }

              //When the choosen song is not the current song
              if m.current != &m.playlist[m.cursor] {
                m.current = &m.playlist[m.cursor]
              }


              m.stopSignal = make(chan int)
                m.current.IsPlaying = true

                return m, func() tea.Msg {
                  m.current.Play(m.stopSignal)
                    return SongEndedMsg {}
                }

        }
  }

  return m, nil
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) View() string {
	// Header
	Head := m.styles.header.Render(m.header)

	// Playlist
	var playlist string
	for i, song := range m.playlist {
		songStyle := lipgloss.NewStyle()
		if m.cursor == i {
        songStyle = m.styles.cursorStyle
		}

		playlist += songStyle.Render(fmt.Sprintf("[%v]. %v", i, song.Name))
    playlist += "\n"
	}

	ListBox := m.styles.playlistBox.Render(playlist)

	// Playing
	var playing string
	playing += fmt.Sprintf("PLAYING SONG:\n")
	for i, song := range m.playlist {
		if song.IsPlaying {
			playing += fmt.Sprintf("%v. %v", i, song.Name)
		}
	}

	PlayingBox := m.styles.playingBox.Render(playing)

	// Render Layout
	return lipgloss.JoinVertical(lipgloss.Center,
		Head,
		lipgloss.JoinHorizontal(lipgloss.Left,
			ListBox, PlayingBox))
}


func DefaultStyles() *Styles {
s := new(Styles)
     s.BorderColor = lipgloss.Color("#16C47F")

     s.header = lipgloss.NewStyle().
     Foreground(lipgloss.Color("#FFFFFF")).  // White text
     Background(lipgloss.Color("#5A3FCF")).  // Background fills text area
     BorderForeground(s.BorderColor).       // Purple border color
     BorderStyle(lipgloss.DoubleBorder()).  // Elegant double border
     Bold(true).
     Align(lipgloss.Center).                // Center the text
     Padding(1, 2).                          // Ensure proper spacing
     Margin(0).                              // Remove extra margin
     Width(80)   

     s.playlistBox = lipgloss.NewStyle().
     BorderForeground(s.BorderColor).
     BorderStyle(lipgloss.ThickBorder()).
     Padding(0).Width(39)


     s.playingBox = lipgloss.NewStyle().
     BorderForeground(s.BorderColor).
     BorderStyle(lipgloss.ThickBorder()).
     Padding(0).Width(39).
     Align(lipgloss.Left)
      


     s.cursorStyle = lipgloss.NewStyle().
     Foreground(lipgloss.Color("#FFFFFF")).  
     Background(lipgloss.Color("#7D56F4")). 
     Padding(0)

     return s
}


//func (t *Timer) SetTime(set int) {
//		t.TimeSet =	set
//		t.TimeCurrent = 0
//}
