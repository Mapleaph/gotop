package widgets

import (
	psProc "github.com/shirou/gopsutil/process"
)

func (self *Proc) update() {
	psProcesses, err := psProc.Processes()
	if err != nil {
		return err
	}
	processes := make([]Process, len(psProcesses))
	for i, psProcess := range psProcesses {
		pid := psProcess.Pid
		command, err := psProcess.Name()
		if err != nil {
			return err
		}
		cpu, err := psProcess.CPUPercent()
		if err != nil {
			return err
		}
		mem, err := psProcess.MemoryPercent()
		if err != nil {
			return err
		}

		processes[i] = Process{
			int(pid),
			command,
			cpu / self.cpuCount,
			float64(mem),
			// getting command args using gopsutil's Cmdline and CmdlineSlice wasn't
			// working the last time I tried it, so we're just reusing 'command'
			command,
		}
	}

	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
}
