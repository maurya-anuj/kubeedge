package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

//KubeCloudInstTool embedes Common struct
//It implements ToolsInstaller interface
type KubeCloudInstTool struct {
	Common
}

const hostname = "master"

//InstallTools downloads KubeEdge for the specified version
//and makes the required configuration changes and initiates edgecontroller.
func (cu *KubeCloudInstTool) InstallTools() error {
	cu.SetOSInterface(GetOSInterface())
	cu.SetKubeEdgeVersion(cu.ToolVersion)

	err := cu.InstallKubeEdge()
	if err != nil {
		return err
	}

	err = cu.xhostandhostname()
	if err != nil {
		return err
	}
	//kubeadm init
	err = cu.StartK8Scluster()
	if err != nil {
		return err
	}
	//update maifests
	err = cu.updateManifests()
	if err != nil {
		return err
	}

	err=cu.clonebeehive()
	if err != nil {
		return err
	}
	fmt.Println("clonebeehive")

	err=cu.cloningkubeedge()
	if err != nil {
		return err
	}
	fmt.Println("cloningkubeedge")

	err=cu.cloningdemo()
	if err != nil {
		return err
	}
	fmt.Println("cloningdemo")

	err=cu.copycerts()
	if err != nil {
		return err
	}
	fmt.Println("cloningdemo")

	err=cu.makemasterready()
	if err != nil {
		return err
	}
	fmt.Println("makemasterready")

	// check master ready

	err=cu.removenoscheduletaint()
	if err != nil {
		return err
	}
	fmt.Println("removenoscheduletaint")

	//time.Sleep(1 * time.Second)
	//err = cu.RunEdgeController()
	//if err != nil {
	//	return err
	//}
	//fmt.Println("Edgecontroller started")

	err=cu.applynode()
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	fmt.Println("applynode")

	// check node ready

	err=cu.applylabel()
	if err != nil {
		return err
	}
	fmt.Println("applylabel")

	err=cu.copydata()
	if err != nil {
		return err
	}
	fmt.Println("copydata")

	//err=cu.buildimage()
	//if err != nil {
	//	return err
	//}
	//fmt.Println("buildimage")

	return nil
}

func (cu *KubeCloudInstTool)clonebeehive() error {
	err := os.MkdirAll("/home/src/github.com/kubeedge",0666)
	if err != nil {
		fmt.Println("in error")
		return fmt.Errorf("%e", err)

	}
	if _, err := os.Stat("/home/src/github.com/kubeedge/beehive"); os.IsNotExist(err) {
		cmd := &Command{Cmd: exec.Command("sh", "-c", "cd /home/src/github.com/kubeedge && git clone https://github.com/kubeedge/beehive.git")}
		err = cmd.ExecuteCmdShowOutput()
		errout := cmd.GetStdErr()
		if err != nil || errout != "" {
			fmt.Println("in error")
			return fmt.Errorf("%s", errout)

		}
	}
	return nil
}



func (cu *KubeCloudInstTool)cloningkubeedge() error {
	err := os.MkdirAll("/home/src/github.com/kubeedge",0666)
	if err != nil {
		fmt.Println("in error")
		return fmt.Errorf("%e", err)

	}
	if _, err := os.Stat("/home/src/github.com/kubeedge/kubeedge"); os.IsNotExist(err) {
		fmt.Println("befor cloning")
		cmd := &Command{Cmd: exec.Command("sh", "-c", " cd /home/src/github.com/kubeedge && git clone https://github.com/maurya-anuj/kubeedge.git -b demo2")}
		err = cmd.ExecuteCmdShowOutput()
		errout := cmd.GetStdErr()
		if err != nil || errout != "" {
			fmt.Println("in error")
			return fmt.Errorf("%s", errout)
		}

	}
	return nil
}

