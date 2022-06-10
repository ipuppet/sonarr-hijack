package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ultagic.com/pkg/flags"
	"ultagic.com/pkg/utils"
)

var (
	logger *log.Logger
)

func init() {
	logger = utils.Logger("config")
}

type Config struct {
	Filename       string
	_path          string
	data           map[string]interface{}
	lastModifyTime int64
	rwLock         sync.RWMutex
	notifyList     []Notifyer
}

func (c *Config) Init() {
	m, err := c.parse()
	if err != nil {
		logger.Println(err)
		return
	}
	c.data = m
	go c.reload()
}

func (c *Config) path() string {
	if c._path == "" {
		c._path = filepath.Join(flags.ConfigPath, c.Filename)
	}

	return c._path
}

func (c *Config) AddNotifyer(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}

func (c *Config) parse() (m map[string]interface{}, err error) {
	m = make(map[string]interface{}, 50)
	configString, err := ioutil.ReadFile(c.path())
	if err != nil {
		return
	}

	if err = json.Unmarshal(configString, &m); err != nil {
		return
	}

	return
}

func (c *Config) reload() {
	// 每 5 秒重新加载一次配置文件
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		func() {
			file, err := os.Open(c.path())
			if err != nil {
				logger.Printf("open %s failed, err: %v\n", c.Filename, err)
				return
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				logger.Printf("stat %s failed, err: %v\n", c.Filename, err)
				return
			}

			curModifyTime := fileInfo.ModTime().Unix()

			// 判断文件的修改时间是否大于最后一次修改时间
			if curModifyTime > c.lastModifyTime {
				m, err := c.parse()
				if err != nil {
					logger.Println("parse failed, err: ", err)
					return
				}

				c.rwLock.Lock()
				c.data = m
				c.rwLock.Unlock()

				c.lastModifyTime = curModifyTime

				for _, n := range c.notifyList {
					n.Callback(c)
				}
			}
		}()
	}
}

func (c *Config) Get(key string) (value interface{}, err error) {
	// 根据字符串获取
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		err = errors.New("key[" + key + "] not found")
		return
	}
	return
}
