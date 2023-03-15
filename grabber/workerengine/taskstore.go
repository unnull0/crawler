package workerengine

import (
	"github.com/unnull0/crawler/grabber"
	"github.com/unnull0/crawler/tasklib/doubantenement"
)

func init() {
	Tkstore.Add(doubantenement.DoubantenementTask)
}

var Tkstore = &TaskStore{
	list: []*grabber.Task{},
	hash: map[string]*grabber.Task{},
}

type TaskStore struct {
	list []*grabber.Task
	hash map[string]*grabber.Task
}

func (t *TaskStore) Add(task *grabber.Task) {
	t.hash[task.Name] = task
	t.list = append(t.list, task)
}
