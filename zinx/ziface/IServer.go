/**
 Server模块的抽象层
*/
package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()

	//提供一个得到链接管理模块的方法
	GetConnMgr() IConnManager

	//添加路由方法  暴露给开发者的
	AddRouter(msgID uint32, router IRouter)

	AddOnConnStart(hookFunc func(conn IConnection))
	AddOnConnStop(hookFunc func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}


