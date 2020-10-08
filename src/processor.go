package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type compressHead struct {
	srclen, dstlen, keymapLen uint32                 //源文件字符个数  压缩文件字符个数   哈夫曼编码字符映射个数
	patchBit                  uint8                  //压缩后不足8bit补0个数
	keysMap                   map[interface{}]uint32 //字符统计构建哈夫曼树
}

//按照小端模式写入文件
func getCompressedBytes(pHead *compressHead, data []byte) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, pHead.srclen); err == nil {
		if err = binary.Write(buf, binary.LittleEndian, pHead.dstlen); err != nil {
			fmt.Println(err.Error())
		}

		if err = binary.Write(buf, binary.LittleEndian, pHead.keymapLen); err != nil {
			fmt.Println(err.Error())
		}

		if err = binary.Write(buf, binary.LittleEndian, pHead.patchBit); err != nil {
			fmt.Println(err.Error())
		}

		for key, value := range pHead.keysMap {
			if err = binary.Write(buf, binary.LittleEndian, key); err != nil {
				fmt.Println(err.Error())
			}

			if err = binary.Write(buf, binary.LittleEndian, value); err != nil {
				fmt.Println(err.Error())
			}
		}

		if err = binary.Write(buf, binary.LittleEndian, data); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
	return buf.Bytes()
}

func Compress(strInFileName, strOutFileName string) bool {
	if data, err := ioutil.ReadFile(strInFileName); err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		keys := make(map[interface{}]uint32)
		for i := 0; i < len(data); i++ {
			if _, ok := keys[data[i]]; ok {
				keys[data[i]]++
			} else {
				keys[data[i]] = 1
			}
		}

		if pTree := CreatHuffman(keys); pTree != nil {
			pHead := &compressHead{srclen: (uint32(len(data))), keysMap: keys}
			compressdata := make([]byte, 0)
			factor := make(map[byte][]byte)
			for i := 0; i < len(data); i++ {
				if value, ok := factor[data[i]]; ok {
					compressdata = append(compressdata, value...)
				} else {
					if code, codeok := GetHuffmanCode(data[i], pTree); codeok {
						factor[data[i]] = code
						compressdata = append(compressdata, code...)
					}
				}
			}

			pHead.keymapLen = (uint32(len(pHead.keysMap)))
			pHead.patchBit = (uint8)(8 - len(compressdata)%8)
			for i := uint8(0); i < pHead.patchBit; i++ {
				compressdata = append(compressdata, 0)
			}

			afterdata := make([]byte, 0)
			for ; len(compressdata) >= 8; compressdata = compressdata[8:] {
				var b byte = 0
				for i := 0; i < 8; i++ {
					b |= compressdata[i] << (7 - i)
				}

				afterdata = append(afterdata, b)
			}

			fmt.Printf("after data:%v\n", afterdata)
			pHead.dstlen = (uint32(len(afterdata)))
			buf := getCompressedBytes(pHead, afterdata)
			ioutil.WriteFile(strOutFileName, buf, 0666)
		}
	}

	return true
}

func getCompressedHead(data []byte) (*compressHead, []byte) {
	pHead := new(compressHead)
	buf := bytes.NewBuffer(data)

	binary.Read(buf, binary.LittleEndian, &pHead.srclen)
	binary.Read(buf, binary.LittleEndian, &pHead.dstlen)
	binary.Read(buf, binary.LittleEndian, &pHead.keymapLen)
	binary.Read(buf, binary.LittleEndian, &pHead.patchBit)
	pHead.keysMap = make(map[interface{}]uint32)
	for i := uint32(0); i < pHead.keymapLen; i++ {
		var key byte
		var value uint32
		binary.Read(buf, binary.LittleEndian, &key)
		binary.Read(buf, binary.LittleEndian, &value)
		pHead.keysMap[key] = value
	}

	dstdata := make([]byte, pHead.dstlen)
	binary.Read(buf, binary.LittleEndian, dstdata[:pHead.dstlen])

	return pHead, dstdata
}

func DeCompress(strInFileName, strOutFileName string) bool {
	if data, err := ioutil.ReadFile(strInFileName); err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		pHead, dst := getCompressedHead(data)
		src := make([]byte, pHead.srclen)
		srcindex := 0
		if pTree := CreatHuffman(pHead.keysMap); pTree != nil {
			pTmpTree := pTree
			bEnd := false
			for len(dst) >= 1 && !bEnd {
				b := dst[0]
				for i := 7; i >= 0; i-- {
					if (b>>i)&1 != 0 {
						pTmpTree = pTmpTree.GetRight()
					} else {
						pTmpTree = pTmpTree.GetLeft()
					}

					if nil != pTmpTree && pTmpTree.Value != nil {
						v, _ := GetHuffmanValue(pTmpTree).(byte)
						src[srcindex] = v
						srcindex++

						if len(dst) == 1 && (uint8(i)) == pHead.patchBit {
							bEnd = true
							break
						}

						pTmpTree = pTree
					}
				}

				dst = dst[1:]
			}
		}

		ioutil.WriteFile(strOutFileName, src, 0666)
	}

	return true
}
