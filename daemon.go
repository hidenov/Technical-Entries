// daemon 本体
package main

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const プロセス名 = "foo"
const 動作ログファイル名 = "./foo.log"

func main () {
	システムログ, err := syslog.New(syslog.LOG_NOTICE |syslog.LOG_USER, プロセス名 )
	if err == nil {
		システムログ.Info( "処理開始" )
		動作ログ, err := os.OpenFile(動作ログファイル名, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(io.MultiWriter(動作ログ, os.Stdout))
			log.SetFlags(log.Ldate | log.Ltime)
			log.Println( "情報:daemon 処理開始" )
// シグナルの取扱い・シグナル受信用の channel を宣言する
			ch_シグナル受信 := make(chan os.Signal, 1)
// シグナルの取扱い・受信するシグナルの種類を宣言する。
// daemon であれば SIGTERM だけで十分なんだけど、開発時にターミナルから CTRL+C で
// 止めることも考慮して、SIGINT も足してある。
			signal.Notify( ch_シグナル受信, syscall.SIGINT, syscall.SIGTERM )
			i := 0
			f_処理終了 := false
			for f_処理終了 == false {
// シグナルの取扱い・シグナルを受信したかどうかの確認
// select { ～ case シグナル　までがひと纏まりで
// シグナルを受信していなければ、default に抜ける。
				select {
					case シグナル, _ := <-ch_シグナル受信 :
						システムログ.Info( fmt.Sprintf( "処理終了要求を受信 / シグナル番号 = %v", シグナル ))
						f_処理終了 = true
					default :
						if i % 100 == 0 { // １０秒に１回処理を起動する
							処理実行( )
							i = 1
						} else {
							i++
						}
// タスクスケジューリング・一番簡単な 100ミリ秒処理を休むパターン
						time.Sleep( 100 * time.Millisecond )
				}
			}
		} else {
			システムログ.Emerg( fmt.Sprintf( "logの生成に失敗 / 理由 = %v", err ))
		}
		defer 動作ログ.Close()
	}
}

func 処理実行() {
	log.Println( "情報:daemon 処理実行" )
}

