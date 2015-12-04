package fdfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ctripcorp/nephele/Godeps/_workspace/src/github.com/ctripcorp/ghost/pool"
	"net"
	"os"
	"time"
)

const (
	STORAGE_MIN_CONN        int           = 5
	STORAGE_MAX_CONN        int           = 5
	STORAGE_MAX_IDLE        time.Duration = 10 * time.Hour
	STORAGE_NETWORK_TIMEOUT time.Duration = 10 * time.Second
)

type storageClient struct {
	host string
	port int
	pool.Pool
}

func newStorageClient(host string, port int) (*storageClient, error) {
	client := &storageClient{host: host, port: port}
	p, err := pool.NewBlockingPool(STORAGE_MIN_CONN, STORAGE_MAX_CONN, STORAGE_MAX_IDLE, client.makeConn)
	if err != nil {
		return nil, err
	}
	client.Pool = p
	return client, nil

}

func (this *storageClient) storageDownload(storeInfo *storageInfo, offset int64, downloadSize int64, fileName string) ([]byte, error) {
	//get a connetion from pool
	conn, err := this.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//request header
	buffer := new(bytes.Buffer)
	//package length:file_offset(8)  download_bytes(8)  group_name(16)  file_name(n)
	binary.Write(buffer, binary.BigEndian, int64(32+len(fileName)))
	//cmd
	buffer.WriteByte(byte(STORAGE_PROTO_CMD_DOWNLOAD_FILE))
	//status
	buffer.WriteByte(byte(0))
	//offset
	binary.Write(buffer, binary.BigEndian, offset)
	//download bytes
	binary.Write(buffer, binary.BigEndian, downloadSize)
	//16 bit groupName
	groupNameBytes := bytes.NewBufferString(storeInfo.groupName).Bytes()
	for i := 0; i < 15; i++ {
		if i >= len(groupNameBytes) {
			buffer.WriteByte(byte(0))
		} else {
			buffer.WriteByte(groupNameBytes[i])
		}
	}
	buffer.WriteByte(byte(0))
	// fileNameLen bit fileName
	fileNameBytes := bytes.NewBufferString(fileName).Bytes()
	for i := 0; i < len(fileNameBytes); i++ {
		buffer.WriteByte(fileNameBytes[i])
	}
	//send request
	if err := tcpSend(conn, buffer.Bytes(), STORAGE_NETWORK_TIMEOUT); err != nil {
		return nil, errors.New(fmt.Sprintf("send to storage server %v fail, error info: %v", conn.RemoteAddr().String(), err.Error()))
	}
	//receive response header
	recvBuff, err := recvResponse(conn, STORAGE_NETWORK_TIMEOUT)
	if err != nil {
		return nil, err
		//try again
		//	if err = tcpSend(conn, buffer.Bytes(), STORAGE_NETWORK_TIMEOUT); err != nil {
		//		return nil, err
		//	}
		//	if recvBuff, err = recvResponse(conn, STORAGE_NETWORK_TIMEOUT); err != nil {
		//		return nil, err
		//	}
	}
	return recvBuff, nil
}

func (this *storageClient) storageDeleteFile(storeInfo *storageInfo, fileName string) error {
	//get a connetion from pool
	conn, err := this.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	//request header
	buffer := new(bytes.Buffer)
	//package length:group_name(16)  file_name(n)
	binary.Write(buffer, binary.BigEndian, int64(FDFS_GROUP_NAME_MAX_LEN+len(fileName)))
	//cmd
	buffer.WriteByte(byte(STORAGE_PROTO_CMD_DELETE_FILE))
	//status
	buffer.WriteByte(byte(0))
	//16 bit groupName
	buffer.WriteString(fixString(storeInfo.groupName, FDFS_GROUP_NAME_MAX_LEN))
	// fileNameLen bit fileName
	buffer.WriteString(fileName)
	//send request
	if err := tcpSend(conn, buffer.Bytes(), STORAGE_NETWORK_TIMEOUT); err != nil {
		return errors.New(fmt.Sprintf("send to storage server %v fail, error info: %v", conn.RemoteAddr().String(), err.Error()))
	}
	//receive response header
	if _, err := recvResponse(conn, STORAGE_NETWORK_TIMEOUT); err != nil {
		return err
	}
	return nil
}

