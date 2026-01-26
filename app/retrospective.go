package app

import (
	"math"
	"sort"
	"time"
)

// TaskRetrospective holds computed retrospective data for a completed task
type TaskRetrospective struct {
	Task          Task
	CreatedDate   time.Time  // From file frontmatter or earliest date
	ScheduledDate *time.Time // Optional - when scheduled
	StartDate     *time.Time // Optional - when started
	CompletedDate time.Time  // Required for retrospective

	// Computed metrics
	TotalDays       int
	WaitDays        int     // Created to Scheduled (or Started if no scheduled)
	LeadDays        int     // Scheduled to Started
	ActiveDays      int     // Started to Completed
	EfficiencyRatio float64 // ActiveDays / TotalDays
}

// TimelineEvent represents a single event in the task timeline
type TimelineEvent struct {
	EventType string    // "created", "scheduled", "started", "completed"
	Date      time.Time
	DayNumber int    // Days since creation
	Icon      string
	Label     string
}

// RetrospectiveMetric represents a computed metric for display
type RetrospectiveMetric struct {
	Label      string
	Value      int
	Unit       string
	Percentage float64
	BarLength  int // For visual bar rendering (0-20)
}

const dateFormat = "2006-01-02"

// NewTaskRetrospective creates a new TaskRetrospective from a task and its creation date.
// It parses the task's date strings and computes all metrics.
func NewTaskRetrospective(task Task, createdDate time.Time) *TaskRetrospective {
	r := &TaskRetrospective{
		Task:        task,
		CreatedDate: createdDate,
	}

	// Parse completed date - required for retrospective
	if task.CompletedDate == "" {
		return nil
	}
	completedDate, err := time.Parse(dateFormat, task.CompletedDate)
	if err != nil {
		return nil
	}
	r.CompletedDate = completedDate

	// Parse optional scheduled date
	if task.ScheduledDate != "" {
		scheduledDate, err := time.Parse(dateFormat, task.ScheduledDate)
		if err == nil {
			r.ScheduledDate = &scheduledDate
		}
	}

	// Parse optional start date
	if task.StartDate != "" {
		startDate, err := time.Parse(dateFormat, task.StartDate)
		if err == nil {
			r.StartDate = &startDate
		}
	}

	r.computeMetrics()
	return r
}

// computeMetrics calculates all the retrospective metrics
func (r *TaskRetrospective) computeMetrics() {
	// TotalDays: Created to Completed
	r.TotalDays = daysBetween(r.CreatedDate, r.CompletedDate)

	// WaitDays: Created to Scheduled (or Started if no scheduled, or Completed if neither)
	if r.ScheduledDate != nil {
		r.WaitDays = daysBetween(r.CreatedDate, *r.ScheduledDate)
	} else if r.StartDate != nil {
		r.WaitDays = daysBetween(r.CreatedDate, *r.StartDate)
	} else {
		r.WaitDays = 0
	}

	// LeadDays: Scheduled to Started (0 if either is missing)
	if r.ScheduledDate != nil && r.StartDate != nil {
		r.LeadDays = daysBetween(*r.ScheduledDate, *r.StartDate)
	} else {
		r.LeadDays = 0
	}

	// ActiveDays: Started to Completed (0 if no start date)
	if r.StartDate != nil {
		r.ActiveDays = daysBetween(*r.StartDate, r.CompletedDate)
	} else {
		r.ActiveDays = 0
	}

	// EfficiencyRatio: ActiveDays / TotalDays
	r.EfficiencyRatio = calculateEfficiencyRatio(r.ActiveDays, r.TotalDays)
}

