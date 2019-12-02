package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

//https://www.youtube.com/watch?v=zVqczFZr124

type Transaccion struct {
	Monto   float64
	Origen  string
	Destino string
}

type Block struct {
	Index        int
	Timestamp    time.Time
	Data         Transaccion
	PreviousHash string
	Hash         string
}

func (block *Block) CalculateHash() string {
	src := fmt.Sprintf("%d-%s-%s", block.Index, block.Timestamp.String(), block.Data)
	return base64.StdEncoding.EncodeToString([]byte(src))
}

type BlockChain struct {
	Chain []Block
}

func (blockChain *BlockChain) CreateGenesisBlock() Block {
	block := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         Transaccion{},
		PreviousHash: "0",
	}
	block.Hash = block.CalculateHash()
	return block
}

func (blockChain *BlockChain) GetLatesBlock() Block {
	n := len(blockChain.Chain)
	return blockChain.Chain[n-1]
}

func (blockChain *BlockChain) AddBlock(block Block) {
	block.Timestamp = time.Now()
	block.Index = blockChain.GetLatesBlock().Index + 1
	block.PreviousHash = blockChain.GetLatesBlock().Hash
	block.Hash = block.CalculateHash()
	blockChain.Chain = append(blockChain.Chain, block)
}

func (blockChain *BlockChain) IsChainValid() bool {
	n := len(blockChain.Chain)
	for i := 1; i < n; i++ {
		currentBlock := blockChain.Chain[i]
		previousBlock := blockChain.Chain[i-1]
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

func CreateBlockChain() BlockChain {
	bc := BlockChain{}
	genesisBlock := bc.CreateGenesisBlock()
	bc.Chain = append(bc.Chain, genesisBlock)
	return bc
}

var localBlockChain BlockChain

func PrintTransacciones() {
	blocks := localBlockChain.Chain[1:]
	for index, block := range blocks {
		transaccion := block.Data
		]

		fmt.Printf("\t\nOrigen:  %s", transaccion.Origen)
		fmt.Printf("\t\nDestino:  %s", transaccion.Destino)
		fmt.Printf("\t\nMonto:  %f", transaccion.Monto)

		//transaccionBytes, _ := json.MarshalIndent(transaccion, "", "\t")
		//jsonTransaccion := string(transaccionBytes)
		//fmt.Println(jsonTransaccion)
	}
}

func PrintHosts() {
	fmt.Println("- - - HOSTS - - -")
	const first = 0
	fmt.Printf("\t%s (Mi host)\n", LOCALHOST)
	for _, host := range HOSTS {
		fmt.Printf("\t%s\n", host)
	}
}

func main() {
	var dest string
	end := make(chan int)
	updatedBlocks := make(chan int)
	fmt.Print("Ingrese su host: ")
	fmt.Scanf("%s\n", &LOCALHOST)
	fmt.Print("Ingrese host de destino(Vacio por ser el primer nodo): ")
	fmt.Scanf("%s\n", &dest)
	go BCIPServer(end, updatedBlocks)
	localBlockChain = CreateBlockChain()
	if dest != "" {
		requestBody := &RequestBody{
			Message:     LOCALHOST,
			MessageType: NEWHOST,
		}
		requestMessage, _ := json.Marshal(requestBody)
		SendMessage(dest, string(requestMessage))
		requestBody.MessageType = NEWBLOCK
		requestMessage, _ = json.Marshal(requestBody)
		SendMessage(dest, string(requestMessage))
		<-updatedBlocks
	}
	var action int
	fmt.Println("App de Transaciones Transparentes para los ciudadnos")
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n\n1. Crear Transaccion\n2. Lista de Transacciones\n3. Lista de Hosts\n")
		fmt.Print("Ingrese opcion: 1 o 2 o 3):")
		fmt.Scanf("%d\n", &action)
		if action == NEWTRANS {
			transaccion := Transaccion{}
			fmt.Println("- - - Registro - - -")

			fmt.Print("Ingrese Origen: ")
			transaccion.Origen, _ = in.ReadString('\n')

			fmt.Print("Ingrese Destino: ")
			transaccion.Destino, _ = in.ReadString('\n')

			fmt.Print("Ingrese Monto: ")
			fmt.Scanf("%f", &transaccion.Monto)

			newBlock := Block{
				Data: transaccion,
			}
			localBlockChain.AddBlock(newBlock)
			BroadcastBlock(newBlock)
			fmt.Println("Registro Satisfactorio")
			time.Sleep(2 * time.Second)
			PrintTransacciones()
		} else if action == LISTTRANS {
			PrintTransacciones()
		} else if action == LISTHOSTS {
			PrintHosts()
		}
	}
	<-end
}

