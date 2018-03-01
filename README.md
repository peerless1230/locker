# locker
基于`golang`，参考`Runc`和`Open Container Initiative`规范v1 [OpenContainer_SPEC.md](https://github.com/peerless1230/locker/blob/master/OpenContainer_SPEC.md)

目前实现功能：
- `Namespace`+`Cgroup`+`pivotRoot`实现的容器隔离运行环境及资源限制，已添加的Cgroup subsystem
    - cpu, 可通过
        - `cpu_share`控制容器进程cpu竞争优先级
        - `cpu-period`和`cpu-quota`联合控制`CFS (Completely Fair Scheduler)`调度机制下的cpu核数，属于对cpu资源绝对的限制
    - memory，限制容器进程内存消耗
    - cpuset，控制多核cpu下core的选择，如：0-3代表在core0到code3运行任务，0,2表示在core0，core2上
- 参照Docker `Overlay2`存储驱动，实现`OverlayFS`挂载容器`RootFS`
- 支持从`Overlay2`存储驱动下的`Docker镜像`启动容器
- 支持ps命令查看现有全面容器的信息，通过json序列/反序列化化,容器进程信息：PID，ContainerID，Name等