// GetTimeline returns ordered timeline events for the task
func (r *TaskRetrospective) GetTimeline() []TimelineEvent {
	events := []TimelineEvent{
		{
			EventType: "created",
			Date:      r.CreatedDate,
			DayNumber: 0,
			Icon:      "+",
			Label:     "Created",
		},
	}

	if r.ScheduledDate != nil {
		events = append(events, TimelineEvent{
			EventType: "scheduled",
			Date:      *r.ScheduledDate,
			DayNumber: daysBetween(r.CreatedDate, *r.ScheduledDate),
			Icon:      ScheduledIcon,
			Label:     "Scheduled",
		})
	}

	if r.StartDate != nil {
		events = append(events, TimelineEvent{
			EventType: "started",
			Date:      *r.StartDate,
			DayNumber: daysBetween(r.CreatedDate, *r.StartDate),
			Icon:      StartedIcon,
			Label:     "Started",
		})
	}

	events = append(events, TimelineEvent{
		EventType: "completed",
		Date:      r.CompletedDate,
		DayNumber: daysBetween(r.CreatedDate, r.CompletedDate),
		Icon:      CompletedIcon,
		Label:     "Completed",
	})

	// Sort by date
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})

	return events
}

// GetMetrics returns computed metrics formatted for display
func (r *TaskRetrospective) GetMetrics() []RetrospectiveMetric {
	metrics := []RetrospectiveMetric{}

	// Total days metric
	totalMetric := RetrospectiveMetric{
		Label:      "Total Time",
		Value:      r.TotalDays,
		Unit:       pluralize(r.TotalDays, "day", "days"),
		Percentage: 100.0,
		BarLength:  20, // Full bar for total
	}
	metrics = append(metrics, totalMetric)

	// Only show wait days if there's meaningful wait time
	if r.WaitDays > 0 || r.TotalDays > 0 {
		waitPct := calculatePercentage(r.WaitDays, r.TotalDays)
		waitMetric := RetrospectiveMetric{
			Label:      "Wait Time",
			Value:      r.WaitDays,
			Unit:       pluralize(r.WaitDays, "day", "days"),
			Percentage: waitPct,
			BarLength:  percentageToBarLength(waitPct),
		}
		metrics = append(metrics, waitMetric)
	}

	// Lead days (scheduled to started)
	if r.ScheduledDate != nil && r.StartDate != nil {
		leadPct := calculatePercentage(r.LeadDays, r.TotalDays)
		leadMetric := RetrospectiveMetric{
			Label:      "Lead Time",
			Value:      r.LeadDays,
			Unit:       pluralize(r.LeadDays, "day", "days"),
			Percentage: leadPct,
			BarLength:  percentageToBarLength(leadPct),
		}
		metrics = append(metrics, leadMetric)
	}

	// Active days (started to completed)
	if r.StartDate != nil {
		activePct := calculatePercentage(r.ActiveDays, r.TotalDays)
		activeMetric := RetrospectiveMetric{
			Label:      "Active Time",
			Value:      r.ActiveDays,
			Unit:       pluralize(r.ActiveDays, "day", "days"),
			Percentage: activePct,
			BarLength:  percentageToBarLength(activePct),
		}
		metrics = append(metrics, activeMetric)
	}

	// Efficiency ratio as a percentage metric
	efficiencyMetric := RetrospectiveMetric{
		Label:      "Efficiency",
		Value:      int(r.EfficiencyRatio * 100),
		Unit:       "%",
		Percentage: r.EfficiencyRatio * 100,
		BarLength:  percentageToBarLength(r.EfficiencyRatio * 100),
	}
	metrics = append(metrics, efficiencyMetric)

	return metrics
}

// calculateEfficiencyRatio computes the ratio of active days to total days.
// Returns 1.0 if total is 0 (same-day completion), otherwise active/total.
func calculateEfficiencyRatio(active, total int) float64 {
	if total == 0 {
		return 1.0 // Same-day completion is 100% efficient
	}
	if active == 0 {
		return 0.0 // No active time means 0% efficiency
	}
	return float64(active) / float64(total)
}

// daysBetween calculates the number of days between two dates.
// Returns 0 for same-day, 1 for consecutive days, etc.
func daysBetween(start, end time.Time) int {
	// Normalize to start of day to avoid time-of-day issues
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	days := int(endDay.Sub(startDay).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// calculatePercentage computes part as a percentage of total
func calculatePercentage(part, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return (float64(part) / float64(total)) * 100.0
}

// percentageToBarLength converts a percentage (0-100) to bar length (0-20)
func percentageToBarLength(pct float64) int {
	if pct <= 0 {
		return 0
	}
	if pct >= 100 {
		return 20
	}
	return int(math.Round(pct / 5.0))
}

// pluralize returns singular or plural form based on count
func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
