package creator

import (
	"encoding/json"
	"github.com/woozhijun/clickhouse_sinker/model"
	"github.com/woozhijun/clickhouse_sinker/service/input"
	"github.com/woozhijun/clickhouse_sinker/service/output"
	"github.com/woozhijun/clickhouse_sinker/service/parser/impl"
	"github.com/woozhijun/clickhouse_sinker/service/task"
	"github.com/woozhijun/clickhouse_sinker/util"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/wswz/go_commons/log"
	"github.com/wswz/go_commons/utils"
)

type Config struct {
	Kafka      map[string]*KafkaConfig
	Clickhouse map[string]*ClickHouseConfig

	Tasks []*model.Task

	Common struct {
		FlushInterval int
		BufferSize    int
		LogLevel      string
	}
}

type KafkaConfig struct {
	Brokers string
	Sasl    struct {
		Password string
		Username string
	}
	Version string
}

type ClickHouseConfig struct {
	Db   string
	Host string
	Port int

	Username    string
	Password    string
	MaxLifeTime int
	DsnParams   string
	DnsLoop     bool
}

var (
	defaultFlushInterval = 3
	defaultBufferSize    = 10000
)

var (
	baseConfig *Config
)

// InitConfig must run before the server start
func InitConfig(dir string) *Config {
	confPath := ""
	if len(dir) > 0 {
		confPath = dir
	}
	var f = "config.json"
	f = filepath.Join(confPath, f)
	s, err := utils.ExtendFile(f)
	if err != nil {
		panic(err)
	}
	baseConfig = &Config{}
	err = json.Unmarshal([]byte(s), baseConfig)
	if err != nil {
		panic(err)
	}
	if baseConfig.Common.FlushInterval < 1 {
		baseConfig.Common.FlushInterval = defaultFlushInterval
	}

	if baseConfig.Common.BufferSize < 1 {
		baseConfig.Common.BufferSize = defaultBufferSize
	}
	err = baseConfig.LoadTasks(filepath.Join(confPath, "tasks"))
	if err != nil {
		panic(err)
	}

	log.SetLevelStr(baseConfig.Common.LogLevel)
	pp.Println(baseConfig)
	return baseConfig
}

func (cfg *Config) LoadTasks(dir string) error {
	//检测配置是否正确
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	cfg.Tasks = make([]*model.Task, 0, len(files))
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			s, err := utils.ExtendFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return err
			}
			taskConfig := &model.Task{}
			err = json.Unmarshal([]byte(s), taskConfig)
			if err != nil {
				return err
			}
			cfg.Tasks = append(cfg.Tasks, taskConfig)
		}
	}
	return nil
}

// GenTasks generate the tasks via config
func (config *Config) GenTasks() []*task.TaskService {
	res := make([]*task.TaskService, 0, len(config.Tasks))
	for _, taskConfig := range config.Tasks {
		kafka := config.GenInput(taskConfig)
		ck := config.GenOutput(taskConfig)
		p := impl.NewParser(taskConfig.Parser)

		taskImpl := task.NewTaskService(kafka, ck, p)

		util.IngestConfig(taskConfig, taskImpl)

		if taskImpl.FlushInterval == 0 {
			taskImpl.FlushInterval = config.Common.FlushInterval
		}

		if taskImpl.BufferSize == 0 {
			taskImpl.BufferSize = config.Common.BufferSize
		}

		res = append(res, taskImpl)
	}
	return res
}

// GenInput generate the input via config
func (config *Config) GenInput(taskCfg *model.Task) *input.Kafka {
	kfkCfg := config.Kafka[taskCfg.Kafka]

	inputImpl := input.NewKafka()
	util.IngestConfig(taskCfg, inputImpl)
	util.IngestConfig(kfkCfg, inputImpl)
	return inputImpl
}

// GenOutput generate the output via config
func (config *Config) GenOutput(taskCfg *model.Task) *output.ClickHouse {
	ckCfg := config.Clickhouse[taskCfg.Clickhouse]

	outputImpl := output.NewClickHouse()

	util.IngestConfig(ckCfg, outputImpl)
	util.IngestConfig(taskCfg, outputImpl)
	util.IngestConfig(config.Common, outputImpl)
	return outputImpl
}