package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool/v2"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type HttpServer struct {
	conn    *tarantool.Connection
	wallets []string
}

func main() {
	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect("127.0.0.1:3301", opts)
	if err != nil {
		fmt.Printf("Connection refused: %s\n", err.Error())
		return
	}
	defer conn.Close()

	srv := HttpServer{
		conn:    conn,
		wallets: saveWallets(conn),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tarantool", srv.updateTransactionTarantool)

	go func() {
		err := http.ListenAndServe(":3334", mux)
		log.Println(err)
	}()

	select {}
}

func (s *HttpServer) updateTransactionTarantool(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /tarantoool request\n")
	start := time.Now()
	id := saveWallet(s.conn)
	resp, err := s.conn.Call("transfer", []interface{}{s.wallets[0], s.wallets[0], 5})
	elapsed := time.Since(start)
	fmt.Println(id, elapsed, resp, err)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func saveWallets(conn *tarantool.Connection) []string {
	ids := make([]string, 10)
	for i := 0; i < 10; i++ {
		id := saveWallet(conn)
		ids[i] = id
	}

	return ids
}

func saveWallet(conn *tarantool.Connection) string {
	minAmount := 10.0
	maxAmount := 100.0
	id := uuid.New().String()
	memberId := uuid.New().String()
	hash := md5.Sum([]byte(uuid.New().String()))
	strHash := hex.EncodeToString(hash[:])
	randAmount := minAmount + rand.Float64()*(maxAmount-minAmount)
	randAmount = math.Round(randAmount)
	q := fmt.Sprintf("INSERT INTO wallets VALUES('%s', '%s', '%s', %d , %d);", id, strHash, memberId, int(randAmount), time.Now().UnixNano())
	resp, err := conn.Eval("box.execute([[START TRANSACTION;]]) box.execute([["+q+"]]) box.execute([[COMMIT;]])", []interface{}{})
	fmt.Println(resp, err)

	return id
}
