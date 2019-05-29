/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/25 17:22
* @Mail: danbing.at@gmail.com
*/
package net

import (
	"fmt"
	"zinx/config"
	"zinx/ziface"
)

type MsgHandler struct {
	//存放路由集合的map
	Apis map[uint32] ziface.IRouter //就是开发者全部的业务，消息ID和业务的对应关系

	//负责Worker取任务的消息队列  一个worker对应一个任务队列
	TaskQueue  []chan ziface.IRequest

	//worker工作池的worker数量
	WorkerPoolSize uint32
}

//初始化方法
func NewMsgHandler() ziface.IMsgHandler {
	//给map开辟头空间
	return &MsgHandler{
		Apis:make(map[uint32]ziface.IRouter),
		WorkerPoolSize:config.GlobalObject.WorkerPoolSize,
		TaskQueue:make([]chan ziface.IRequest, config.GlobalObject.WorkerPoolSize),//切片的初始化
	}
}

//添加路由到map集合中
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断新添加的msgID key是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//msgId已经注册
		fmt.Println("repeat Api msgID = ", msgID)
		return
	}
	//2 添加msgID 和 router的对应关系
	mh.Apis[msgID] = router
	fmt.Println("Apd api MsgID = ", msgID, " succ!")
}

//调度路由， 根据MsgID
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1 从Request 取到MsgiD
	router, ok := mh.Apis[request.GetMsg().GetMsgId()]
	if !ok {
		fmt.Println("api MsgID = ", request.GetMsg().GetMsgId(), " Not Found! Need Add！")
		return
	}
	//2 根据msgID  找到对应的router 进行调用
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)

}

//一个worker真正处理业务的 goroutine函数
func (mh *MsgHandler) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println(" worker ID = ", workerID , " is starting... ")

	//不断的从对应的管道 等待数据
	for {
		select {
			case req := <-taskQueue:
				mh.DoMsgHandler(req)
		}
	}
}

//启动Worker工作池 (在整个server服务中 只启动一次)
func (mh *MsgHandler) StartWorkerPool() {
	fmt.Println("WorkPool is  started..")

	//根据WorkerPoolSize 创建worker goroutine
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//开启一个workergoroutine

		//1 给当前Worker所绑定消息channel对象 开辟空间  第0个worker 就用第0个Channel
		//给channel 进行开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, config.GlobalObject.MaxWorkerTaskLen)

		//2 启动一个Worker，阻塞等待消息从对应的管道中进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}

}


//将消息添加到Worker工作池中 （将消息发送给对应的消息队列）
//应该是Reader来调用的
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 将消息 平均分配给worker 确定当前的request到底要给哪个worker来处理
	//1个客户端绑定一个worker来处理
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	//2 直接将 request 发送给对应的worker的taskqueue
	mh.TaskQueue[workerID] <- request
}