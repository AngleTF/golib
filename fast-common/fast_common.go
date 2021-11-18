package fastCommon

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"golib/fast-check"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func JoinUrl(u1, u2 string) string {
	return strings.TrimRight(u1, "/") + "/" + strings.TrimLeft(u2, "/")
}

func TrimSpaces(args ...*string) {
	for _, v := range args {
		*v = strings.TrimSpace(*v)
	}
}

func Encrypt(sb []byte, slat string) string {
	var (
		c        byte
		i        uint
		size     uint
		slatsb   []byte
		slatsize uint
	)

	size = uint(len(sb))
	slatsb, slatsize = GetByteAndSize(slat)

	if sb == nil || size <= 0 {
		return string(sb)
	}

	for i = 0; i < size; i++ {
		c = (sb[i] << (i % 8)) + (sb[i] >> (8 - i%8))
		sb[i] = c ^ slatsb[i%slatsize]
	}

	return string(sb)
}

func Decrypt(s, slat string) string {
	var (
		c        byte
		i        uint
		sb       []byte
		size     uint
		slatsb   []byte
		slatsize uint
	)
	sb, size = GetByteAndSize(s)
	slatsb, slatsize = GetByteAndSize(slat)

	if sb == nil || size <= 0 {
		return s
	}

	for i = 0; i < size; i++ {
		c = sb[i] ^ slatsb[i%slatsize]
		sb[i] = (c >> (i % 8)) + c<<(8-i%8)
	}
	return string(sb)
}

func GetByteAndSize(s string) ([]byte, uint) {
	var sb = []byte(s)
	var size = uint(len(sb))
	return sb, size
}

func UnZip(zipFileName, deposit string) error {
	var (
		err    error
		rc     *zip.ReadCloser
		bts    []byte
		irc    io.ReadCloser
		file   *os.File
		u8Name []byte
		rd     io.Reader
		gbToU8 *transform.Reader
	)

	if rc, err = zip.OpenReader(zipFileName); err != nil {
		return err
	}

	for _, v := range rc.File {
		if irc, err = v.Open(); err != nil {
			return err
		}

		if bts, err = ioutil.ReadAll(irc); err != nil {
			return err
		}

		rd = bytes.NewReader([]byte(v.Name))

		gbToU8 = transform.NewReader(rd, simplifiedchinese.GBK.NewDecoder())

		if u8Name, err = ioutil.ReadAll(gbToU8); err != nil {
			return err
		}

		v.Name = JoinUrl(deposit, string(u8Name))

		if err = os.MkdirAll(v.Name[:strings.LastIndex(v.Name, "/")], os.ModeDir); err != nil {
			return err
		}

		if file, err = os.Create(v.Name); err != nil {
			return err
		}

		if _, err = file.Write(bts); err != nil {
			return err
		}

		file.Close()
	}

	rc.Close()

	return nil
}

func compress(dirName, deposit string, zWriter *zip.Writer, prefix string) error {
	var (
		fileInfoList []os.FileInfo
		err          error
		fw           io.Writer
		fr           *os.File
		fName        string
		zName        string
	)

	if fileInfoList, err = ioutil.ReadDir(dirName); err != nil {
		return err
	}

	for _, fi := range fileInfoList {
		fName = JoinUrl(dirName, fi.Name())

		if !fastCheck.IsEmpty(prefix) {
			zName = JoinUrl(prefix, fi.Name())
		} else {
			zName = fi.Name()
		}

		if fi.IsDir() {
			if err = compress(fName, deposit, zWriter, zName); err != nil {
				return err
			}
			continue
		}

		if fw, err = zWriter.Create(zName); err != nil {
			return err
		}

		if fr, err = os.Open(fName); err != nil {
			return err
		}

		if _, err = io.Copy(fw, fr); err != nil {
			return err
		}

		fr.Close()
	}

	return nil
}

func StrToZip(fileName, data, deposit string) error {
	var (
		err       error
		zipWriter *zip.Writer
		w         io.Writer
		r         io.Reader
	)

	if zipWriter, err = getZipWrite(deposit); err != nil {
		return err
	}

	if w, err = zipWriter.Create(fileName); err != nil {
		return err
	}

	r = bytes.NewReader([]byte(data))

	if _, err = io.Copy(w, r); err != nil {
		return err
	}

	defer zipWriter.Close()

	return nil
}

func getZipWrite(deposit string) (*zip.Writer, error) {
	var (
		file    *os.File
		err     error
		baseDir string
	)

	baseDir = deposit[:strings.LastIndex(deposit, "/")]

	if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, err
	}

	if file, err = os.OpenFile(deposit, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		return nil, err
	}

	return zip.NewWriter(file), nil
}

func Zip(dirName, deposit string) error {
	var (
		err       error
		zipWriter *zip.Writer
	)

	if zipWriter, err = getZipWrite(deposit); err != nil {
		return err
	}

	if err = compress(dirName, deposit, zipWriter, ""); err != nil {
		return err
	}

	defer zipWriter.Close()
	return nil
}

func Instance2join(structure interface{}, split string, filter []string) string {

	refVal := reflect.ValueOf(structure)
	refTyp := reflect.TypeOf(structure)

	var str []string
	flag := true

	for i := 0; i < refVal.NumField(); i++ {
		val := fmt.Sprint(refVal.Field(i))
		for _, v := range filter {
			if v == refTyp.Field(i).Name {
				flag = false
			}
		}
		if flag {
			str = append(str, val)
		} else {
			flag = true
		}
	}

	joinStr := strings.Join(str, split)

	return joinStr
}

const (
	Y = "2006"
	m = "01"
	d = "02"
	H = "15"
	i = "04"
	s = "05"
)

var template = map[string]string{"Y": Y, "m": m, "d": d, "H": H, "i": i, "s": s}

func Date(mode string, t interface{}) string {

	for k, v := range template {
		mode = strings.Replace(mode, k, v, -1)
	}

	switch t.(type) {
	case time.Time:
		return t.(time.Time).Format(mode)
	case int:
		unix := time.Unix(int64(t.(int)), 0)
		return unix.Format(mode)
	default:
		return strconv.FormatInt(time.Now().Unix(), 10)
	}

}

func Md5(str string, lower bool) string {
	md5Bytes := md5.Sum([]byte(str))
	if lower {
		return fmt.Sprintf("%x", md5Bytes)
	} else {
		return fmt.Sprintf("%X", md5Bytes)
	}
}
