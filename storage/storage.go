package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

type SeqC struct {
	CurSeq uint64
	MidSeq uint64
}

var seqMap map[string]*SeqC

type JsonStruct struct {
	ID     uint64 `json:id`
	MidSeq uint64 `json:midseq`
}

func LoadSeq(seqPath string) (map[string]*SeqC, error) {
	var fileNoExist bool
	seqMap = make(map[string]*SeqC, 0)
	seqFileName := "seq.dat"

	fileinfo, serr := os.Stat(seqPath)
	if serr != nil {
		if serr, ok := serr.(*os.PathError); ok && os.IsNotExist(serr.Err) {
			fileNoExist = true
		} else {
			fileNoExist = false
		}
	}

	// condition 1 ,stat error, file exist
	if serr != nil && fileNoExist == false {
		log.Fatal("Get file stat error: ", serr)
		return nil, serr
	}
	// condition 2, stat error, file no exist
	if serr != nil && fileNoExist == true {
		_, cerr := os.Create(seqPath)
		if cerr != nil {
			log.Fatal("seq file no exist, create new one occur error: ", cerr)
			return nil, cerr
		}
		seqMap = make(map[string]*SeqC, 0)
		log.Println("seq file no exist, create new one")
		return seqMap, nil
	}
	// condition 3, stat ok, file is dir
	if fileinfo.IsDir() {
		seqPath = seqPath + string(os.PathSeparator) + seqFileName
		_, cerr := os.Create(seqPath)
		if cerr != nil {
			log.Fatal("seq path is dir, create new one occur error: ", cerr)
			return nil, cerr
		}
		log.Println("seq file is dir, create new one")
		seqMap = make(map[string]*SeqC, 0)
		return seqMap, nil
	}
	//condition 4, stat ok, file ok
	f, oerr := os.Open(seqPath)
	defer f.Close()
	if oerr != nil {
		log.Fatal("open seq path error: ", oerr)
		return nil, oerr
	}
	r := bufio.NewReader(f)
	seqJson, rerr := r.ReadBytes('\n')
	if rerr != nil && rerr != io.EOF {
		log.Fatal("read seq path error: ", rerr)
	}
	if rerr == io.EOF {
		seqMap = make(map[string]*SeqC, 0)
		log.Println("seq file was empty")
		return seqMap, nil
	}

	seqData := make([]JsonStruct, 0)
	juerr := json.Unmarshal(seqJson, &seqData)
	if juerr != nil {
		log.Fatal("ummarshal seq json error: ", juerr)
		return nil, juerr
	}
	for _, vo := range seqData {
		seqMap[strconv.Itoa(int(vo.ID))] = &SeqC{CurSeq: vo.MidSeq, MidSeq: vo.MidSeq}
	}
	log.Println("log seqMap data: ", seqMap)
	return seqMap, nil
}
