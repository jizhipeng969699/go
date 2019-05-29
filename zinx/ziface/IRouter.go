/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/23 15:24
* @Mail: danbing.at@gmail.com
*/
package ziface

/*
   抽象的路由模块
 */
type IRouter interface {
	//处理业务之前的方法
	PreHandle(request IRequest)
	//真正处理业务的方法
	Handle(request IRequest)
	//处理业务之后的方法
	PostHandle(request IRequest)
}