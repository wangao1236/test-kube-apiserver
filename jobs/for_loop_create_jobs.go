package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"k8s.io/klog"
)

func (j *JobInfo) process(ctx context.Context, idx int64, wg *sync.WaitGroup) {
	job := buildJobConfig(idx)
	defer func() {
		klog.Infof("Stop executing task of %+v", idx)
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(rand.Intn(MaxWaitMillisecond)) * time.Millisecond):
			j.CreateJob(job)
			j.DeleteJob(job.Name)
		}
	}
}

// ParallelizeUntil executes the tasks in parallel and decides whether to stop based on the context
func (j *JobInfo) ParallelizeUntil() {
	j.ClearJobs()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	var i int64
	for i = 0; i < j.concurrentNum; i++ {
		wg.Add(1)
		go j.process(ctx, i, wg)
	}

	defer func() {
		cancel()
		wg.Wait()
		klog.Info("All task done")
		time.Sleep(3 * time.Second)
	}()

	<-j.exit
	time.Sleep(time.Second)
}

func main() {
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())
	if jobInfo, err := NewJobInfo(); err == nil {
		jobInfo.ParallelizeUntil()
	} else {
		klog.Errorf("PullInfo create error: %v", err)
	}
}