func (cu *KubeCloudInstTool)cloningdemo() error {
	if _, err := os.Stat("/home/src/fr_demo"); os.IsNotExist(err) {
		cmd := &Command{Cmd: exec.Command("sh", "-c", "cd /home/src && git clone https://github.com/maurya-anuj/fr_demo.git")}
		err := cmd.ExecuteCmdShowOutput()
		errout := cmd.GetStdErr()
		if err != nil || errout != "" {
			fmt.Println("in error")
			return fmt.Errorf("%s", errout)

		}
	}

	//cmd := &Command{Cmd: exec.Command("sh", "-c", "cp -rf /home/src/fr_demo/http_router /home/src/github.com/kubeedge/kubeedge/cloud/pkg")}
	//err := cmd.ExecuteCmdShowOutput()
	//errout := cmd.GetStdErr()
	//if err != nil || errout != "" {
	//	fmt.Println("in error")
	//	return fmt.Errorf("%s", errout)
	//
	//}
	//cmd = &Command{Cmd: exec.Command("sh", "-c", "sed -i 's|enabled: .*|enabled: [devicecontroller, controller, cloudhub, http_router]|g' /home/src/github.com/kubeedge/kubeedge/cloud/conf/modules.yaml")}
	//err = cmd.ExecuteCmdShowOutput()
	//errout = cmd.GetStdErr()
	//if err != nil || errout != "" {
	//	fmt.Println("in error")
	//	return fmt.Errorf("%s", errout)
	//
	//}
	//
	return nil
}

func (cu *KubeCloudInstTool)copycerts() error {
	cmd := &Command{Cmd: exec.Command("sh", "-c", "cp -rf /home/src/fr_demo/certs /etc/kubeedge/")}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)
	}
	return nil
}

//updateManifests - Kubernetes Manifests file will be updated by necessary parameters
func (cu *KubeCloudInstTool) updateManifests() error {
	input, err := ioutil.ReadFile(KubeCloudApiserverYamlPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	output := bytes.Replace(input, []byte("insecure-port=0"), []byte("insecure-port=8080"), -1)

	if err = ioutil.WriteFile(KubeCloudApiserverYamlPath, output, 0666); err != nil {
		fmt.Println(err)
		return err
	}

	lines, err := file2lines(KubeCloudApiserverYamlPath)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == KubeCloudReplaceIndex {
			fileContent += KubeCloudReplaceString
		}
		fileContent += line
		fileContent += "\n"
	}

	ioutil.WriteFile(KubeCloudApiserverYamlPath, []byte(fileContent), 0644)

	if err = ioutil.WriteFile(KubeEdgeControllerYaml, ControllerYaml, 0666); err != nil {
		return err
	}
	if err = ioutil.WriteFile(KubeEdgeControllerModulesYaml,ControllerModulesYaml, 0666); err != nil {
		return err
	}
	return nil

}

