## Target workload
- Pod

## How-to
- create a new mpod(monitoring pod), scheduing to the same node of the target pod
- get pids of the target pods, maybe contain multiple pods
- run ebpf snooping those pids and report

## CRD
- namespaced
- pod object name(namespaced unique)
- pod uid(cluster unique)

## Controller
- get the tpod by pod object name
- create a mpod scheduling to the node running opod
