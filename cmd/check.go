package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"

	"github.com/minio/minio-go"
)

type CheckResult struct {
	Offset int64
	SrcMD5 string
	minio.ObjectPart
}

func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getObjectSize(objectParts map[int]minio.ObjectPart) (size int64) {
	size = 0
	index := 1
	for true {
		value, ok := objectParts[index]
		if ok == false {
			break
		}
		size += value.Size
		index++
	}
	return size
}

func doCheck(objectParts map[int]minio.ObjectPart, path string) (checkResultList []*CheckResult, err error) {
	if pathExist(path) == false {
		err = fmt.Errorf("Path is not exist")
		return checkResultList, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var offset int64 = 0
	index := 1
	for true {
		value, ok := objectParts[index]
		if ok == false {
			break
		}
		b := make([]byte, value.Size)
		n, _ := f.ReadAt(b, offset)

		h := md5.New()
		io.WriteString(h, string(b[:n]))
		md5Str := hex.EncodeToString(h.Sum(nil))
		if md5Str != value.ETag {
			var m *CheckResult = &CheckResult{offset, md5Str, value}
			checkResultList = append(checkResultList, m)
			color.Red("Wrong: localMd5[%s] remoteMd5[%s] partNumber[%d] \n", md5Str, value.ETag, value.PartNumber)
		} else {
			color.Green("Right: localMd5[%s] remoteMd5[%s] partNumber[%d] \n", md5Str, value.ETag, value.PartNumber)
		}

		offset += value.Size
		index++
	}
	if offset < fileInfo.Size() {
		color.Red("Wrong: Local file Size is bigger than file in s3 \n")
	}

	return checkResultList, nil
}

func doRepair(urlStr string, checkResultList []*CheckResult, path string, objectSize int64) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer f.Close()
	color.Blue("\nStart repail size\n")
	err = f.Truncate(objectSize)
	if err != nil {
		return err
	}

	if len(checkResultList) <= 0 {
		return nil
	}

	color.Blue("\nStart repail part\n")
	client, perr := newClient(urlStr)
	if perr != nil {
		err = fmt.Errorf("New client faild")
		return err
	}

	var repairFaild int = 0

	for _, value := range checkResultList {
		b := make([]byte, value.Size)
		err = client.ReadAt(b, value.Offset)
		if err != nil {
			repairFaild++
			color.Red("Repair faild: partNumber[%d], can not get data for s3, err: %s\n", value.PartNumber, err.Error())
			continue
		}
		h := md5.New()
		io.WriteString(h, string(b[:value.Size]))
		md5Str := hex.EncodeToString(h.Sum(nil))
		if md5Str != value.ETag {
			color.Red("Wrong: localMd5[%s] remoteMd5[%s] partNumber[%d] \n", md5Str, value.ETag, value.PartNumber)
			repairFaild++
			continue
		} else {
			color.Green("Right: localMd5[%s] remoteMd5[%s] partNumber[%d] \n", md5Str, value.ETag, value.PartNumber)
		}

		_, err := f.WriteAt(b, value.Offset)
		if err != nil {
			repairFaild++
			color.Red("Repair faild: partNumber[%d], write to file faild, err: %s \n", value.PartNumber, err.Error())
			continue
		}
		color.Green("Repair succeed: partNumber[%d] \n", value.PartNumber)
	}

	if repairFaild > 0 {
		err = fmt.Errorf("%d part Repail faild", repairFaild)
		return err
	}
	return nil
}
