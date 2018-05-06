package engine

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/blkiodev"
	"github.com/docker/docker/client"
	"time"
)

type Client struct {
	api       *client.Client
	cgroup    string
	container *types.ContainerJSON

	cancelEvents context.CancelFunc

	closed bool
}

type Component struct {
	Image           string
	Entrypoint      interface{}
	Command         interface{}
	WorkingDir      string      `yaml:"working_dir"`
	EnvFile         interface{} `yaml:"env_file"`
	Environment     interface{}
	Labels          interface{}
	Privileged      bool
	ReadOnly        bool `yaml:"read_only"`
	StdinOpen       bool `yaml:"stdin_open"`
	Tty             bool
	StopSignal      string        `yaml:"stop_signal"`
	StopGracePeriod time.Duration `yaml:"stop_grace_period"`
	User            string

	Healthcheck *Healthcheck

	OomScoreAdj    *int  `yaml:"oom_score_adj"`
	OomKillDisable *bool `yaml:"oom_kill_disable"`

	CapAdd  []string `yaml:"cap_add"`
	CapDrop []string `yaml:"cap_drop"`

	Devices           []string
	DeviceCgroupRules []string `yaml:"device_cgroup_rules"`

	MemoryLimit       string  `yaml:"mem_limit"`
	MemoryReservation *int64  `yaml:"mem_reservation"`
	MemorySwapLimit   string  `yaml:"memswap_limit"`
	MemorySwappiness  *int64  `yaml:"mem_swappiness"`
	ShmSize           *string `yaml:"shm_size"`

	CPUShares          int64   `yaml:"cpu_shares"`
	CPUs               float64 `yaml:"cpus"`
	CPUPeriod          int64   `yaml:"cpu_period"`
	CPUQuota           int64   `yaml:"cpu_quota"`
	CPURealtimePeriod  int64   `yaml:"cpu_rt_period"`
	CPURealtimeRuntime int64   `yaml:"cpu_rt_runtime"`
	CpusetCpus         string  `yaml:"cpuset"`
	CPUCount           int64   `yaml:"cpu_count"`
	CPUPercent         int64   `yaml:"cpu_percent"`

	BlkioConfig *BlkioConfig `yaml:"blkio_config"`

	GroupAdd   []string `yaml:"group_add"`
	UsernsMode string   `yaml:"userns_mode"`

	// the parent client to the engine
	client *Client `yaml:"-"`

	// the name and container ID set in runtime
	Name      string               `yaml:"-"`
	container *types.ContainerJSON `yaml:"-"`
}

type Healthcheck struct {
	Test        interface{}
	Interval    time.Duration
	Timeout     time.Duration
	StartPeriod time.Duration `yaml:"start_period"`
	Retries     int
	Disable     bool
}

type BlkioConfig struct {
	Weight          uint16
	WeightDevice    []*blkiodev.WeightDevice   `yaml:"weight_device"`
	DeviceReadBps   []*blkiodev.ThrottleDevice `yaml:"device_read_bps"`
	DeviceWriteBps  []*blkiodev.ThrottleDevice `yaml:"device_read_iops"`
	DeviceReadIOps  []*blkiodev.ThrottleDevice `yaml:"device_write_bps"`
	DeviceWriteIOps []*blkiodev.ThrottleDevice `yaml:"device_write_iops"`
}

type ComponentExited struct {
	Component *Component

	StatusCode int64
	Error      error
}
