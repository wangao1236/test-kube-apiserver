package main

import (
	"context"
	"flag"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	// MaxWaitMillisecond represents the maximum waiting time for each round of execution
	MaxWaitMillisecond = 10
	// MaxConcurrentNum represents the maximum number of concurrent tasks
	MaxConcurrentNum = 20
	// JobBatchSize represents the number of jobs sent per goroutine
	JobBatchSize = 50
	// Client represents the number of ClientSets
	ClientSetNum = 1
	// JobParallelism represents the maximum number of a job in parallel
	JobParallelism = 2
	// JobCompletions represents the minimum completion of a job
	JobCompletions = 4
	// ConfigQps indicates the maximum QPS to the master from this client
	ConfigQps = 4000
	// ConfigBurst represents the maximum burst for throttle
	ConfigBurst = 8000
)

// JobInfo records test jobs' information
type JobInfo struct {
	exit          chan os.Signal
	clientSet     *kubernetes.Clientset
	jobsClient    v1.JobInterface
	concurrentNum int64
}

// PrintLocalDial outputs dialing information, including TCP connection status
func PrintLocalDial(ctx context.Context, network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	klog.Info("connect done, use ", conn.LocalAddr().String())

	return conn, err
}

// NewConfig reads the kubeconfig file and creates the Config object
func NewConfig() (*restclient.Config, error) {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		klog.Warningf("NewConfig: err: %+v", err)
		return nil, err
	} else if config != nil {
		config.QPS = float32(ConfigQps)
		config.Burst = ConfigBurst
		config.Dial = PrintLocalDial
		config.Timeout = 300 * time.Second
	}
	klog.Infof("config: %+v", config)
	return config, nil
}

// NewJobInfo creates a new JobInfo
func NewJobInfo() (*JobInfo, error) {
	jobInfo := &JobInfo{}
	config, err := NewConfig()
	if err != nil {
		klog.Warningf("NewJobInfo: err: %+v", err)
		return nil, err
	}

	jobInfo.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		klog.Warningf("NewJobInfo: err: %+v", err)
		return nil, err
	}

	jobInfo.jobsClient = jobInfo.clientSet.BatchV1().Jobs(apiv1.NamespaceDefault)
	jobInfo.exit = make(chan os.Signal, 1)

	signal.Notify(jobInfo.exit, os.Interrupt, os.Kill, syscall.SIGTERM)
	jobInfo.FlagParse()
	return jobInfo, nil
}

func buildJobConfig(idx int64) *batchv1.Job {
	parallelism := int32(JobParallelism)
	completions := int32(JobCompletions)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-job-" + strconv.FormatInt(idx, 10) + "-" +
				strconv.FormatInt(time.Now().UnixNano(), 10),
			Namespace: apiv1.NamespaceDefault,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Completions: &completions,
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "pi",
							Image:   "resouer/ubuntu-bc",
							Command: []string{"sh", "-c", "echo 'scale=10; 4*a(1)' | bc -l "},
						},
					},
					RestartPolicy: apiv1.RestartPolicyOnFailure,
				},
			},
		},
	}
	return job
}

// FlagParse parses command line parameters for JobInfo
func (j *JobInfo) FlagParse() {
	var concurrent int64
	flag.Int64Var(&concurrent, "concurrent", MaxConcurrentNum, "If empty, use this number")
	flag.Parse()
	//
	runtime.GOMAXPROCS(10)

	j.concurrentNum = concurrent
}

// CreateJob executes task of job's creation
func (j *JobInfo) CreateJob(job *batchv1.Job) {
	klog.Infof("Creating job: %v...", job.Name)
	result, err := j.jobsClient.Create(job)
	if err != nil {
		klog.Errorf("Create job: %v err: %v", job.Name, err)
		return
	}
	klog.Infof("Created job: %v", result.Name)
}

// DeleteJob executes task of job's deletion
func (j *JobInfo) DeleteJob(jobName string) {
	klog.Infof("Deleting job: %v...", jobName)
	deletePolicy := metav1.DeletePropagationForeground
	if err := j.jobsClient.Delete(jobName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		klog.Errorf("Delete job: %v err: %v", jobName, err)
	}
	klog.Infof("Deleted job %v", jobName)
}

// ListJobs executes task of jobs' enumeration, and returns all jobs' name
func (j *JobInfo) ListJobs() []string {
	klog.Infof("Listing jobs in namespace %v", apiv1.NamespaceDefault)
	list, err := j.jobsClient.List(metav1.ListOptions{})
	if err != nil {
		klog.Errorf("List jobs err: %v", err)
	}
	klog.Infof("Items: %+v", list.Items)
	jobNames := make([]string, 0)
	for _, d := range list.Items {
		jobNames = append(jobNames, d.Name)
	}
	return jobNames
}

// ClearJobs executes task of jobs' cleanup
func (j *JobInfo) ClearJobs() {
	jobNames := j.ListJobs()
	for _, name := range jobNames {
		j.DeleteJob(name)
	}
}

func (j *JobInfo) processBatch(idx int64, wg *sync.WaitGroup) {
	defer func() {
		klog.Infof("Stop executing task of %+v", idx)
		wg.Done()
	}()

	for i := 0; i < JobBatchSize; i++ {
		job := buildJobConfig(idx)
		j.CreateJob(job)
	}
	klog.Infof("finish task: %v", idx)
}

// ParallelizeBatch performs batch tasks in parallel
func (j *JobInfo) ParallelizeBatch() {
	//j.ClearJobs()
	wg := &sync.WaitGroup{}

	var i int64
	for i = 0; i < j.concurrentNum; i++ {
		wg.Add(1)
		go j.processBatch(i, wg)
	}

	defer func() {
		klog.Info("All task done")
		time.Sleep(3 * time.Second)
	}()

	wg.Wait()
}

func main() {
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())
	if jobInfo, err := NewJobInfo(); err == nil {
		jobInfo.ParallelizeBatch()
	} else {
		klog.Errorf("PullInfo create error: %v", err)
	}
}