func file2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
//
////RunEdgeController starts edgecontroller process
//func (cu *KubeCloudInstTool) RunEdgeController() error {
//
//	cmd := &Command{Cmd: exec.Command("sh", "-c", "cd /home/src/github.com/kubeedge/kubeedge && make all WHAT=cloud")}
//	err := cmd.ExecuteCmdShowOutput()
//	errout := cmd.GetStdErr()
//	if err != nil || errout != "" {
//		fmt.Println("in error")
//		return fmt.Errorf("%s", errout)
//
//	}
//
//	cmd = &Command{Cmd: exec.Command("sh", "-c", "cp /home/src/github.com/kubeedge/kubeedge/cloud/edgecontroller /etc/kubeedge/kubeedge/cloud/")}
//	err = cmd.ExecuteCmdShowOutput()
//	errout = cmd.GetStdErr()
//	if err != nil || errout != "" {
//		fmt.Println("in error")
//		return fmt.Errorf("%s", errout)
//
//	}
//
//	binExec := fmt.Sprintf("cd  chmod +x edgecontroller && ./edgecontroller > %s/kubeedge/cloud/%s.log 2>&1 &", KubeEdgePath, KubeCloudBinaryName)
//	cmd = &Command{Cmd: exec.Command("sh", "-c", binExec)}
//	cmd.Cmd.Env = os.Environ()
//	env := fmt.Sprintf("GOARCHAIUS_CONFIG_PATH=%skubeedge/cloud", KubeEdgePath)
//	cmd.Cmd.Env = append(cmd.Cmd.Env, env)
//	err = cmd.ExecuteCmdShowOutput()
//	errout = cmd.GetStdErr()
//	if err != nil || errout != "" {
//		return fmt.Errorf("%s", errout)
//	}
//	fmt.Println(cmd.GetStdOutput())
//	fmt.Println("KubeEdge controller is running, For logs visit", KubeEdgePath+"cloud/")
//	return nil
//}
func (cu *KubeCloudInstTool) RunEdgeController() error {
	cmd := &Command{Cmd: exec.Command("sh", "-c", "cd /home/src/github.com/kubeedge/kubeedge && make all WHAT=cloud")}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}

	//cmd = &Command{Cmd: exec.Command("sh", "-c", "cp -rf /home/src/github.com/kubeedge/kubeedge/cloud/edgecontroller /etc/kubeedge/kubeedge/cloud/")}
	//err = cmd.ExecuteCmdShowOutput()
	//errout = cmd.GetStdErr()
	//if err != nil || errout != "" {
	//	fmt.Println("in error")
	//	return fmt.Errorf("%s", errout)
	//
	//}
	time.Sleep(1*time.Second)
	//filetoCopy := fmt.Sprintf("cp -rf %s/kubeedge/cloud/%s /usr/local/bin/", KubeEdgePath, KubeCloudBinaryName)
	//cmd = &Command{Cmd: exec.Command("sh", "-c", filetoCopy)}
	//err = cmd.ExecuteCmdShowOutput()
	//errout = cmd.GetStdErr()
	//if err != nil || errout != "" {
	//	fmt.Println("in error")
	//	return fmt.Errorf("%s", errout)
	//
	//}
	cmd = &Command{Cmd: exec.Command("sh", "-c", "export PATH=$PATH:/home/src/github.com/kubeedge/kubeedge/cloud")}
	err = cmd.ExecuteCmdShowOutput()
	errout = cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}

	binExec := fmt.Sprintf("chmod +x /home/src/github.com/kubeedge/kubeedge/cloud/%s && /home/src/github.com/kubeedge/kubeedge/cloud/%s > /home/src/github.com/kubeedge/kubeedge/cloud/%s.log 2>&1 &", KubeCloudBinaryName, KubeCloudBinaryName,  KubeCloudBinaryName)
	cmd = &Command{Cmd: exec.Command("sh", "-c", binExec)}
	cmd.Cmd.Env = os.Environ()
	env := fmt.Sprintf("GOARCHAIUS_CONFIG_PATH=/home/src/github.com/kubeedge/kubeedge/cloud", KubeEdgePath)
	cmd.Cmd.Env = append(cmd.Cmd.Env, env)
	err = cmd.ExecuteCmdShowOutput()
	errout = cmd.GetStdErr()
	if err != nil || errout != "" {
		return fmt.Errorf("%s", errout)
	}
	fmt.Println(cmd.GetStdOutput())
	fmt.Println("KubeEdge controller is running, For logs visit", KubeEdgePath + "cloud/")
	return nil
}

//TearDown method will remove the edge node from api-server and stop edgecontroller process
func (cu *KubeCloudInstTool) TearDown() error {

        cu.SetOSInterface(GetOSInterface())

	//Stops kubeadm
	binExec := fmt.Sprintf("echo 'y' | kubeadm reset &&  rm -rf ~/.kube")
	cmd := &Command{Cmd: exec.Command("sh", "-c", binExec)}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		return fmt.Errorf("kubeadm reset failed %s", errout)
	}

	//Kill edgecontroller process
	cu.KillKubeEdgeBinary(KubeCloudBinaryName)

	return nil
}

func (cu *KubeCloudInstTool) xhostandhostname() error {
	cmd := &Command{Cmd: exec.Command("sh", "-c", "xhost +local:user:root")}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)
	}
	cmd = &Command{Cmd: exec.Command("sh", "-c", fmt.Sprintf("hostnamectl set-hostname %s", hostname))}
	err = cmd.ExecuteCmdShowOutput()
	errout = cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)
	}
	return nil
}

