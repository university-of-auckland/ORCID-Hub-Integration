package main

import (
	"strconv"
	"strings"
	"time"
)

// Task - ORCID Hub affiliation registration batch task
type Task struct {
	ID          int      `json:"id,omitempty"`
	CompletedAt string   `json:"completed-at,omitempty"`
	CreatedAt   string   `json:"created-at,omitempty"`
	ExpiresAt   string   `json:"expires-at,omitempty"`
	Filename    string   `json:"filename,omitempty"`
	Status      string   `json:"status,omitempty"`
	Type        string   `json:"task-type,omitempty"`
	Records     []Record `json:"records"`
}

// Record - ORCID Hub affiliation registration batch task recordr
type Record struct {
	ID                  int    `json:"id,omitempty"`
	AffiliationType     string `json:"affiliation-type"`
	City                string `json:"city,omitempty"`
	Country             string `json:"country,omitempty"`
	Department          string `json:"department,omitempty"`
	DisambiguatedID     string `json:"disambiguated-id,omitempty"`
	DisambiguatedSource string `json:"disambiguated-source,omitempty"`
	Email               string `json:"email,omitempty"`
	EndDate             string `json:"end-date,omitempty"`
	LocalID             string `json:"local-id,omitempty"`
	FirstName           string `json:"first-name,omitempty"`
	IsActive            bool   `json:"is-active,omitempty"`
	LastName            string `json:"last-name,omitempty"`
	Orcid               string `json:"orcid,omitempty"`
	Organisation        string `json:"organisation,omitempty"`
	ProcessedAt         string `json:"processed-at,omitempty"`
	PutCode             int    `json:"put-code,omitempty"`
	Role                string `json:"role,omitempty"`
	StartDate           string `json:"start-date,omitempty"`
	State               string `json:"state,omitempty"`
	Status              string `json:"status,omitempty"`
}

func (t *Task) activate() {
	var task Task
	log.Debugf("Activate the task %q (ID: %d)", t.Filename, t.ID)
	err := oh.patch("api/v1/tasks/"+strconv.Itoa(t.ID), map[string]string{"status": "ACTIVE"}, &task)
	if err != nil {
		log.Errorf("ERROR: Failed to activate task %d: %q", t.ID, err)
	}
}

func newTask() {

	taskFilename := taskFilenamePrefix + strconv.FormatInt(time.Now().Unix(), 36) + ".json"
	var task = Task{Filename: taskFilename, Type: "AFFILIATION", Records: []Record{}}
	err := oh.post("api/v1/affiliations?filename="+taskFilename, task, &task)
	if err != nil {
		logFatal("failed to create a new affiliation task", err)
	}
	taskID = task.ID
	taskCreatedAt, err = time.Parse("2006-01-02T15:04:05", task.CreatedAt)
	if err != nil {
		log.Errorf("failed to parse date %q: %s", task.CreatedAt, err)
	}
	log.Debugf("*** New affiliation task created (ID: %d, filename: %q)", task.ID, task.Filename)
}

// Either get the task ID or activate outstanding tasks and start a new one
func setupTask() (err error) {

	taskIDMutex.Lock()
	defer taskIDMutex.Unlock()

	now := time.Now()
	if taskID == 0 {
		var tasks []Task
		// Make sure the access token acquired
		log.Debug("=======================================================================================")
		oh.get("api/v1/tasks?type=AFFILIATION&status=INACTIVE", &tasks)
		for _, t := range tasks {
			log.Debugf("TASK: %+v", t)
			if t.Status == "ACTIVE" || t.Status == "RESET" || t.CompletedAt != "" || !strings.HasPrefix(t.Filename, taskFilenamePrefix) {
				continue
			}
			var createdAt time.Time
			createdAt, err = time.Parse("2006-01-02T15:04:05", t.CreatedAt)
			if err != nil {
				log.Error(err)
				return
			}
			if now.Sub(createdAt).Minutes() > taskRetentionMin && len(t.Records) > batchSize {
				t.activate()
				continue
			}
			taskID = t.ID
			taskCreatedAt = createdAt
			taskRecordCount = len(t.Records)
			return
		}
		newTask()

	} else if now.Sub(taskCreatedAt).Minutes() > taskRetentionMin && taskRecordCount > batchSize {
		log.Debug(now.Sub(taskCreatedAt).Minutes(), taskRetentionMin, taskRecordCount, batchSize)
		(&Task{ID: taskID}).activate()
		newTask()
	}
	return
}
