package log

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func init() {
	conf := Config{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "time",
		Encoding:   "json",
		CallerKey:  "caller",
		Level:      -1,
	}
	X, _ = GetLoggerByConf(&conf)
}

func Init(configPath string) (err error) {
	logger, err := getLogger(configPath)
	if err != nil {
		return
	}
	X = logger
	return
}

// getLogger get X
func getLogger(path string) (logger *Logger, err error) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("read log configuration %v, error: %v", path, err)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		err = fmt.Errorf("unmarshal log configuration %v, error: %v", path, err)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	return GetLoggerByConf(&config)

}