func (cu *KubeCloudInstTool) makemasterready() error {
	time.Sleep(1 * time.Second)
	cmd := &Command{Cmd: exec.Command("sh", "-c", "kubectl apply -f /home/src/fr_demo/kube-flannel.yml")}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}
	time.Sleep(5 * time.Second)
	cmd = &Command{Cmd: exec.Command("sh", "-c", "kubectl apply -f /home/src/fr_demo/kube-flannel-rbac.yml")}
	err = cmd.ExecuteCmdShowOutput()
	errout = cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}
	return nil
}

func (cu *KubeCloudInstTool) removenoscheduletaint() error {
	//cmd := &Command{Cmd: exec.Command("sh", "-c", "hostname")}
	//err := cmd.ExecuteCmdShowOutput()
	masternode := hostname
	fmt.Printf("master : %s", masternode)
	cmd := &Command{Cmd: exec.Command("sh", "-c", fmt.Sprintf("kubectl taint nodes %s node-role.kubernetes.io/master:NoSchedule-", masternode))}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}
	return nil
}

func (cu *KubeCloudInstTool) applynode() error {
	//kubectl apply -f $GOPATH/src/github.com/kubeedge/kubeedge/build/node.json -s http://192.168.20.50:8080
	nodeJSONApply := fmt.Sprintf("kubectl apply -f /home/src/github.com/kubeedge/kubeedge/build/node.json")
	cmd := &Command{Cmd: exec.Command("sh", "-c", nodeJSONApply)}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
	return fmt.Errorf("%s", errout)
	}
	fmt.Println(cmd.GetStdOutput())

	return nil
}

func (cu *KubeCloudInstTool) applylabel() error {
	//cmd := &Command{Cmd: exec.Command("sh", "-c", "hostname")}
	masternode := hostname

	cmd := &Command{Cmd: exec.Command("sh", "-c", fmt.Sprintf("kubectl label nodes %s dedicated=master", masternode))}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}
	nodename := KubeEdgeToReplaceKey1
	cmd = &Command{Cmd: exec.Command("sh", "-c", fmt.Sprintf("kubectl label nodes %s dedicated=edge", nodename))}
	err = cmd.ExecuteCmdShowOutput()
	errout = cmd.GetStdErr()
	if err != nil || errout != "" {
		fmt.Println("in error")
		return fmt.Errorf("%s", errout)

	}
	return nil
}


func (cu *KubeCloudInstTool) copydata() error {
	var cmd *Command
	if _, err := os.Stat("/etc/kubeedge/data"); os.IsNotExist(err) {
		cmd = &Command{Cmd: exec.Command("sh", "-c", "cp -rf /home/src/fr_demo/demo/data /etc/kubeedge")}
	}else {
		cmd = &Command{Cmd: exec.Command("sh", "-c", "rm -rf /etc/kubeedge/data && cp -rf /home/src/fr_demo/demo/data /etc/kubeedge")}
	}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		return fmt.Errorf("%s", errout)
	}
	fmt.Println(cmd.GetStdOutput())

	return nil
}

func (cu *KubeCloudInstTool) buildimage() error {

	cmd := &Command{Cmd: exec.Command("sh", "-c", "cd /home/src/fr_demo/demo/demo_docker && docker build -t demo-docker:1 ./")}
	err := cmd.ExecuteCmdShowOutput()
	errout := cmd.GetStdErr()
	if err != nil || errout != "" {
		return fmt.Errorf("%s", errout)
	}
	fmt.Println(cmd.GetStdOutput())

	cmd = &Command{Cmd: exec.Command("sh", "-c", "cd /home/src/fr_demo/demo/create_docker && docker build -t training:1 ./")}
	err = cmd.ExecuteCmdShowOutput()
	errout = cmd.GetStdErr()
	if err != nil || errout != "" {
		return fmt.Errorf("%s", errout)
	}
	fmt.Println(cmd.GetStdOutput())

	return nil
}
