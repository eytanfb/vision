package mindmap

import (
	"time"

	"github.com/charmbracelet/log"
)

// NullMindMapUpdater implements MindMapUpdaterInterface but does nothing
type NullMindMapUpdater struct{}

func NewNullUpdater() MindMapUpdaterInterface {
	log.Info("Creating new NullMindMapUpdater")
	return &NullMindMapUpdater{}
}

func (n *NullMindMapUpdater) Stop() {}

func (n *NullMindMapUpdater) InitializeDailyMindMap(date time.Time) {}

func (n *NullMindMapUpdater) AppendTask(taskID, title string) {
	log.Debug("NullMindMapUpdater: Ignoring AppendTask", "taskID", taskID, "title", title)
}

func (n *NullMindMapUpdater) AppendSubtask(parentTaskID, subtaskID, title string) {
	log.Debug("NullMindMapUpdater: Ignoring AppendSubtask", "parentTaskID", parentTaskID, "subtaskID", subtaskID, "title", title)
}

func (n *NullMindMapUpdater) UpdateTaskStatus(taskID, newStatus string) {
	log.Debug("NullMindMapUpdater: Ignoring UpdateTaskStatus", "taskID", taskID, "newStatus", newStatus)
}

func (n *NullMindMapUpdater) UpdateSubtaskStatus(parentTaskID, taskID, newStatus string) {
	log.Debug("NullMindMapUpdater: Ignoring UpdateSubtaskStatus", "parentTaskID", parentTaskID, "taskID", taskID, "newStatus", newStatus)
}

func (n *NullMindMapUpdater) ProcessedCount() int32 {
	return 0
}
