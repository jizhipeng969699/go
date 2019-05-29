/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/26 15:33
* @Mail: danbing.at@gmail.com
*/
package ziface

type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(connID uint32)
	//根据链接ID得到链接
	Get(connID uint32) (IConnection, error)
	//得到目前服务器的链接总个数
	Len() int
	//清空全部链接的方法
	ClearConn()
}
