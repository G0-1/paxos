package deal

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/n0_1/paxos/protocol"
	"github.com/n0_1/paxos/storage"
)

const IDLength = 19
const SeqStep = 10000

var Data map[string]*storage.SeqC
var Cmx sync.Mutex

func init() {
	var err error
	Data, err = storage.LoadSeq("./seq.dat")
	if err != nil {
		log.Fatal(err)
	}
}

/* 100->fail; 101->success; 102->exist     lock defant over*/
func CreateSeq(strID string) []byte {
	Cmx.Lock()
	defer Cmx.Unlock()
	_, isExist := Data[strID]
	if isExist {
		return protocol.AsmResp(102, 0, "ID exist")
	} else {
		Data[strID] = &storage.SeqC{CurSeq: 1, MidSeq: 1}
		return protocol.AsmResp(101, 1, "create seq success")
	}
}

/* 100->fail, no exist */
func AddSeq(strID string) []byte {
	Cmx.Lock()
	defer Cmx.Unlock()
	_, isExist := Data[strID]
	if isExist {
		curseq := atomic.AddUint64(&Data[strID].CurSeq, 1)
		if Data[strID].CurSeq >= Data[strID].MidSeq {
			Data[strID].MidSeq += SeqStep
			SaveSeq(Data, "./seq.dat")
		}
		return protocol.AsmResp(101, curseq, "inc seq success")
	} else {
		return protocol.AsmResp(100, 0, "ID no exist")
	}
}

/* 100->fail, no exist */
func GetSeq(strID string) []byte {
	Cmx.Lock()
	defer Cmx.Unlock()
	_, isExist := Data[strID]
	if isExist {
		return protocol.AsmResp(101, atomic.LoadUint64(&Data[strID].CurSeq), "get seq success")
	} else {
		return protocol.AsmResp(100, 0, "ID no exist")
	}
}

// when CurSeq == MidSeq call this func
func SaveSeq(seqMap map[string]*storage.SeqC, seqPath string) error {
	if seqMap == nil || seqPath == "" {
		log.Fatal("save seq error, seqmap or seqpath is nil ")
		return nil
	}
	seqData := make([]storage.JsonStruct, 0)
	for id, vo := range seqMap {
		intID, _ := strconv.Atoi(id)
		seqData = append(seqData, storage.JsonStruct{ID: uint64(intID), MidSeq: vo.MidSeq})
	}

	seqJson, jmerr := json.Marshal(&seqData)
	if jmerr != nil {
		log.Fatal("seq json marshal error: ", jmerr)
		return jmerr
	}
	seqJson = append(seqJson, '\n')

	file, oerr := os.Create(seqPath)
	defer file.Close()
	if oerr != nil {
		log.Fatal("open seq file error: ", oerr)
		return oerr
	}

	_, werr := file.Write(seqJson)
	if werr != nil {
		log.Fatal("write seq file error: ", werr)
		return werr
	}

	return nil
}
