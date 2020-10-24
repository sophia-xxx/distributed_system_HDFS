package detector

import (
	"github.com/golang/protobuf/proto"
	"net"
	"strings"

	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/config"
	"cs425_mp2/logger"
)

// socket to listen TCP message
func ListenMessage() {
	addressString := GetLocalIPAddr().String() + config.TCPPORT
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		logger.ErrorLogger.Println("Cannot resolve TCP address!")
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP connection!")
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.ErrorLogger.Println("Cannot open TCP connection!")
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	// read message data
	buf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		logger.ErrorLogger.Println("Unable to read data!")
	}
	messageBytes := buf[0:n]
	remoteMsg, _ := DecodeTCPMessage(messageBytes)

	// deal with all PUT relevant message
	if remoteMsg.Type <= config.PUT {
		putMessageHandler(remoteMsg)
	}
	// deal with all GET relevant message
	if remoteMsg.Type > config.PUT && remoteMsg.Type <= config.GET {
		getMessageHandler(remoteMsg)
	}
	// deal with all DELETE relevant message
	if remoteMsg.Type > config.GET && remoteMsg.Type <= config.DELETE {
		deleteMessageHandle(remoteMsg)
	}
	// deal with other message
	if remoteMsg.Type > config.DELETE {
		if isIntroducer && remoteMsg.Type == pbm.MsgType_LIST {
			ListReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
		}
		if remoteMsg.Type == pbm.MsgType_LIST_REP {
			nodeList := remoteMsg.PayLoad
			var fileString strings.Builder
			for _, node := range nodeList {
				fileString.WriteString(node + "\t")
			}
			logger.PrintInfo(remoteMsg.FileName + " : " + fileString.String())
		}
	}
}

// send TCP message
func SendMessage(dest string, message []byte) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.TCPPORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote address!")
	}
	_, err = conn.Write(message)
}

func EncodeTCPMessage(fileMessage *pbm.TCPMessage) ([]byte, error) {
	message, err := proto.Marshal(fileMessage)
	return message, err
}
func DecodeTCPMessage(message []byte) (*pbm.TCPMessage, error) {
	list := &pbm.TCPMessage{}
	err := proto.Unmarshal(message, list)

	return list, err
}

/*//The get request initiator tell file source ip that it get the file size info successfully
func sendReadReqAckAck(targetIp string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_P2P_SIZE_ACK,
		FileName: sdfsFileName,
		SenderIP: detector.GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}*/