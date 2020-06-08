package main

import (
	"flag"
	"fmt"
	"github.com/clearcodecn/cargo/cluster"
	"github.com/clearcodecn/cargo/proto"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	connections int
)

func init() {
	flag.IntVar(&connections, "n", 100, "how much connection to create")
}

func main() {

	cluster.RegisterHandle(proto.Text, func(ctx *cluster.Context, v *proto.Message) error {
		return nil
	})

	wg := sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())
	n := connections - 1
	for i := 0; i < n; i++ {
		wg.Add(1)
		j := i
		go func() {
			defer wg.Done()
			tk := time.NewTicker(2 * time.Second)
			defer tk.Stop()
			client := newClient(j)
			defer client.Close()

			for {
				select {
				case <-tk.C:
					client.WriteMessage(&proto.Message{
						Type: proto.Text,
						Body: []byte(randomStr(rand.Intn(30))),
					})
				}
			}
		}()
	}
	wg.Wait()

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//
	//})
	//http.ListenAndServe(":1111", nil)

}

func newClient(n int) *cluster.Context {
	d := websocket.Dialer{}
	conn, _, err := d.Dial("ws://0.0.0.0:6300/ws", nil)
	ctx := cluster.NewClientContext(conn)
	ctx.WriteMessage(&proto.Message{
		Type: proto.Auth,
		Body: []byte(fmt.Sprintf("client - %d", n)),
	})

	if err != nil {
		log.Fatal(err)
	}
	return ctx
}

var text = []rune(`（观察者网讯）继上月28日二读通过后，香港立法会于4日继续审议《国歌条例草案》。据“东网”等港媒报道称，现场反对派议员提出了21项修正案，但全部都被否决。按照立法会主席梁君彦所订立的时间表，预计4日傍晚将进行草案的三读表决。此外，4日的立法会现场再度发生反对派议员撒泼闹事的情况，中午12时50分许，议员陈志全和朱凯迪突然冲向主席台，并向会议厅内投掷有臭味的不明液体，里面甚至还有白色的虫，这导致会议被迫暂停。朱凯迪事后称不明液体是“有机肥料”，大批警员和消防员先后抵达会议厅调查。
反对派议员所投掷的含白色虫不明液体 图自“东网”
会上，建制派议员葛珮帆指出，设立《国歌法》是国际的一贯规则，反对派议员的反对理由没有道理。她表示，4年前就有议员在立法会宣誓时“耍花招”，令香港市民十分反感。葛珮帆表示，要香港市民尊重国歌，不能单单靠《国歌法》立法，立法只是国民教育的一部分。
多名反对派议员提出21项修正案，香港特区政制及内地事务局局长曾国卫表示反对这21项修正案，并逐项解释反对理由。对于判断侮辱国歌行为的问题，他表示有关部门会在搜集证据后再作检控，最终由法庭作出公平判决。《国歌条例草案》已列明构成罪行的要素，若市民无意图侮辱国歌，根本不会触犯法例。
香港特区政制及内地事务局局长曾国卫发言 图自香港中通社
建制派议员陈建波则认为，公众对香港的法庭和司法制度有信心，反对派议员所提出的修正案毫无道理，只是为了反对而反对。
建制派议员邵家辉称，近年来在香港进行的足球比赛中经常出现有人对国歌发出嘘声以及做出侮辱手势的情况，这令香港声誉受损。他认为嘘国歌属于“港独”思想之一，有着想和国家分离的意味。
此外，在今天的会议中再度发生了不和谐的场面。
据港媒“橙新闻”报道，4日中午12时50分左右，在建制派议员陆颂雄发言期间，反对派议员陈志全和朱凯迪走向会场主席台抗议，随后又投掷有臭味并附有白色虫的不明液体。会议随即被迫暂停，工作人员开始清理现场。
而在上月28日《国歌条例草案》的二读辩论现场，当时反对派议员就曾向主席台丢掷气味刺鼻难闻的“臭弹”，大闹立法会会议厅。`)

func randomStr(n int) string {
	var d []rune
	var l = len(text)
	for i := 0; i < n; i++ {
		d = append(d, text[rand.Intn(l)])
	}
	return string(d)
}