var HOSTS []string
var LOCALHOST string

type MessageType int32

const (
	NEWHOST   MessageType = 0
	ADDHOST   MessageType = 1
	ADDBLOCK  MessageType = 2
	NEWBLOCK  MessageType = 3
	SETBLOCKS MessageType = 4
	PROTOCOL              = "tcp"
	NEWTRANS              = 1
	LISTTRANS             = 2
	LISTHOSTS             = 3
)

type RequestBody struct {
	Message     string
	MessageType MessageType
}

func GetMessage(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	data, _ := reader.ReadString('\n')
	return strings.TrimSpace(data)
}

func SendMessage(toHost string, message string) {
	conn, _ := net.Dial(PROTOCOL, toHost)
	defer conn.Close()
	fmt.Fprintln(conn, message)
}

func SendMessageWithReply(toHost string, message string) string {
	conn, _ := net.Dial(PROTOCOL, toHost)
	defer conn.Close()
	fmt.Fprintln(conn, message)
	return GetMessage(conn)
}

func RemoveHost(index int, hosts []string) []string {
	n := len(hosts)
	hosts[index] = hosts[n-1]
	hosts[n-1] = ""
	return hosts[:n-1]
}

func RemoveHostByValue(ip string, hosts []string) []string {
	for index, host := range hosts {
		if host == ip {
			return RemoveHost(index, hosts)
		}
	}
	return hosts
}

func Broadcast(newHost string) {
	for _, host := range HOSTS {
		data := append(HOSTS, newHost, LOCALHOST)
		data = RemoveHostByValue(host, data)
		requestBroadcast := RequestBody{
			Message:     strings.Join(data, ","),
			MessageType: ADDHOST,
		}
		broadcastMessage, _ := json.Marshal(requestBroadcast)
		SendMessage(host, string(broadcastMessage))
	}
}

func BroadcastBlock(newBlock Block) {
	for _, host := range HOSTS {
		data, _ := json.Marshal(newBlock)
		requestBroadcast := RequestBody{
			Message:     string(data),
			MessageType: ADDBLOCK,
		}
		broadcastMessage, _ := json.Marshal(requestBroadcast)
		SendMessage(host, string(broadcastMessage))
	}
}

func BCIPServer(end chan<- int, updatedBlocks chan<- int) {
	ln, _ := net.Listen(PROTOCOL, LOCALHOST)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		defer conn.Close()
		request := RequestBody{}
		data := GetMessage(conn)
		_ = json.Unmarshal([]byte(data), &request)
		if request.MessageType == NEWHOST {
			message := strings.Join(append(HOSTS, LOCALHOST), ",")
			requestClient := RequestBody{
				Message:     message,
				MessageType: ADDHOST,
			}
			clientMessage, _ := json.Marshal(requestClient)
			SendMessage(request.Message, string(clientMessage))
			Broadcast(request.Message)
			HOSTS = append(HOSTS, request.Message)
		} else if request.MessageType == ADDHOST {
			HOSTS = strings.Split(request.Message, ",")
		} else if request.MessageType == NEWBLOCK {
			blocksMessage, _ := json.Marshal(localBlockChain.Chain)
			setBlocksRequest := RequestBody{
				Message:     string(blocksMessage),
				MessageType: SETBLOCKS,
			}
			setBlocksMessage, _ := json.Marshal(setBlocksRequest)
			SendMessage(request.Message, string(setBlocksMessage))
		} else if request.MessageType == SETBLOCKS {
			_ = json.Unmarshal([]byte(request.Message), &localBlockChain.Chain)
			updatedBlocks <- 0
		} else if request.MessageType == ADDBLOCK {
			block := Block{}
			src := []byte(request.Message)
			json.Unmarshal(src, &block)
			localBlockChain.Chain = append(localBlockChain.Chain, block)
		}
	}
	end <- 0
}
