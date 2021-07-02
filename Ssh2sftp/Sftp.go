/***
	ssh2 sftp
 */
package Ssh2sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

type Ssh2sftp struct{
	SftpClient   *sftp.Client
	Sshcon  *ssh.Client
}

func NewFTP() *Ssh2sftp {
	//返回日志对象
	Ssh2sftp := &Ssh2sftp{
	}
	return Ssh2sftp
}

//遍历目录
func (a *Ssh2sftp) GetList(path string)([]os.FileInfo,error) {
	fstat,err :=a.SftpClient.ReadDir(path)
	if err!=nil  {
		return  nil, err
	}
	return  fstat, nil
}
//连接
func (a *Ssh2sftp) Connect(user, password, host string, port int) error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}


	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if a.Sshcon, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	// create sftp client
	if a.SftpClient, err = sftp.NewClient(a.Sshcon); err != nil {
		return err
	}
	return nil
}

/**
断点续传
 */
func (a *Ssh2sftp) DownloadResumeFile(remotefilepath string,localpath string) error {
	var f *os.File
	var err error
	f, err = os.OpenFile(localpath, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()

	if err != nil {
		return err
	}

	var stat os.FileInfo
	stat, err = f.Stat()
	if err != nil{
		return err
	}

	offset := stat.Size()

	srcFile, err := a.SftpClient.Open(remotefilepath)
	if err != nil {
		return  err
	}
	defer srcFile.Close()
	_,err=srcFile.Seek(offset,0)

	if err!=nil{
		return err
	}
	dstFile, err := os.OpenFile(localpath, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return  err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return  err
	}

	return err
}

//获取大小
func (a *Ssh2sftp) Size(remotefilepath string) (int,error) {

	srcFile, err := a.SftpClient.Open(remotefilepath)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	stat, err :=srcFile.Stat()
	if err != nil{
		return 0, err
	}

	return  int(stat.Size()),nil
}
//断开
func (a *Ssh2sftp) Quit() error{
	a.SftpClient.Close()
	a.Sshcon.Close()
	return  nil
}