package mindmap

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"
)

// EventKind represents the type of mind-map mutation requested.
type EventKind int

const (
	EventInit EventKind = iota
	EventAppendTask
	EventAppendSubtask
	EventUpdateStatus
)

// Event encapsulates an action that mutates the daily mind-map.
type Event struct {
	Kind         EventKind
	TaskID       string
	ParentTaskID string
	Title        string
	NewStatus    string
	Date         time.Time
}

// MindMapUpdaterInterface defines the contract for mind map updates
type MindMapUpdaterInterface interface {
	Stop()
	InitializeDailyMindMap(date time.Time)
	AppendTask(taskID, title string)
	AppendSubtask(parentTaskID, subtaskID, title string)
	UpdateTaskStatus(taskID, newStatus string)
	UpdateSubtaskStatus(parentTaskID, taskID, newStatus string)
	ProcessedCount() int32
}

// Ensure MindMapUpdater implements MindMapUpdaterInterface
var _ MindMapUpdaterInterface = (*MindMapUpdater)(nil)

// MindMapUpdater receives task-related events and applies them to
// the corresponding daily mind-map file from a single background goroutine.
type MindMapUpdater struct {
	events chan Event
	done   chan struct{}
	wg     sync.WaitGroup

	rootDir    string
	retryLimit int

	// processedCount is used only for unit-test introspection.
	processedCount int32
	processedDates map[string]struct{}
}

// NewUpdater starts the updater's worker goroutine and returns the instance.
func NewUpdater(rootDir string) *MindMapUpdater {
	u := &MindMapUpdater{
		events:         make(chan Event, 256),
		done:           make(chan struct{}),
		rootDir:        rootDir,
		retryLimit:     1,
		processedDates: make(map[string]struct{}),
	}

	u.wg.Add(1)
	go u.loop()

	return u
}

// Stop gracefully shuts down the updater, waiting for pending events to finish.
func (u *MindMapUpdater) Stop() {
	close(u.done)
	u.wg.Wait()
}

// InitializeDailyMindMap enqueues an initialization event for the given date.
func (u *MindMapUpdater) InitializeDailyMindMap(date time.Time) {
	u.enqueue(Event{Kind: EventInit, Date: date})
}

// AppendTask enqueues a new task node under the Work section.
func (u *MindMapUpdater) AppendTask(taskID, title string) {
	// Clean the text before creating the task
	cleanTaskID := u.cleanTaskText(taskID)
	cleanTitle := u.cleanTaskText(title)

	u.enqueue(Event{
		Kind:   EventAppendTask,
		TaskID: cleanTaskID,
		Title:  cleanTitle,
	})
}

// AppendSubtask enqueues a new subtask node under its parent task.
func (u *MindMapUpdater) AppendSubtask(parentTaskID, subtaskID, title string) {
	// Clean the text before creating the subtask
	cleanParentID := u.cleanTaskText(parentTaskID)
	cleanSubtaskID := u.cleanTaskText(subtaskID)
	cleanTitle := u.cleanTaskText(title)

	u.enqueue(Event{
		Kind:         EventAppendSubtask,
		TaskID:       cleanSubtaskID,
		ParentTaskID: cleanParentID,
		Title:        cleanTitle,
	})
}

// UpdateTaskStatus enqueues a status update for an existing task node.
func (u *MindMapUpdater) UpdateTaskStatus(taskID, newStatus string) {
	u.enqueue(Event{Kind: EventUpdateStatus, TaskID: taskID, NewStatus: newStatus})
}

// UpdateSubtaskStatus enqueues a status update for an existing subtask node.
func (u *MindMapUpdater) UpdateSubtaskStatus(parentTaskID, taskID, newStatus string) {
	u.enqueue(Event{Kind: EventUpdateStatus, TaskID: taskID, ParentTaskID: parentTaskID, NewStatus: newStatus})
}

// enqueue pushes an event into the channel (non-blocking up to channel capacity).
func (u *MindMapUpdater) enqueue(ev Event) {
	log.Info("Enqueuing mind map event", "kind", ev.Kind, "taskID", ev.TaskID, "parentTaskID", ev.ParentTaskID, "title", ev.Title, "status", ev.NewStatus)
	select {
	case u.events <- ev:
		log.Info("Successfully enqueued mind map event")
	case <-u.done:
		log.Warn("Mind map updater shutting down, dropping event")
	}
}

