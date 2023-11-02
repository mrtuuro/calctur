package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"calctur/backend/helpers"
	"calctur/backend/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("/Users/tozay/go/src/calctur/frontend/static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", receiver.handleWS)
	http.HandleFunc("/coords", receiver.handleCoords)
	http.ListenAndServe(":8080", nil)

}

type DataReceiver struct {
	dataCh chan types.Coordinate
	conn   *websocket.Conn
	prod   DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = os.Getenv("KAFKA_TOPIC")
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		dataCh: make(chan types.Coordinate, 128),
		prod:   p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.Coordinate) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  128,
		WriteBufferSize: 128,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	localDR, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	localDR.conn = conn

	go localDR.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New Coordinate client connected.")
	for {
		var coords types.Coordinate
		if err := dr.conn.ReadJSON(&coords); err != nil {
			log.Println("Websocket read error:", err)
			break
		}
		fmt.Println("Received coordinates:", coords)

		if err := dr.produceData(coords); err != nil {
			fmt.Println("Kafka produce error: ", err)
		}
	}
}

func (dr *DataReceiver) handleCoords(w http.ResponseWriter, r *http.Request) {
	var coords types.Coordinate
	if err := json.NewDecoder(r.Body).Decode(&coords); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	fmt.Println("Received coordinates:", coords)
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Coordinates received successfully"})
	return
}
