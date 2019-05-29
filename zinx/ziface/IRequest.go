/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/23 11:56
* @Mail: danbing.at@gmail.com
*/
package ziface

/*
 抽象 IRequest 一次性请求的数据封装
 */

 type IRequest interface {
 	//得到当前的请求的链接
 	GetConnection() IConnection

 	//得到请求的消息
 	GetMsg() IMessage
 }
