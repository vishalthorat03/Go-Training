package utils

type Task struct {
	Action string
	Data   map[string]string
}

var taskChannel = make(chan Task, 10)

// StartWorker starts a goroutine to process tasks
func StartWorker() {
	go func() {
		for task := range taskChannel {
			switch task.Action {
			case "add":
				AddEntry(task.Data)
			case "delete":
				DeleteEntry(task.Data["key"], task.Data["value"])
			}
		}
	}()
}

// AddTask adds a task to the channel
func AddTask(action string, data map[string]string) {
	taskChannel <- Task{Action: action, Data: data}
}
