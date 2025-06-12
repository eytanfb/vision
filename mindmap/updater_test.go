package mindmap

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestUpdaterProcessesEventsSerially(t *testing.T) {
	u := NewUpdater(t.TempDir())

	const events = 100
	var wg sync.WaitGroup
	wg.Add(events)

	// Fire events concurrently
	for i := 0; i < events; i++ {
		go func(idx int) {
			defer wg.Done()
			u.AppendTask("task-"+strconv.Itoa(idx), "Title")
		}(i)
	}

	wg.Wait()
	// Give the updater some time to drain
	time.Sleep(100 * time.Millisecond)

	u.Stop()

	if got := u.ProcessedCount(); got != events {
		t.Fatalf("expected %d processed events, got %d", events, got)
	}
}

func TestInitializeCreatesFile(t *testing.T) {
	dir := t.TempDir()
	u := NewUpdater(dir)

	today := time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC)
	u.InitializeDailyMindMap(today)
	u.InitializeDailyMindMap(today) // second call should be idempotent

	// Allow async processing
	time.Sleep(100 * time.Millisecond)
	u.Stop()

	path := filepath.Join(dir, "2025-06-12.txt")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected mind-map file to be created: %v", err)
	}

	data, _ := os.ReadFile(path)
	expected := "2025-06-12\n\tPersonal\n\tWork\n"
	if string(data) != expected {
		t.Fatalf("unexpected file content:\n%s", string(data))
	}
}

func TestAppendTask(t *testing.T) {
	dir := t.TempDir()
	u := NewUpdater(dir)

	date := time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC)
	u.InitializeDailyMindMap(date)
	u.enqueue(Event{Kind: EventAppendTask, Title: "Task A", Date: date})

	time.Sleep(100 * time.Millisecond)
	u.Stop()

	path := filepath.Join(dir, "2025-06-12.txt")
	data, _ := os.ReadFile(path)
	expectedContains := "\t\tTask A"
	if !strings.Contains(string(data), expectedContains) {
		t.Fatalf("expected file to contain %q, got:\n%s", expectedContains, string(data))
	}
}

func TestAppendSubtask(t *testing.T) {
	dir := t.TempDir()
	u := NewUpdater(dir)

	date := time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC)
	u.InitializeDailyMindMap(date)
	u.enqueue(Event{Kind: EventAppendTask, Title: "Parent Task", Date: date})
	u.enqueue(Event{Kind: EventAppendSubtask, ParentTaskID: "Parent Task", Title: "Child Subtask", Date: date})

	time.Sleep(200 * time.Millisecond)
	u.Stop()

	path := filepath.Join(dir, "2025-06-12.txt")
	dataLines, _ := os.ReadFile(path)
	s := string(dataLines)

	idxParent := strings.Index(s, "\t\tParent Task")
	idxChild := strings.Index(s, "\t\t\tChild Subtask")
	if idxParent == -1 || idxChild == -1 || idxChild < idxParent {
		t.Fatalf("expected child subtask below parent; content:\n%s", s)
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	dir := t.TempDir()
	u := NewUpdater(dir)
	date := time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC)
	u.InitializeDailyMindMap(date)
	u.enqueue(Event{Kind: EventAppendTask, Title: "My Task", Date: date})
	u.enqueue(Event{Kind: EventUpdateStatus, TaskID: "My Task", NewStatus: "scheduled", Date: date})

	time.Sleep(200 * time.Millisecond)
	u.Stop()

	path := filepath.Join(dir, "2025-06-12.txt")
	data, _ := os.ReadFile(path)
	s := string(data)
	matched, _ := regexp.MatchString(`My Task\s+⏳\s+2025-06-12`, s)
	if !matched {
		t.Fatalf("expected scheduled marker, got:\n%s", s)
	}
}
