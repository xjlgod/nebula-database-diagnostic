package remote

import (
	"github.com/pkg/sftp"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
)

// TODO
func GetFileInRemotePath(scid string, conf config.SSHConfig, remotePath string, localPath string) error {

	client, err := GetSSHClient(scid, conf)

	ftpClient, err := GetFtpClient(client.Client)
	if err != nil {
		return err
	}
	defer ftpClient.Close()

	src, err := ftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func GetFilesInRemoteDir(scid string, conf config.SSHConfig, remoteDir string, localDir string) error {

	client, err := GetSSHClient(scid, conf)

	ftpClient, err := GetFtpClient(client.Client)
	if err != nil {
		return err
	}
	defer ftpClient.Close()

	p, _ := filepath.Abs(localDir)
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		os.Mkdir(p, os.ModePerm)
	}



	filesInfo, err := ftpClient.ReadDir(remoteDir)
	for _, fileInfo := range filesInfo {

		if fileInfo.IsDir() {
			continue
		}
		srcPath := remoteDir + "/" + fileInfo.Name()
		src, err := ftpClient.OpenFile(srcPath, os.O_RDONLY)
		if err != nil {
			return err
		}
		defer src.Close()
		dstPath := filepath.Join(localDir, fileInfo.Name())
		dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			continue
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			continue
		}

	}

	return nil
}

func GetFtpClient(client *ssh.Client) (*sftp.Client, error) {
	ftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	return ftpClient, nil
}
