// ChatParser project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	//	"os/exec"
	//	"path/filepath"
	//	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gempir/go-twitch-irc"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"
)

var mykey = ""

type masschat struct {
	*walk.MainWindow
	edit *walk.TextEdit
	chat *walk.TextEdit
	clog *walk.TextEdit
	path string
}
type spy struct {
	*walk.MainWindow
	info *walk.TextEdit
}

func main() {
	GetKey()
	myfont := new(Font)
	myfont.Bold = true
	myfont.PointSize = 10
	var rooms []string
	var wg sync.WaitGroup
	mw := &masschat{}
	flag2 := 1

	MW := MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Parse chat",
		MinSize:  Size{500, 50},
		Size:     Size{500, 50},
		Layout:   VBox{},

		Children: []Widget{
			TextEdit{
				Text:     "Тут будут логи по сообщениям",
				MinSize:  Size{500, 100},
				Font:     *myfont,
				AssignTo: &mw.chat, ReadOnly: true,
			},
			TextEdit{
				Text:     "Логи для тебя.",
				MinSize:  Size{500, 20},
				Font:     *myfont,
				AssignTo: &mw.clog, ReadOnly: true,
			},

			TextEdit{
				Text:     "Название канала",
				AssignTo: &mw.edit, ReadOnly: false,
			},

			PushButton{
				MinSize: Size{50, 20},
				MaxSize: Size{50, 20},
				Text:    "Начать парсинг",
				OnClicked: func() {
					if parseUrl() == 1 {

						flag := 1
						for _, value := range rooms {
							if mw.edit.Text() == value {
								flag = 0
								mw.clog.SetText(mw.edit.Text() + " Уже был запущен ранее...")
							}
						}
						if flag == 1 {
							rooms = append(rooms, mw.edit.Text())
							mw.clog.SetText(mw.edit.Text() + "Запущен")
							wg.Add(1)
							go func(room string, flag2 *int, chat *walk.TextEdit) {
								*flag2 = 1
								defer wg.Done()
								spyRoom(room, flag2, chat)

							}(mw.edit.Text(), &flag2, mw.chat)
						}
					} else {
						mw.chat.SetText(mykey + " Отдай ключ разработчику.")
					}
				},
			},
			PushButton{
				MinSize:   Size{100, 30},
				MaxSize:   Size{100, 30},
				Text:      "Очистить дубли",
				OnClicked: mw.pbClicked,
			},
			PushButton{
				MinSize: Size{100, 30},
				MaxSize: Size{100, 30},
				Text:    "выход(отключит всех ботов)",
				OnClicked: func() {
					os.Exit(1)
				},
			},
		},
	}

	if _, err := MW.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func spyRoom(room string, flag *int, chat *walk.TextEdit) {

	f2, err2 := os.OpenFile(room+".txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err2 != nil {
		f, err3 := os.Create(room + ".txt")
		if err3 != nil {
			fmt.Println("Unable to create file:", err3)
			os.Exit(1)
		} else {
			defer f.Close()
			chat.SetText("Процесс идёт.")
			client := twitch.NewClient("Meemal6412", "oauth:3e139brdfy7jxpki47ptdpcguzpn5o")
			client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
				if user.Username != "nightbot" {
					if strings.Contains(message.Text, "http") || strings.Contains(message.Text, ".ru") || strings.Contains(message.Text, ".com") || strings.Contains(message.Text, ".gd") || strings.Contains(message.Text, ".click") || strings.Contains(message.Text, ".do") || strings.Contains(message.Text, ".ly") || strings.Contains(message.Text, ".to") || strings.Contains(message.Text, ".us") || strings.Contains(message.Text, ".me") || strings.Contains(message.Text, ".org") {
					} else {
						if _, err3 = f.WriteString(message.Text + "\r\n"); err3 != nil {
							panic(err3)
						}
						chat.SetText("Процесс идёт.\r\nЗаписали в файл :" + message.Text)
					}
				}
			})
			client.Join(room)
			err2 := client.Connect()
			if err2 != nil {
				panic(err2)
			}
		}
	} else {
		defer f2.Close()
		chat.SetText("Процесс идёт.")
		client := twitch.NewClient("Meemal6412", "oauth:3e139brdfy7jxpki47ptdpcguzpn5o")
		client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
			if user.Username != "nightbot" {
				if strings.Contains(message.Text, "http") || strings.Contains(message.Text, ".ru") || strings.Contains(message.Text, ".com") || strings.Contains(message.Text, ".gd") || strings.Contains(message.Text, ".click") || strings.Contains(message.Text, ".do") || strings.Contains(message.Text, ".ly") || strings.Contains(message.Text, ".to") || strings.Contains(message.Text, ".us") || strings.Contains(message.Text, ".me") || strings.Contains(message.Text, ".org") {
				} else {
					if _, err2 = f2.WriteString(message.Text + "\r\n"); err2 != nil {
						panic(err2)
					}
					chat.SetText("Процесс идёт.\r\nЗаписали в файл :" + message.Text)
				}
			}
		})

		client.Join(room)
		err4 := client.Connect()
		if err4 != nil {
			panic(err4)
		}
	}

}

func (mw *masschat) pbClicked() {
	GetKey()
	dlg := new(walk.FileDialog)

	dlg.FilePath = mw.path
	dlg.Title = "Select File"
	dlg.Filter = "Exe files (*.txt)|*.txt|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("Error : File Open\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("Cancel\r\n")
		return
	}
	mw.path = dlg.FilePath

	zz, err := ioutil.ReadFile(dlg.FilePath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	m := make(map[string]string)

	for _, value := range strings.Split(string(zz), "\n") {
		m[value] = value
	}
	fmt.Println(len(m))
	f, err2 := os.Create(dlg.FilePath)
	if err2 != nil {
	} else {
		for _, value := range m {
			if _, err2 = f.WriteString(value + "\n"); err2 != nil {
				panic(err2)
			}
		}
	}
}

func SetKey() string {
	file, err := os.Create("key.txt")
	if err != nil {
		return "bad"
	}
	defer file.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	r := randomString(30)
	file.WriteString(r)
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run\`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err := k.SetStringValue("getchat", r); err != nil {
		fmt.Println(err)
	}
	z, _, err := k.GetStringValue("getchat")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(z)
	if err := k.Close(); err != nil {
		fmt.Println(err)
	}
	return (z)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseUrl() int {
	url := "https://crazyhomeless.livejournal.com/835.html"
	doc, err := goquery.NewDocument(url)
	_check(err)
	tkey := ""
	flag := 0
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		if flag == 1 {
			tkey = strings.TrimSpace(s.Text())
		}
		flag++
	})
	ttkey := strings.Split(tkey, "-")
	for _, value := range ttkey {
		if value == mykey {
			return 1
		}
	}
	return 0
}

func GetKey() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run\`, registry.QUERY_VALUE|registry.SET_VALUE)
	//k, err := registry.OpenKey(registry.CURRENT_USER, `Software`, registry.QUERY_VALUE|registry.SET_VALUE)
	z, _, err := k.GetStringValue("getchat")
	if err != nil {
		mykey = SetKey()
	} else {
		mykey = string(z)
	}

}
