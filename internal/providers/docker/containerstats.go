package docker

type ContainerStats struct {
	BlkioStats struct {
		IoServiceBytesRecursive []struct {
			Major int64  `json:"major"`
			Minor int64  `json:"minor"`
			Op    string `json:"op"`
			Value int64  `json:"value"`
		} `json:"io_service_bytes_recursive"`
	} `json:"blkio_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        uint64 `json:"total_usage"`
			UsageInKernelmode uint64 `json:"usage_in_kernelmode"`
			UsageInUsermode   uint64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		OnlineCpus     int64 `json:"online_cpus"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttled_periods"`
			ThrottledTime    int64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	ID          string `json:"id"`
	MemoryStats struct {
		Limit int64 `json:"limit"`
		Stats struct {
			ActiveAnon            int64 `json:"active_anon"`
			ActiveFile            int64 `json:"active_file"`
			Anon                  int64 `json:"anon"`
			AnonThp               int64 `json:"anon_thp"`
			File                  int64 `json:"file"`
			FileDirty             int64 `json:"file_dirty"`
			FileMapped            int64 `json:"file_mapped"`
			FileWriteback         int64 `json:"file_writeback"`
			InactiveAnon          int64 `json:"inactive_anon"`
			InactiveFile          int64 `json:"inactive_file"`
			KernelStack           int64 `json:"kernel_stack"`
			Pgactivate            int64 `json:"pgactivate"`
			Pgdeactivate          int64 `json:"pgdeactivate"`
			Pgfault               int64 `json:"pgfault"`
			Pglazyfree            int64 `json:"pglazyfree"`
			Pglazyfreed           int64 `json:"pglazyfreed"`
			Pgmajfault            int64 `json:"pgmajfault"`
			Pgrefill              int64 `json:"pgrefill"`
			Pgscan                int64 `json:"pgscan"`
			Pgsteal               int64 `json:"pgsteal"`
			Shmem                 int64 `json:"shmem"`
			Slab                  int64 `json:"slab"`
			SlabReclaimable       int64 `json:"slab_reclaimable"`
			SlabUnreclaimable     int64 `json:"slab_unreclaimable"`
			Sock                  int64 `json:"sock"`
			ThpCollapseAlloc      int64 `json:"thp_collapse_alloc"`
			ThpFaultAlloc         int64 `json:"thp_fault_alloc"`
			Unevictable           int64 `json:"unevictable"`
			WorkingsetActivate    int64 `json:"workingset_activate"`
			WorkingsetNodereclaim int64 `json:"workingset_nodereclaim"`
			WorkingsetRefault     int64 `json:"workingset_refault"`
		} `json:"stats"`
		Usage uint64 `json:"usage"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	Networks struct {
		Eth0 struct {
			RxBytes   int64 `json:"rx_bytes"`
			RxDropped int64 `json:"rx_dropped"`
			RxErrors  int64 `json:"rx_errors"`
			RxPackets int64 `json:"rx_packets"`
			TxBytes   int64 `json:"tx_bytes"`
			TxDropped int64 `json:"tx_dropped"`
			TxErrors  int64 `json:"tx_errors"`
			TxPackets int64 `json:"tx_packets"`
		} `json:"eth0"`
	} `json:"networks"`
	NumProcs  int64 `json:"num_procs"`
	PIDsStats struct {
		Current int64 `json:"current"`
		Limit   int64 `json:"limit"`
	} `json:"pids_stats"`
	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			UsageInKernelmode int64 `json:"usage_in_kernelmode"`
			UsageInUsermode   int64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		OnlineCPUs     int64 `json:"online_cpus"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttled_periods"`
			ThrottledTime    int64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	Preread      string   `json:"preread"`
	Read         string   `json:"read"`
	StorageStats struct{} `json:"storage_stats"`
}
