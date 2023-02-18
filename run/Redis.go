package run

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"golin/config"
	"os"
	"strings"
	"time"
)

func Redis(cmd *cobra.Command, args []string) {
	//获取分隔符，默认是||
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "Redis")
		wg.Wait()
		config.Log.Info("单次运行Redis模式完成！")
		return
	}
	//下面是多线程的模式
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断redis.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), pem, ippath)

	// 运行share文件中的函数
	Rangefile(ippath, spr, "Redis")
	wg.Wait()
	//完成前最后写入文件
	Deffile("Redis", count, count-len(errhost), errhost)

}

func Runredis(myname, myhost, mypasswd, myport1 string) {
	defer wg.Done()
	Port := strings.Replace(myport1, "\r", "", -1)
	ctx := context.Background()
	adr := myhost + ":" + Port
	client := redis.NewClient(&redis.Options{
		Addr:            adr,
		Password:        mypasswd,
		DB:              0,
		DialTimeout:     1 * time.Second,
		MinRetryBackoff: 1 * time.Second,
		ReadTimeout:     1 * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		//wg.Done()
		errhost = append(errhost, myhost)
		return
	}
	client.Get(ctx, "config").Val()
	ipaddr := client.ConfigGet(ctx, "bind").Val()
	lofile := client.ConfigGet(ctx, "logfile").Val()
	loglevel := client.ConfigGet(ctx, "loglevel").Val()
	pass := client.ConfigGet(ctx, "requirepass").Val()
	redistimout := client.ConfigGet(ctx, "timeout").Val()
	redisport := client.ConfigGet(ctx, "port").Val()
	redisdir := client.ConfigGet(ctx, "dir").Val()
	//confinfo := client.Info(ctx).Val()

	_, err = os.Stat(succpath)
	if os.IsNotExist(err) {
		os.Mkdir(succpath, pem)
	}
	fire := "采集完成目录//" + myname + "_" + myhost + "(redis).log"
	file, _ := os.OpenFile(fire, os.O_CREATE|os.O_APPEND, pem)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("-----基本信息------\n")
	write.WriteString(fmt.Sprintf("地址限制策略为:%s\n", ipaddr[1]))
	write.WriteString(fmt.Sprintf("日志存储为:%s  日志等级为:%s\n", lofile[1], loglevel[1]))
	write.WriteString(fmt.Sprintf("密码信息为:%s\n", pass[1]))
	write.WriteString(fmt.Sprintf("超时时间为:%s\n", redistimout[1]))
	write.WriteString(fmt.Sprintf("redis运行端口为:%s\n", redisport[1]))
	write.WriteString(fmt.Sprintf("redis运行位置为:%s\n", redisdir[1]))
	write.WriteString("\n-----info信息------\n")
	write.WriteString(client.Info(ctx).Val())
	write.Flush()
}
