package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

func esperarNotificacao(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Dados recebidos do canal[", n.Channel, "]:")

			//// Prepara a carga de notificação para impressão bonita
			//var prettyJson bytes.Buffer
			//err := json.Indent(&prettyJson, []byte(n.Extra), "", "\t")
			//if err != nil {
			//	fmt.Println("Erro ao processar JSON", err)
			//	return
			//}
			//fmt.Println(string(prettyJson.Bytes()))
			fmt.Println(n.Extra)
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Não recebeu nenhum evento por 90 segundos, verificando a conexão")
			go func() {
				l.Ping()
			}()
			return

		}
	}

}

func main() {
	var conninfo string = "dbname=postgres user=postgres password=dev123456 sslmode=disable"

	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Monitorando o PostgreSQL")
	for {
		esperarNotificacao(listener)
	}
}