func (u *MindMapUpdater) loop() {
	defer u.wg.Done()
	log.Info("Starting mind map updater loop")

	for {
		select {
		case ev := <-u.events:
			log.Info("Processing mind map event", "kind", ev.Kind, "taskID", ev.TaskID)
			u.handle(ev)
			atomic.AddInt32(&u.processedCount, 1)
		case <-u.done:
			log.Info("Mind map updater shutting down, draining events")
			// Drain remaining events before exiting.
			for ev := range u.events {
				u.handle(ev)
				atomic.AddInt32(&u.processedCount, 1)
			}
			return
		}
	}
}

// handle performs the actual file mutation logic.
func (u *MindMapUpdater) handle(ev Event) {
	switch ev.Kind {
	case EventInit:
		log.Info("Handling init event", "date", ev.Date)
		u.ensureDailyMindMap(ev.Date)
	case EventAppendTask:
		log.Info("Handling append task event", "taskID", ev.TaskID, "title", ev.Title)
		u.ensureDailyMindMap(u.resolveDate(ev))
		if err := u.appendTaskLine(ev); err != nil {
			log.Error("mindmap: append task", "err", err)
		}
	case EventAppendSubtask:
		log.Info("Handling append subtask event", "taskID", ev.TaskID, "parentTaskID", ev.ParentTaskID, "title", ev.Title)
		u.ensureDailyMindMap(u.resolveDate(ev))
		if err := u.appendSubtaskLine(ev); err != nil {
			log.Error("mindmap: append subtask", "err", err)
		}
	case EventUpdateStatus:
		log.Info("Handling update status event", "taskID", ev.TaskID, "status", ev.NewStatus)
		u.ensureDailyMindMap(u.resolveDate(ev))
		if err := u.updateTaskStatusLine(ev); err != nil {
			log.Error("mindmap: update status", "err", err)
		}
	default:
		log.Warn("Unknown event kind", "kind", ev.Kind)
	}
}

func (u *MindMapUpdater) ensureDailyMindMap(date time.Time) {
	// Always use today's date for new files
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	log.Info("Ensuring daily mind map exists", "date", dateStr)

	if _, done := u.processedDates[dateStr]; done {
		log.Info("Mind map already processed", "date", dateStr)
		return
	}

	// Ensure root dir exists
	log.Info("Creating root directory", "path", u.rootDir)
	if err := os.MkdirAll(u.rootDir, 0o755); err != nil {
		log.Error("mindmap: unable to create root dir", "err", err)
		return
	}

	path := filepath.Join(u.rootDir, dateStr+".txt")
	log.Info("Checking mind map file", "path", path)

	// If file already exists, mark as processed and return.
	if _, err := os.Stat(path); err == nil {
		log.Info("Mind map file already exists", "path", path)
		u.processedDates[dateStr] = struct{}{}
		return
	}

	contentLines := []string{
		dateStr,
		"\tPersonal",
		"\tWork",
		"", // trailing newline
	}

	data := []byte(strings.Join(contentLines, "\n"))
	log.Info("Creating new mind map file", "path", path)

	try := func() error {
		return os.WriteFile(path, data, 0o644)
	}

	if err := try(); err != nil {
		log.Error("First attempt to write mind map file failed", "err", err)
		// retry once
		if err2 := try(); err2 != nil {
			log.Error("mindmap: failed writing file", "path", path, "err", err2)
			return
		}
	}

	log.Info("Successfully created mind map file", "path", path)
	u.processedDates[dateStr] = struct{}{}
}

// ProcessedCount returns how many events were handled – useful for tests.
func (u *MindMapUpdater) ProcessedCount() int32 {
	return atomic.LoadInt32(&u.processedCount)
}

