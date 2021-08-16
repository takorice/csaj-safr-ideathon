package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	interval = 10
)

type EventsJson struct {
	Events []struct {
		EventId      string  `json:"eventId"`
		StartTime    int64   `json:"startTime"`
		EndTime      int64   `json:"endTime"`
		EventType    string  `json:"type"`
		PersonId     string  `json:"personId"`
		MaxSentiment float32 `json:"maxSentiment"`
		MinSentiment float32 `json:"minSentiment"`
		AvgSentiment float32 `json:"avgSentiment"`
	} `json:"events"`
}

type ReactionSummary struct {
	Id             int
	UserId         int
	LessonId       int
	EmotionalValue float32
	ReactedAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func gormConnect() *gorm.DB {
	//connection := os.Getenv("DATABASE_URL")
	// ローカルでテストする際の接続先を設定
	connection := "host=localhost port=5432 user=postgres dbname=default password=postgres sslmode=disable"
	val, ret := os.LookupEnv("SSL_MODE")
	if ret == false {
		val = "false"
	}

	sslmode, _ := strconv.ParseBool(val)
	if sslmode {
		connection += "?sslmode=require"
	}
	//fmt.Println(connection)
	db, err := gorm.Open("postgres", connection)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

// SAFRのAPIを実行し、eventを取得する
func getEvents(personId string, sinceTime time.Time, untilTime time.Time) []byte {
	// リクエスト先のURLを設定
	reqUrl := "https://cv-event.real.com/events"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	// クエリパラメータを設定
	values := url.Values{}
	values.Add("combineActiveEvents", "false")
	values.Add("personId", personId)
	values.Add("sinceTime", strconv.FormatInt(sinceTime.Unix()*1000, 10))
	values.Add("untilTime", strconv.FormatInt(untilTime.Unix()*1000, 10))
	req.URL.RawQuery = values.Encode()
	// ヘッダ情報を設定
	req.Header.Set("accept", "application/json;charset=UTF-8")
	req.Header.Set("X-RPC-DIRECTORY", "main")
	req.Header.Set("X-RPC-AUTHORIZATION", os.Getenv("AUTH_INFO"))

	//// ヘッダ情報を出力
	//dump, _ := httputil.DumpRequestOut(req, true)
	//fmt.Printf("%s", dump)

	client := new(http.Client)
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		fmt.Println("Status code is not 200.")
	}
	defer resp.Body.Close()
	// req 送信
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

func main() {
	// SAFR の person_id を設定
	var personIds [5]string
	personIds[0] = os.Getenv("SAFR_PERSON_ID_A") // a san
	personIds[1] = os.Getenv("SAFR_PERSON_ID_B") // b san
	personIds[2] = os.Getenv("SAFR_PERSON_ID_C") // c san
	personIds[3] = os.Getenv("SAFR_PERSON_ID_D") // d san
	personIds[4] = os.Getenv("SAFR_PERSON_ID_E") // e san
	// person_id の一覧
	//for _, personId := range personIds {
	//	fmt.Printf("person_id : %s\n", personId)
	//}

	// 開始時刻、終了時刻を設定する
	sinceTime := time.Now().Add(time.Duration(-interval) * time.Minute)
	untilTime := time.Now()

	// 処理実行
	db := gormConnect()
	for i, personId := range personIds { // ユーザ数の分、繰り返す
		// event データを取得する
		events := getEvents(personId, sinceTime, untilTime)
		// レスポンスデータをサンプル出力
		//fmt.Println(string(events))

		// JSONのパースを実施
		var eventsJson EventsJson
		errs := json.Unmarshal(events, &eventsJson)
		if errs != nil {
			fmt.Println("JSON parse is failed.")
		}
		if len(eventsJson.Events) == 0 { // 何も取得できなかった場合は、次のユーザへ
			continue
		}

		// event の内容を解析する
		var emotionalValue float32 = 0
		for _, event := range eventsJson.Events {
			if event.EventType == "person" { // Type が person のデータのみ取り扱う
				if event.StartTime <= untilTime.Unix()*1000 && (event.EndTime == 0 || untilTime.Unix()*1000 <= event.EndTime) {
					// untilTime が、開始時間と終了時間の間にある、または、終了していない場合のデータのみを利用する
					emotionalValue += event.AvgSentiment
					// JSONデータをサンプル出力
					fmt.Printf("personId : %s, eventId : %s, startTime : %d, endTime : %d, avgSentiment:%f\n",
						event.PersonId, event.EventId, event.StartTime, event.EndTime, event.AvgSentiment)
				}
			}
		}
		// insert
		if emotionalValue != 0 { // 感情値が取れる場合のみDBに格納
			reactionSummary := ReactionSummary{
				Id:             0,
				UserId:         i + 1, // ユーザの順番に設定
				LessonId:       1,     // TODO:LessonIDをいい感じに設定する
				EmotionalValue: emotionalValue,
				ReactedAt:      untilTime,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			db.Create(&reactionSummary)
		}
	}
	defer db.Close()
}
