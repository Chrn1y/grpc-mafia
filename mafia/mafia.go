package mafia

import (
	"errors"
	"fmt"
	mafiaproto "github.com/Chrn1y/grpc-mafia/proto"
	"math/rand"
)

type InputChan chan *mafiaproto.Request
type OutputChan chan *mafiaproto.Response


type Mafia struct {
	players map[string]InputChan
	output map[string]OutputChan
	mafia string
	num   int
	started chan struct{}
}

func (m *Mafia) Join(player string) (InputChan, OutputChan, error) {
	//println("connecting:", player)
	select{
	case <-m.started:
		return nil, nil, errors.New("game already started")
	default:
	}
	if _, ok := m.players[player]; ok {
		return nil, nil, errors.New("player with this name already exists")
	}
	m.sendInfoAll(player + " присоединяется к игре", false)
	m.players[player] = make(chan *mafiaproto.Request, 1)
	m.output[player] = make(chan *mafiaproto.Response, 1)
	if len(m.players) == m.num {
		//println("enough")
		go m.start()
	} else {
		m.sendInfoAll(fmt.Sprintf("Ждем необходимое количество игроков: %d\nСейчас: %d", m.num, len(m.players)), false)
	}
	return m.players[player], m.output[player], nil
}

func (m *Mafia) sendInfo(to, info string, end bool) {
	m.output[to] <- &mafiaproto.Response{Data: &mafiaproto.Response_Info_{Info: &mafiaproto.Response_Info{
		Text: info,
		End: end}},
		Type: mafiaproto.ResponseType_info,
	}
}

func (m *Mafia) sendInfoAll(info string, end bool) {
	for p := range m.players {
		m.sendInfo(p, info, end)
	}
}

func (m *Mafia) sendVote(to, msg string, choose []string) string {
	m.output[to] <- &mafiaproto.Response{Data: &mafiaproto.Response_Vote_{Vote: &mafiaproto.Response_Vote{
		Text: msg,
		Choose: choose,
	}},
	Type: mafiaproto.ResponseType_vote_response,
	}
	temp := <- m.players[to]
	return temp.Data.(*mafiaproto.Request_Vote_).Vote.Name
}

func (m *Mafia) getCitizens() []string {
	out := make([]string, 0)
	for p := range m.players {
		if p != m.mafia {
			out = append(out, p)
		}
	}
	return out
}

func (m *Mafia) getPlayers() []string {
	out := make([]string, 0)
	for p := range m.players {
		out = append(out, p)
	}
	return out
}

func (m *Mafia) start() {
	close(m.started)
	defer m.close()
	//println("starting")
	temp := m.getPlayers()
	m.mafia = temp[rand.Intn(len(temp))]
	//println("mafia")
	m.sendInfo(m.mafia, "Ты мафия", false)
	//println("sent to", m.mafia)
	for p := range m.players{
		m.sendInfo(p, "Игра начинается", false)
	}
	for {
		//println("начинаем раунд")
		for p := range m.players{
			m.sendInfo(p, "Начинается ночь", false)
		}
		ok := false
		killed := ""
		for !ok {
			killed = m.sendVote(m.mafia, "Мафия, выбери мирного жителя:", m.getCitizens())
			_, ok = m.players[killed]
			if !ok {
				m.sendInfo(m.mafia, "Выбери существующего мирного жителя", false)
				continue
			}
		}
		m.sendInfoAll("Город просыпается без " + killed, false)
		m.sendInfo(killed, "Для тебя игра закончена", true)
		delete(m.players, killed)
		if len(m.players) <= 2 {
			m.sendInfoAll("Мафия победила!\nКонец игры", true)
			break
		}

		m.sendInfoAll("Начинается голосование", false)

		voting := make(map[string]int)
		for p := range m.players {
			ok = false
			for !ok {
				voted := m.sendVote(p, "Проголосуй за мафию", m.getPlayers())
				_, ok = m.players[voted]
				if !ok {
					m.sendInfo(p, "Выбери существующего игрока", false)
					continue
				}
				voting[voted]++
			}
		}
		voted := ""
		votes := 0
		for v := range voting {
			if voting[v] >= votes {
				voted = v
			}
		}
		m.sendInfoAll("Игроки выбрали " + voted, false)
		m.sendInfo(voted, "Для тебя игра закончена", true)
		delete(m.players, voted)
		if voted == m.mafia {
			m.sendInfoAll("Мирные жители победили!\nКонец игры", true)
			break
		}
	}
}

func (m *Mafia) close() {
	for p := range m.players {
		close(m.players[p])
	}
	for out := range m.output {
		close(m.output[out])
	}
	m.mafia = ""
	m.players = make(map[string]InputChan)
	m.output = make(map[string]OutputChan)
	m.started = make(chan struct{})
}

func New(num int) *Mafia{
	return &Mafia{
		players:  make(map[string]InputChan),
		output:   make(map[string]OutputChan),
		num:      num,
		started:  make(chan struct{}),
	}
}