package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/voioc/geo/cz"
	"github.com/voioc/geo/lite"
	"gopkg.in/natefinch/lumberjack.v2"
)

var czGeo *cz.CZ
var liteGeo *lite.Geo

func init() {
	viper.SetConfigFile("./config.toml") // 指定配置文件路径
	err := viper.ReadInConfig()          // 查找并读取配置文件
	if err != nil {                      // 处理读取配置文件的错误
		fmt.Printf("Fatal error config file: %s \n", err.Error())
	}

	czGeo = &cz.CZ{}
	liteGeo = &lite.Geo{}

}

func main() {
	path := viper.GetString("log.path")
	fmt.Println(path)

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err.Error())
		os.Exit(0)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Fatal error watcher: %s \n", err.Error())
		os.Exit(0)
	}

	if err = watcher.Add(path); err != nil {
		fmt.Printf("Fatal error add watcher: %s \n", err.Error())
		os.Exit(0)
	}

	log := &lumberjack.Logger{
		Filename:   viper.GetString("log.target"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	defer func() {
		file.Close()
		watcher.Close()
		log.Close()
	}()

	buf := bufio.NewReader(file)
	for {
		// n, err := f.File.Read(b)
		line, err := buf.ReadBytes('\n')
		if err == io.EOF { // 现有文件已经读取完毕
			select {
			case _, ok := <-watcher.Events: // 事件到来时唤醒，重新尝试读
				if !ok { // 通道已关闭
					fmt.Println("Events chan got: ", io.EOF.Error())
				}
			case err, ok := <-watcher.Errors:
				if !ok { // 通道已关闭
					fmt.Println("Reads chan error: ", io.EOF.Error())
				}

				if err != nil {
					fmt.Println("11111", err.Error())
				}
			}
		} else {
			Loc(line)
		}
	}
}

func Loc(log []byte) string {
	ip := jsoniter.Get(log, "remote_addr").ToString()

	// var location *model.Location
	// source := viper.GetString("geo.source")
	// if source == "cz" {
	// 	location = czGeo.Analyze(ip)
	// } else if source == "geo" {
	// 	location, _ = liteGeo.Analyze(ip)
	// } else { // source == "all"
	// 	czLoc := czGeo.Analyze(ip)
	// 	geoLoc, _ := liteGeo.Analyze(ip)
	// 	primary := viper.GetString("geo.source")

	// 	if primary == "cz" {
	// 		location = geoLoc
	// 		if czLoc != nil && czLoc.Province != "" {
	// 			location = czLoc
	// 		}
	// 	} else {
	// 		location = czLoc
	// 		if geoLoc != nil {
	// 			location = geoLoc
	// 		}
	// 	}
	// }

	fmt.Printf("%+v\n", location)
	return ""
}

// var version = "NaLi-Go 1.5.0\n" +
// 	"Source: https://github.com/Mikubill/nali-go\n" +
// 	"Git Commit Hash: %s\n"

// func main() {
// 	info, err := os.Stdin.Stat()
// 	if err != nil {
// 		panic(err)
// 	}

// 	if args := os.Args; len(args) > 1 {
// 		ip := cz.CZ{}
// 		ip.Cmd()
// 		for i := range args {
// 			item := args[i]
// 			if !strings.Contains(item, os.Args[0]) {
// 				fmt.Println(ip.Analyze(item))
// 			}
// 		}

// 		// fmt.Printf("\n")
// 		os.Exit(0)
// 	}

// 	fmt.Println(info.Mode(), os.ModeCharDevice)
// 	if (info.Mode() & os.ModeCharDevice) != 0 {
// 		self := os.Args[0]
// 		fmt.Println("1111", cz.Helper, "2222", self)
// 		os.Exit(0)
// 	}

// 	reader := bufio.NewReader(os.Stdin)
// 	ip := cz.CZ{}
// 	for {
// 		line, err := reader.ReadString('\n')
// 		fmt.Printf("%s", ip.Analyze(line))
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			} else {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 		}
// 	}
// 	// fmt.Printf("\n")
// }

// func contains(array []string, flag string) bool {
// 	for i := 0; i < len(array); i++ {
// 		if array[i] == flag {
// 			return true
// 		}
// 	}
// 	return false
// }

// func execute(cmd string) {
// 	runner := exec.Command("sh", "-c", cmd)
// 	fmt.Println(runner.Args)
// 	stdout, err := runner.StdoutPipe()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	_ = runner.Start()
// 	reader := bufio.NewReader(stdout)
// 	for {
// 		line, err := reader.ReadString('\n')
// 		fmt.Printf("%s", line)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			} else {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 		}
// 	}
// }