func (u *MindMapUpdater) resolveDate(ev Event) time.Time {
	// Always use today's date
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// cleanTaskText removes all status markers, icons, and dates from task text
func (u *MindMapUpdater) cleanTaskText(text string) string {
	// Remove status icons and dates
	statusRegex := regexp.MustCompile(`\s*[✅🛫⏳]\s+\d{4}-\d{2}-\d{2}$`)
	text = statusRegex.ReplaceAllString(text, "")

	// Remove any other icons that might be present
	text = strings.TrimSpace(text)
	return text
}

// findTaskLine looks for a task or subtask in the mind map lines and returns its index and whether it was found
func (u *MindMapUpdater) findTaskLine(lines []string, taskID string, tabLevel int) (int, bool) {
	cleanTaskID := u.cleanTaskText(taskID)

	for i, line := range lines {
		if leadingTabs(line) == tabLevel {
			cleanLine := u.cleanTaskText(line)
			if strings.Contains(cleanLine, cleanTaskID) {
				return i, true
			}
		}
	}

	return -1, false
}

func (u *MindMapUpdater) appendTaskLine(ev Event) error {
	date := u.resolveDate(ev)
	dateStr := date.Format("2006-01-02")
	path := filepath.Join(u.rootDir, dateStr+".txt")

	lines, err := u.readLines(path)
	if err != nil {
		return err
	}

	// Find Work section index
	workIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "Work" && leadingTabs(line) == 1 {
			workIdx = i
			break
		}
	}
	if workIdx == -1 {
		return nil // Work section missing
	}

	// Determine insertion index after last \t\t line following workIdx
	insertIdx := workIdx + 1
	for insertIdx < len(lines) && leadingTabs(lines[insertIdx]) == 2 {
		insertIdx++
	}

	// Clean the task title before adding
	cleanTitle := u.cleanTaskText(ev.Title)
	newLine := "\t\t" + cleanTitle

	// Avoid duplicate
	for _, l := range lines {
		if u.cleanTaskText(l) == u.cleanTaskText(newLine) {
			return nil
		}
	}

	lines = append(lines[:insertIdx], append([]string{newLine}, lines[insertIdx:]...)...)
	return u.writeLines(path, lines)
}

func (u *MindMapUpdater) appendSubtaskLine(ev Event) error {
	date := u.resolveDate(ev)
	dateStr := date.Format("2006-01-02")
	path := filepath.Join(u.rootDir, dateStr+".txt")

	lines, err := u.readLines(path)
	if err != nil {
		return err
	}

	// find parent task line index
	parentIdx := -1
	cleanParentID := u.cleanTaskText(ev.ParentTaskID)
	for i, line := range lines {
		if leadingTabs(line) == 2 && u.cleanTaskText(line) == cleanParentID {
			parentIdx = i
			break
		}
	}
	if parentIdx == -1 {
		return nil // parent not found
	}

	// find insertion index after existing subtasks
	insertIdx := parentIdx + 1
	for insertIdx < len(lines) && leadingTabs(lines[insertIdx]) == 3 {
		insertIdx++
	}

	// Clean the task title before adding
	cleanTitle := u.cleanTaskText(ev.Title)
	newLine := "\t\t\t" + cleanTitle

	// Avoid duplicate
	for _, l := range lines {
		if u.cleanTaskText(l) == u.cleanTaskText(newLine) {
			return nil
		}
	}

	lines = append(lines[:insertIdx], append([]string{newLine}, lines[insertIdx:]...)...)
	return u.writeLines(path, lines)
}

func (u *MindMapUpdater) readLines(path string) ([]string, error) {
	log.Info("Reading mind map file", "path", path)
	f, err := os.Open(path)
	if err != nil {
		log.Error("Failed to open mind map file", "path", path, "err", err)
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Error("Failed to scan mind map file", "path", path, "err", err)
		return nil, err
	}
	log.Info("Successfully read mind map file", "path", path, "lineCount", len(lines))
	return lines, nil
}

func (u *MindMapUpdater) writeLines(path string, lines []string) error {
	log.Info("Writing mind map file", "path", path, "lineCount", len(lines))
	data := []byte(strings.Join(lines, "\n"))
	try := func() error {
		return os.WriteFile(path, append(data, '\n'), 0o644)
	}
	if err := try(); err != nil {
		log.Error("First attempt to write mind map file failed", "err", err)
		if err2 := try(); err2 != nil {
			log.Error("Failed to write mind map file", "path", path, "err", err2)
			return err2
		}
	}
	log.Info("Successfully wrote mind map file", "path", path)
	return nil
}

