/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/25 9:59
* @Mail: danbing.at@gmail.com
*/
package ziface

type IDataPack interface {
	//获取二进制包的头部长度  固定返回8
	GetHeadLen() uint32

 	//封包方法  ---- 将 Message  打包成 |datalen|dataID|data|\
 	Pack(msg IMessage) ([]byte, error)

	//拆包方法  ---  将|datalen|dataID|data|   拆解到 Message 结构体中
	UnPack([]byte) (IMessage, error)
}
