package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"scrm-sync-data/app/config"
	"sync"
	"time"
)

type TaskFunc func(args interface{})

type Task struct {
	loopNum int
	exec    TaskFunc
	params  interface{}
}

type DelayTaskQueue struct {
	lock      *sync.RWMutex //防止map读写冲突
	curIndex  int
	slots     [3600]map[string]*Task
	closed    chan bool
	taskRange chan bool
	taskClose chan bool
	timeClose chan bool
	startTime time.Time
}

type DelayTaskQueueService struct {
	lock    *sync.Mutex
	context context.Context
	cancel  context.CancelFunc
	config  *config.Config
	logger  *LoggerService
	queues  map[string]*DelayTaskQueue
	State   bool
}

func (delayTaskQueueService *DelayTaskQueueService) Init(parentContext context.Context, config *config.Config, logger *LoggerService) {
	delayTaskQueueService.config = config
	delayTaskQueueService.logger = logger
	delayTaskQueueService.State = false
	delayTaskQueueService.context, delayTaskQueueService.cancel = context.WithCancel(parentContext)
	delayTaskQueueService.Start()
}

func (delayTaskQueueService *DelayTaskQueueService) Start() {
	fmt.Println("start delay queue service...", time.Now())
	delayTaskQueueService.logger.Debug(fmt.Sprintf("%s,%v", "start delay queue service...", time.Now()))
	delayTaskQueueService.lock = &sync.Mutex{}
	delayTaskQueueService.queues = make(map[string]*DelayTaskQueue)
	delayTaskQueueService.State = true
}

func (delayTaskQueueService *DelayTaskQueueService) Stop() {
	delayTaskQueueService.lock.Lock()
	defer delayTaskQueueService.lock.Unlock()
	fmt.Println("stop delay queue service...", time.Now())
	delayTaskQueueService.logger.Debug(fmt.Sprintf("%s,%v", "stop delay queue service...", time.Now()))
	for _, delayTaskQueue := range delayTaskQueueService.queues {
		delayTaskQueue.Stop()
	}
	delayTaskQueueService.queues = make(map[string]*DelayTaskQueue)
	delayTaskQueueService.cancel()
	delayTaskQueueService.State = false
}

func (delayTaskQueueService *DelayTaskQueueService) NewDelayTaskQueue(queueName string) (delayTaskQueue *DelayTaskQueue) {
	delayTaskQueueService.lock.Lock()
	defer delayTaskQueueService.lock.Unlock()
	dq := &DelayTaskQueue{
		lock:      &sync.RWMutex{},
		curIndex:  0,
		closed:    make(chan bool),
		taskRange: make(chan bool),
		taskClose: make(chan bool),
		timeClose: make(chan bool),
		startTime: time.Now(),
	}
	for i := 0; i < 3600; i++ {
		dq.slots[i] = make(map[string]*Task)
	}
	go dq.Start()
	delayTaskQueueService.queues[queueName] = dq
	return dq
}

func (delayTaskQueueService *DelayTaskQueueService) DelayTaskQueue(queueName string) (delayTaskQueue *DelayTaskQueue, ok bool) {
	delayTaskQueue, ok = delayTaskQueueService.queues[queueName]
	return delayTaskQueue, ok
}

func (dq *DelayTaskQueue) Start() {
	go dq.TaskLoop()
	go dq.TimeLoop()
	for {
		select {
		case <-dq.closed:
			{
				dq.taskClose <- true
				dq.timeClose <- true
				break
			}
		}
	}
}

func (dq *DelayTaskQueue) Stop() {
	dq.closed <- true
}

func (dq *DelayTaskQueue) TaskLoop() {
	defer func() { fmt.Println("task loop exit") }()
	for {
		select {
		case <-dq.taskClose:
			break
		case <-dq.taskRange:
			dq.lock.RLock()
			tasks := dq.slots[dq.curIndex]
			if len(tasks) > 0 {
				for index, task := range tasks {
					if task.loopNum == 0 {
						go task.exec(task.params)
						delete(tasks, index)
					} else {
						task.loopNum--
					}
				}
			}
			dq.lock.RUnlock()
		}
	}
}

func (dq *DelayTaskQueue) TimeLoop() {
	defer func() { fmt.Println("time loop exit") }()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-dq.timeClose:
			break
		case <-ticker.C:
			if dq.curIndex == 3599 {
				dq.curIndex = 0
			} else {
				dq.curIndex++
			}
			dq.taskRange <- true
			break
		}
	}
}

func (dq *DelayTaskQueue) AddTask(t time.Time, key string, exec TaskFunc, params interface{}) error {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	if dq.startTime.After(t) {
		return errors.New("task time error")
	}
	subSecond := t.Unix() - dq.startTime.Unix()
	loopNum := int(subSecond / 3600)
	index := subSecond % 3600
	tasks := dq.slots[index]
	if _, ok := tasks[key]; ok {
		return errors.New("task exists in queue")
	}
	tasks[key] = &Task{
		loopNum: loopNum,
		exec:    exec,
		params:  params,
	}
	return nil
}

func (dq *DelayTaskQueue) DelTask(key string) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	for index := range dq.slots {
		if len(dq.slots[index]) == 0 {
			continue
		}
		tasks := dq.slots[index]
		if _, ok := tasks[key]; ok {
			delete(dq.slots[index], key)
		}
	}
}
