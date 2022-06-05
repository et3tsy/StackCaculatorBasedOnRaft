# StackCaculatorBasedOnRaft
 A stack caculator based on Raft. 西电2022年秋季分布式作业

Author: 谭升阳 19030100265

项目地址：https://github.com/et3tsy/StackCaculatorBasedOnRaft



## 设计需求

![WnqBC.png](https://s1.328888.xyz/2022/06/05/WnqBC.png)



## 设计结构

| **第四层** | **TCP处理客户请求**  |
| :--------: | :------------------: |
| **第三层** | **栈计算器具体实现** |
| **第二层** |   **Raft基本框架**   |
| **第一层** | **基于TCP的RPC框架** |

### 第一层

Go官方提供了一个[RPC库](https://golang.org/pkg/net/rpc/): `net/rpc`。

对于所有暴露的方法，均采用如下格式：

```
func (t *T) MethodName(argType T1, replyType *T2) error
```

我们主要将`AppendEntriesRPC` 和 `RequestVoteRPC` 进行注册。但是，`net\rpc` 是利用反射机制将所有大写方法遍历完成注册的，所以这里我们如果在第二层中暴露了其他大写字母方法时，干脆就让 `net\rpc` 报错，不处理即可。

### 第二层

基本完成了 Raft 协议的精髓部分，即Leader选举，日志复制。但是，没有实现持久化，快照，动态增减集群（热插拔）。

进行第二层开发时，为确保正确性，我的代码是基于麻省理工 MIT6.824 lab2 移植过来的。该实验有比较成熟的测试框架，方便模拟各类真实网络环境，设计了很多corner cases，保证代码的正确性。根据实验要求，不能进行开源。

实验地址：

```
https://pdos.csail.mit.edu/6.824/labs/lab-raft.html
```

我独立实现的lab2A，lab2B，在大多数情况下，代码没有产生异常，但是在压力测试下，2B在一些测试中，产生了数据竞争，进而导致的访问越界问题。由于时间关系，没有进一步的调试和修正，等有空再回头做。

我在这一阶段前前后后花了7天，有些内容是论文中 firgure 2 中没有提到的（但是在后面有提及）。也是我在实验中碰到的很多坑。

（1）什么时候将 leader Append 到的新数据复制给其他 peers？

​	方案一：触发机制，一有 client 添加数据，就开始复制。

​	方案二：定时检查。

对于方案一，有个致命问题，如果说有从机宕机，当他恢复以后，若长时间没有用户写入，那么容易产生问题。

因此，采用定时的方式。

（2）复制包与心跳包的关系？

每轮心跳时间一到，检测 Append 的新数据，并完成复制，这次视为一次心跳包。

（3）为什么在进行回退 nextIndex、向 peers 进行 matching prev entry 的时候要捎带 nextIndex的数据块？

### 第三层

设计栈计算器，并且实现基于 Raft 设计用户交互。Raft 保证读写都是向 leader 发出的。

设计读请求，需要保证当前读到的数据是最新的（not return stale data），按照论文提到的方案需要做到两点保证。一是，leader当选以后，主动 Append 一个空的 Entry，并且等到这个 Entry 被提交。这是为了当前的 leader 得到最新的 commitIndex。 当前 leader 可能虽然拥有最新的 Entry， 但是可能不知道它已经被提交了，旧 leader 提交的消息没传递过来（在选举的时候，他用当前最新的 Entry 向大多数获得的选票，但是，旧 leader 更新 commitIndex 的消息并没有通过 AppendEntries 传递过来。

第二是，为了保证当前 leader 依然是 leader，它需要主动发一次心跳包，大多数的从机的 term 没有变化，就说明肯定没有新 leader 上任。

### 第四层

利用TCP协程池处理用户请求。