func leadingTabs(s string) int {
	count := 0
	for _, ch := range s {
		if ch == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}

// cleanTaskID removes status markers and dates from a task ID
func (u *MindMapUpdater) cleanTaskID(taskID string) string {
	statusRegex := regexp.MustCompile(`\s*[✅🛫⏳]\s+\d{4}-\d{2}-\d{2}$`)
	return statusRegex.ReplaceAllString(strings.TrimSpace(taskID), "")
}

// findAndUpdateTask looks for a task/subtask and updates its status if found
func (u *MindMapUpdater) findAndUpdateTask(lines []string, taskID string, status string) ([]string, bool) {
	// Try as subtask first, then as main task
	for _, tabLevel := range []int{3, 2} {
		if idx, found := u.findTaskLine(lines, taskID, tabLevel); found {
			log.Info("Found task in mind map", "line", lines[idx], "tabLevel", tabLevel)
			// Keep only the clean task text with proper indentation
			cleanText := u.cleanTaskText(lines[idx])
			prefix := strings.Repeat("\t", tabLevel)
			lines[idx] = prefix + cleanText
			return lines, true
		}
	}
	return lines, false
}

// createAndUpdateTask creates a new task/subtask with clean text
func (u *MindMapUpdater) createAndUpdateTask(path string, taskID string, parentTaskID string, status string) error {
	// Create parent task if this is a subtask
	if parentTaskID != "" {
		if err := u.appendTaskLine(Event{Kind: EventAppendTask, TaskID: parentTaskID, Title: parentTaskID}); err != nil {
			return err
		}
		if err := u.appendSubtaskLine(Event{Kind: EventAppendSubtask, TaskID: taskID, ParentTaskID: parentTaskID, Title: taskID}); err != nil {
			return err
		}
	} else {
		if err := u.appendTaskLine(Event{Kind: EventAppendTask, TaskID: taskID, Title: taskID}); err != nil {
			return err
		}
	}

	return nil
}

// removeTaskLine removes a task or subtask from the mind map
func (u *MindMapUpdater) removeTaskLine(lines []string, taskID string) ([]string, bool) {
	// Try as subtask first, then as main task
	for _, tabLevel := range []int{3, 2} {
		if idx, found := u.findTaskLine(lines, taskID, tabLevel); found {
			log.Info("Removing task from mind map", "line", lines[idx], "tabLevel", tabLevel)
			// If this is a main task, also remove all its subtasks
			if tabLevel == 2 {
				// Find the end of subtasks section
				endIdx := idx + 1
				for endIdx < len(lines) && leadingTabs(lines[endIdx]) == 3 {
					endIdx++
				}
				// Remove the task and all its subtasks
				lines = append(lines[:idx], lines[endIdx:]...)
			} else {
				// Just remove the single subtask
				lines = append(lines[:idx], lines[idx+1:]...)
			}
			return lines, true
		}
	}
	return lines, false
}

func (u *MindMapUpdater) updateTaskStatusLine(ev Event) error {
	path := filepath.Join(u.rootDir, time.Now().Format("2006-01-02")+".txt")
	lines, err := u.readLines(path)
	if err != nil {
		return err
	}

	taskID := u.cleanTaskID(ev.TaskID)
	log.Info("Looking for task in mind map", "taskID", taskID)

	// If task is unscheduled, remove it from the mind map
	if ev.NewStatus == "unscheduled" {
		if lines, removed := u.removeTaskLine(lines, taskID); removed {
			return u.writeLines(path, lines)
		}
		return nil
	}

	// For other statuses, ensure task exists with clean text
	if lines, found := u.findAndUpdateTask(lines, taskID, ev.NewStatus); found {
		return u.writeLines(path, lines)
	}

	// Task not found and not unscheduled, create it with clean text
	return u.createAndUpdateTask(path, taskID, ev.ParentTaskID, ev.NewStatus)
}

func statusIcon(status string) string {
	switch status {
	case "scheduled":
		return "⏳"
	case "started":
		return "🛫"
	case "completed":
		return "✅"
	default:
		return ""
	}
}