//stroage upload by buffer
func (this *storageClient) storageUploadByBuffer(storeInfo *storageInfo, fileBuffer []byte,
	fileExtName string) (string, error) {
	bufferSize := len(fileBuffer)

	return this.storageUploadFile(storeInfo, fileBuffer, int64(bufferSize),
		STORAGE_PROTO_CMD_UPLOAD_FILE, "", "", fileExtName)
}

//storage upload slave by buffer
func (this *storageClient) storageUploadSlaveByBuffer(storeInfo *storageInfo, fileBuffer []byte,
	remoteFileId string, prefixName string, fileExtName string) (string, error) {
	bufferSize := len(fileBuffer)

	return this.storageUploadFile(storeInfo, fileBuffer, int64(bufferSize),
		STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, remoteFileId, prefixName, fileExtName)
}

//stroage upload file
func (this *storageClient) storageUploadFile(storeInfo *storageInfo, fileBuff []byte, fileSize int64, cmd int8,
	masterFileName string, prefixName string, fileExtName string) (string, error) {
	var (
		uploadSlave bool  = false
		headerLen   int64 = 15
	)
	//get a connetion from pool
	conn, err := this.Get()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	masterFilenameLen := int64(len(masterFileName))
	if len(storeInfo.groupName) > 0 && len(masterFileName) > 0 {
		uploadSlave = true
		//master_len(8) file_size(8) prefix_name(16) file_ext_name(6) master_name(master_filename_len)
		headerLen = int64(38) + masterFilenameLen
	}

	//request header
	buffer := new(bytes.Buffer)
	//package length
	binary.Write(buffer, binary.BigEndian, int64(headerLen+fileSize))
	//cmd
	buffer.WriteByte(byte(cmd))
	//status
	buffer.WriteByte(byte(0))

	if uploadSlave {
		// master file name len
		binary.Write(buffer, binary.BigEndian, masterFilenameLen)
		// file size
		binary.Write(buffer, binary.BigEndian, fileSize)
		// 16 bit prefixName
		buffer.WriteString(fixString(prefixName, FDFS_FILE_PREFIX_MAX_LEN))
		// 6 bit fileExtName
		buffer.WriteString(fixString(fileExtName, FDFS_FILE_EXT_NAME_MAX_LEN))
		// master_file_name
		buffer.WriteString(masterFileName)
	} else {
		//store_path_index
		buffer.WriteByte(byte(uint8(storeInfo.storePathIndex)))
		// file size
		binary.Write(buffer, binary.BigEndian, fileSize)
		// 6 bit fileExtName
		buffer.WriteString(fixString(fileExtName, FDFS_FILE_EXT_NAME_MAX_LEN))
	}
	//send header
	if err := tcpSend(conn, buffer.Bytes(), STORAGE_NETWORK_TIMEOUT); err != nil {
		return "", errors.New(fmt.Sprintf("send to storage server %v fail, error info: %v", conn.RemoteAddr().String(), err.Error()))
	}
	//send file buffer
	if err := tcpSend(conn, fileBuff, STORAGE_NETWORK_TIMEOUT); err != nil {
		return "", errors.New(fmt.Sprintf("send to storage server %v fail, error info: %v", conn.RemoteAddr().String(), err.Error()))
	}
	//receive response header
	recvBuff, err := recvResponse(conn, STORAGE_NETWORK_TIMEOUT)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(recvBuff)
	groupName, err := readCstr(b, FDFS_GROUP_NAME_MAX_LEN)
	if err != nil {
		return "", err
	}
	remoteFilename := string(recvBuff[len(recvBuff)-b.Len():])
	remoteFileId := groupName + string(os.PathSeparator) + remoteFilename
	return remoteFileId, nil
}

//factory method used to dial
func (this *storageClient) makeConn() (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", this.host, this.port)
	event := globalCat.NewEvent("DialStorage", addr)
	defer func() {
		event.Complete()
	}()
	conn, err := net.DialTimeout("tcp", addr, STORAGE_NETWORK_TIMEOUT)
	if err != nil {
		errMsg := fmt.Sprintf("dial storage %v fail, error info: %v", addr, err.Error())
		event.AddData("detail", errMsg)
		event.SetStatus("ERROR")
		return nil, errors.New(errMsg)
	}
	event.SetStatus("0")
	return conn, nil
}
