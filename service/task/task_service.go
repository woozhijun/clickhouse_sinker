package task

import (
	"fmt"
	"github.com/woozhijun/clickhouse_sinker/model"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/woozhijun/clickhouse_sinker/service/input"
	"github.com/woozhijun/clickhouse_sinker/service/output"
	"github.com/woozhijun/clickhouse_sinker/service/parser"

	"github.com/wswz/go_commons/log"
)

type TaskService struct {
	stopped    chan struct{}
	kafka      *input.Kafka
	clickhouse *output.ClickHouse
	p          parser.Parser
	sync.Mutex

	FlushInterval int
	BufferSize    int
	ItemSchemas []*model.ItemSchema
}

func NewTaskService(kafka *input.Kafka, clickhouse *output.ClickHouse, p parser.Parser) *TaskService {
	return &TaskService{
		stopped:    make(chan struct{}),
		kafka:      kafka,
		clickhouse: clickhouse,
		p:          p,
	}
}

func (service *TaskService) Init() error {
	err := service.kafka.Init()
	if err != nil {
		return err
	}
	return service.clickhouse.Init()
}

func (service *TaskService) Run() {
	if err := service.kafka.Start(); err != nil {
		panic(err)
	}

	log.Infof("TaskService %s TaskService has started", service.clickhouse.GetName())
	tick := time.NewTicker(time.Duration(service.FlushInterval) * time.Second)
	defer tick.Stop()

	msgs := make([]parser.IParser, 0, 100000)
FOR:
	for {
		select {
		case msg, more := <-service.kafka.Msgs():
			if !more {
				break FOR
			}
			msgs = append(msgs, service.parse(msg))
			if len(msgs) >= service.BufferSize {
				service.Lock()
				service.flush(msgs)
				msgs = make([]parser.IParser, 0, 100000)
				tick = time.NewTicker(time.Duration(service.FlushInterval) * time.Second)
				service.Unlock()
			}
		case <-tick.C:
			log.Info(service.clickhouse.GetName() + " tick")
			if len(msgs) == 0 {
				continue
			}
			service.Lock()
			service.flush(msgs)
			msgs = make([]parser.IParser, 0, 100000)
			service.Unlock()
		}
	}
	service.flush(msgs)
	service.stopped <- struct{}{}
	return
}

func (service *TaskService) parse(data []byte) parser.IParser {
	return service.p.Parse(data)
}
func (service *TaskService) flush(metrics []parser.IParser) {
	log.Info("buf size:", len(metrics))
	service.clickhouse.LoopWrite(metrics)
}

func (service *TaskService) Stop() {
	log.Info("close TaskService size:")
	if err := service.kafka.Stop(); err != nil {
		panic(err)
	}
	<-service.stopped
	service.clickhouse.Close()
	log.Info("closed TaskService size:")
}

//获取goroutine的id
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}