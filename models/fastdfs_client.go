package models

import (
	"github.com/weilaihui/fdfs_client"
	"github.com/astaxie/beego"
)

//import (
//	"github.com/weilaihui/fdfs_client"
//	"github.com/astaxie/beego"
//)

////通过文件名的方式进行上传
//func TestUploadByFilename(t *testing.T) {
//
//	//通过配置文件创一个客户端fdfs的句柄
//	fdfsClient, err := NewFdfsClient("client.conf")
//	//判断
//	if err != nil {
//		t.Errorf("New FdfsClient error %s", err.Error())
//		return
//	}
//	//通过句柄调用，使用文件名上传文件的方法（客户端配置文件）
//	uploadResponse, err = fdfsClient.UploadByFilename("client.conf")
//	//判断
//	if err != nil {
//		t.Errorf("UploadByfilename error %s", err.Error())
//	}
//
//	//打印存放地址
//	//store_path0=/home/itcast/go/src/29q/day9/ihome/fastdfs/storage_data
//	t.Log(uploadResponse.GroupName)
//	//打印fileid
//	//group1/M00/00/00/wKgLg1txSyiAE7I_AAaInSze-cQ989.jpg
//	t.Log(uploadResponse.RemoteFileId)
//	//通过句柄对文件进行删除
//	fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
//}

func UploadByFilename( filename string)(GroupName,RemoteFileId string ,err error ) {
	//通过配置文件创建fdfs操作句柄
	fdfsClient, thiserr :=fdfs_client.NewFdfsClient("./conf/client.conf")
	if thiserr  !=nil{
		beego.Info("UploadByFilename( ) fdfs_client.NewFdfsClient  err",err)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}

	//unc (this *FdfsClient) UploadByFilename(filename string) (*UploadFileResponse, error)
	//通过句柄上传文件（被上传的文件）

	uploadResponse, thiserr := fdfsClient.UploadByFilename(filename)
	if thiserr !=nil{
		beego.Info("UploadByFilename( ) fdfsClient.UploadByFilename(filename)  err",err)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}

	beego.Info(uploadResponse.GroupName)
	beego.Info(uploadResponse.RemoteFileId)



	return  uploadResponse.GroupName , uploadResponse.RemoteFileId ,nil

}



/*//通过二进制文件的方式进行上传
func TestUploadByBuffer(t *testing.T) {
	//通过配置文件创建一个fastdfs的句柄
	fdfsClient, err := NewFdfsClient("client.conf")
	//判断是否有误
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}


	//打开文件
	file, err := os.Open("testfile") // For read access.
	if err != nil {
		t.Fatal(err)
	}
	//创建一个变量用来存放文件的大小
	var fileSize int64 = 0
	//通过stat获取 文件信息 赋值给我们的fileinfo
	if fileInfo, err := file.Stat(); err == nil {
		//获取文件的大小
		fileSize = fileInfo.Size()
	}

	//根据文件的大小创建一个byte切片
	fileBuffer := make([]byte, fileSize)
	//读取文件存放到创建的byte的切片当中
	_, err = file.Read(fileBuffer)
	if err != nil {
		t.Fatal(err)
	}
	//通过句柄上传二进制文件（二进制的切片 ，文件的后缀名）
	uploadResponse, err = fdfsClient.UploadByBuffer(fileBuffer, "txt")
	if err != nil {
		t.Errorf("TestUploadByBuffer error %s", err.Error())
	}
	//打印GroupName
	t.Log(uploadResponse.GroupName)
	//打印RemoteFileId
	t.Log(uploadResponse.RemoteFileId)
	//删除文件
	fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
}*/
//功能函数 操作fdfs上传二进制文件
func UploadByBuffer(filebuffer []byte, fileExtName string)(GroupName,RemoteFileId string ,err error ){

	//通过配置文件创建fdfs操作句柄
	fdfsClient, thiserr :=fdfs_client.NewFdfsClient("./conf/client.conf")
	if thiserr  !=nil{
		beego.Info("UploadByBuffer( ) fdfs_client.NewFdfsClient  err",err)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}
	//func (this *FdfsClient) UploadByBuffer(filebuffer []byte, fileExtName string) (*UploadFileResponse, error) {
	//通过句柄上传二进制的文件
	uploadResponse, thiserr :=fdfsClient.UploadByBuffer(filebuffer,fileExtName)
	if thiserr  !=nil{
		beego.Info("UploadByBuffer( ) fdfs_client.UploadByBuffer  err",err)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}
	beego.Info(uploadResponse.GroupName)
	beego.Info(uploadResponse.RemoteFileId)

	return uploadResponse.GroupName,uploadResponse.RemoteFileId,nil

}