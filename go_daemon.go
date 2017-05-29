// プロセスを daemon として起動する
// 基本的に下記のソースコードのパクりであるｗ
// https://stackoverflow.com/questions/23031752/start-a-process-in-go-and-detach-from-it

package main

import (
	"fmt"
	"log/syslog"
	"os"
	"os/user"
	"strconv"
	"syscall" )

const プロセス名 = "foo"
const フルパス名 = "/usr/local/foo/" + プロセス名

func main() {
	システムログ, err := syslog.New(syslog.LOG_NOTICE |syslog.LOG_USER, プロセス名 )
	if err == nil {
// ここでユーザ名を指定。ユーザ名は root も使用できる。
		ユーザ, err := user.Lookup( "daemon" )
		if err == nil {
			uid,_ := strconv.ParseInt(ユーザ.Uid, 10, 32)
			gid,_ := strconv.ParseInt(ユーザ.Gid, 10, 32)
			var 認証情報 = &syscall.Credential{ uint32(uid), uint32(gid), []uint32{} }
// Noctty フラグは、子プロセスを親プロセスの tty から切り離すために使用する。
			var システムプロセス属性 = &syscall.SysProcAttr{ Credential:認証情報, Noctty:true }
			var プロセス属性 = os.ProcAttr{
				Dir: ".",
				Env: os.Environ(),
// os.File には標準入出力のリダイレクト先を指定している。daemon なので何も指定しない。
				Files: []*os.File{ nil, nil, nil, },
				Sys:システムプロセス属性,
			}
// '[]string{ フルパス名 }' は、'[]string{ フルパス名, パラメータ1 }' のようにパラメータを含めることもできる。
			子プロセス, err := os.StartProcess(フルパス名, []string{ フルパス名 }, &プロセス属性)
			if err == nil {
				fmt.Printf( "%d", 子プロセス.Pid )
				err = システムログ.Info( fmt.Sprintf( "PID %d / UID %d / GID %d で起動しますた", 子プロセス.Pid, uid, gid ))
				if err == nil {
// ドキュメントにははっきり書かれていないが、実際 Realease() は子プロセスの切り離し（fork()）を行う。
					err = 子プロセス.Release();
// ここの時点では 子プロセス.Pid の値は -1 になっている。
					if err != nil {
						システムログ.Emerg( fmt.Sprintf( " は、子プロセスの Release() に失敗しますた… / 理由 = %v", err))
					}
				} else {
					fmt.Printf( "システムログに書き込めンゴ… / 理由 = %v", err)
				}
			} else {
// 実行ファイルが無い・権限不足と言った一般的な要因で起こるエラーはここに落ちる。
				システムログ.Emerg( fmt.Sprintf( "子プロセスの起動に失敗したンゴ… / 理由 = %v", err))
			}
		} else {
			システムログ.Emerg( fmt.Sprintf( "ユーザ情報の取得に失敗したンゴ… / 理由 = %v", err))
		}
	} else {
		fmt.Printf( "システムログが開けンゴ… / 理由 = %v", err)
	}
}