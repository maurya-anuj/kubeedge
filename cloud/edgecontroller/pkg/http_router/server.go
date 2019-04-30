package http_router

import (
	"fmt"
	"net/http"
	"io/ioutil"

	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/beehive/pkg/common/log"
	"github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/cloud/edgecontroller/pkg/controller/messagelayer"
	"github.com/kubeedge/kubeedge/edge/pkg/servicebus/util"
)

const nodeID  = "fb4ebb70-2783-42b8-b3ef-63e2fd6d242e"
const demoPort = 8000
const serverPort  = 8001
const messageSource  = "router_rest"
const routeGroup = "user"
var cml messagelayer.ContextMessageLayer
var httpRequest util.HTTPRequest

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		header := r.Header
		log.LOGGER.Infof("header: %+v", header)
		log.LOGGER.Infof("name: %s", header.Get("name"))
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.LOGGER.Infof("data=%+v \n", data)
		msg := buildMessage(data, header)
		if err := cml.Send(*msg); err != nil {
			log.LOGGER.Warnf("send message failed with error: %s, operation: %s, resource: %s", err, msg.GetOperation(), msg.GetResource())
		} else {
			log.LOGGER.Infof("send message successfully, operation: %s, resource: %s", msg.GetOperation(), msg.GetResource())
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main(c *context.Context) {
	cml = messagelayer.ContextMessageLayer{SendModuleName: "cloudhub", ReceiveModuleName: "", ResponseModuleName: "", Context: c}
	http.HandleFunc("/", handler)
	log.LOGGER.Infof("SERVER STARTED\n")
	port := fmt.Sprintf(":%d", serverPort)
	err := http.ListenAndServe(port, nil)
	if err!=nil{
		log.LOGGER.Fatalf("error listen ans serve :%e", err)
	}
}

func buildMessage(data []byte, header http.Header) *model.Message{
	msg := model.NewMessage("")
	path := fmt.Sprintf("%d:%s",demoPort, header.Get("name"))
	log.LOGGER.Infof("path:%s", path)
	resource := fmt.Sprintf("%s%s%s%s%s", "node", "/", nodeID, "/", path)
	msg.BuildRouter(messageSource, routeGroup, resource, "POST")
	var err error
	httpRequest.Body = data
	if err != nil {
		log.LOGGER.Warnf("build message resource failed with error: %s", err)
	}
	msg.FillBody(httpRequest)
	// log.LOGGER.Infof("message = %+v , httprequest : %+v", msg, httpRequest)
	return msg
}