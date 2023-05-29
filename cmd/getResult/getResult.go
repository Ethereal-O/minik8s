package getResult

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"os"
	"path"
	"time"
)

var key string

var getResultCmd = &cobra.Command{
	Use:   "getResult",
	Short: "get the result of gpuJob of k8s",
	Long:  "this is the main cmd to get the the result of gpuJob in k8s",
	Run:   doit,
}

func connect(user string, password string, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func printResult(ip string, file string) {

	var (
		err        error
		sftpClient *sftp.Client
	)

	sftpClient, err = connect("root", "cloudOS2023", ip, 22)
	if err != nil {
		fmt.Println(err)
	}
	defer sftpClient.Close()

	var remoteFilePath = config.GPU_NODE_DIR_PATH + "/" + file + "/result.out"
	var localDir = "/home/tmpResult"

	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	var localFileName = path.Base(remoteFilePath)
	dstFile, err := os.Create(path.Join(localDir, localFileName))
	if err != nil {
		fmt.Println(err)
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		fmt.Println(err)
	}
	content, _ := ioutil.ReadFile(dstFile.Name())
	fmt.Println(string(content))

}

func doit(cmd *cobra.Command, args []string) {
	var gpuObject object.GpuJob
	gpujob := client.Get_object(key, config.GPUJOB_TYPE)[0]
	json.Unmarshal([]byte(gpujob), &gpuObject)

	var pod = client.Get_object(object.GpuJobPodFullName(gpuObject), config.POD_TYPE)[0]
	var podObject object.Pod
	json.Unmarshal([]byte(pod), &podObject)

	ip := podObject.Runtime.Bind[5:]
	file := gpuObject.Metadata.Name
	printResult(ip, file)
}

func init() {
	getResultCmd.MarkFlagRequired("key")
	getResultCmd.Flags().StringVarP(&key, "key", "k", config.EMPTY_FLAG, "Name of API object to inspect, refers to all API objects of specified type if not set")
}

func GetResult() *cobra.Command {
	return getResultCmd
}
