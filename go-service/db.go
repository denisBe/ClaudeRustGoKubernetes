package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JobInfo struct {
	ID     string `json:"id"`
	Filter string `json:"filter"`
	Status string `json:"status"`
}

type JobQueueMessage struct {
	ID        string `json:"id"`
	Filter    string `json:"filter"`
	ImageData []byte `json:"image_data"`
}

type DbContext struct {
	redisClient *redis.Client
}

func (db *DbContext) CreateJob(image []byte, filter string) (JobInfo, error) {
	job := JobInfo{
		ID:     uuid.New().String(),
		Filter: filter,
		Status: "pending",
	}

	ctx := context.Background()

	// Store job status in Redis (for querying via handleGetJob)
	jobJson, err := json.Marshal(job)
	if err != nil {
		fmt.Println("Failed to marshal job info:", err)
		return job, err
	}
	err = db.redisClient.Set(ctx, "job:"+job.ID, string(jobJson), 0).Err()
	if err != nil {
		fmt.Println("Failed to store job info in Redis:", err)
		return job, err
	}

	// Push job to queue for Rust worker to consume
	queueMsg := JobQueueMessage{
		ID:        job.ID,
		Filter:    filter,
		ImageData: image,
	}
	queueJson, err := json.Marshal(queueMsg)
	if err != nil {
		fmt.Println("Failed to marshal queue message:", err)
		return job, err
	}
	err = db.redisClient.RPush(ctx, "jobs:queue", string(queueJson)).Err()
	if err != nil {
		fmt.Println("Failed to push job to queue:", err)
		return job, err
	}

	return job, nil
}

func (db *DbContext) GetJobInfo(jobID string) (JobInfo, error) {
	ctx := context.Background()
	val, err := db.redisClient.Get(ctx, "job:"+jobID).Result()
	if err != nil {
		return JobInfo{}, err
	}

	var info JobInfo
	if err = json.Unmarshal([]byte(val), &info); err != nil {
		return JobInfo{}, fmt.Errorf("stored job info is not valid JSON: %w", err)
	}
	return info, nil
}
