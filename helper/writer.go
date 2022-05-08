package helper

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/lwinmgmg/http_data_store/environ"
)

var env *environ.Environ = environ.NewEnviron()

func selectFile(preFixs []string, name string) (file *os.File, fileName string, err error) {
	for _, v := range preFixs {
		fileName = path.Join(v, name)
		file, err = os.Create(fileName)
		if err != nil {
			continue
		}
		return
	}
	return
}

func copyFiles(originalPath, fileName string, othersPaths []string) {
	reader, err := os.Open(originalPath)
	if err != nil {
		return
	}
	for _, v := range othersPaths {
		reader.Seek(0, 0)
		fPath := path.Join(v, fileName)
		if fPath == originalPath {
			continue
		}
		writer, err := os.Create(fPath)
		if err != nil {
			continue
		}
		if _, err := io.Copy(writer, reader); err != nil {
			continue
		}
	}
}

func ArchiveWriter(bkFileName, fileName string, writedPathList [2]string) {
	wr, err := os.Create(bkFileName)
	if err != nil {
		return
	}
	defer wr.Close()
	gz := gzip.NewWriter(wr)
	defer gz.Close()
	af := tar.NewWriter(gz)
	defer af.Close()
	hdr := &tar.Header{}
	for k, writedPath := range writedPathList {
		fRdr, err := os.Open(writedPath)
		if err != nil {
			return
		}
		fStat, err := fRdr.Stat()
		if err != nil {
			return
		}
		hdr.Name = fmt.Sprintf("%v.%v", k, fileName)
		hdr.Mode = 0600
		hdr.Size = fStat.Size()
		af.WriteHeader(hdr)
		io.Copy(af, fRdr)
		if err = fRdr.Close(); err != nil {
			return
		}
	}
}

type WriterManager struct {
	Reader    io.Reader
	TotalSize int64
	FileName  string
	FirstDirs []string
	LastDirs  []string
	BackupDir string
}

func (wMgr *WriterManager) WriteOriginal() (int64, int64, error) {
	var firstSize, lastSize int64
	chunkSize := 10 * 1024 * 1024
	halfSize := wMgr.TotalSize / 2
	var buf []byte
	if halfSize < int64(chunkSize) {
		buf = make([]byte, halfSize/2)
	} else {
		buf = make([]byte, chunkSize)
	}
	firstFile, firstFileName, err := selectFile(wMgr.FirstDirs, wMgr.FileName)
	if err != nil {
		return firstSize, lastSize, err
	}
	defer firstFile.Close()
	lastFile, lastFileName, err := selectFile(wMgr.LastDirs, wMgr.FileName)
	if err != nil {
		return firstSize, lastSize, err
	}
	defer lastFile.Close()
	var writer []*os.File = []*os.File{
		firstFile,
		lastFile,
	}
	var readSize int64
	var tmpReadSize int
	var control int = 0
	for {
		tmpReadSize, err = wMgr.Reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return firstSize, lastSize, err
		}
		_, err := writer[control].Write(buf[:tmpReadSize])
		if err != nil {
			return firstSize, lastSize, err
		}
		readSize += int64(tmpReadSize)
		if readSize > halfSize && control == 0 {
			firstSize = readSize
			readSize = 0
			control = 1
		}
	}
	lastSize = readSize
	go ArchiveWriter(path.Join(wMgr.BackupDir, wMgr.FileName), wMgr.FileName, [2]string{firstFileName, lastFileName})
	go copyFiles(firstFileName, wMgr.FileName, wMgr.FirstDirs)
	go copyFiles(lastFileName, wMgr.FileName, wMgr.LastDirs)
	return firstSize, lastSize, nil
}

func NewWriterManager(reader io.Reader, totalSize int64) WriterManager {
	return WriterManager{
		Reader:    reader,
		TotalSize: totalSize,
		FileName:  fmt.Sprintf("%v-%v", GetUniqueKey(), time.Now().Unix()),
		FirstDirs: env.HDS_FIRST_DIRS,
		LastDirs:  env.HDS_LAST_DIRS,
		BackupDir: env.HDS_BACKUP_DIR,
	}
}

func CheckDirectoryExistOrNot(filePath string) error {
	fInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if !fInfo.IsDir() {
		return fmt.Errorf("%v file does not exist", filePath)
	}
	return nil
}

func DeleteFiles(fileName string, firstDirs []string, lastDirs []string, bkDir string) {
	prefixPaths := make([]string, 0, 5)
	prefixPaths = append(prefixPaths, firstDirs...)
	prefixPaths = append(prefixPaths, lastDirs...)
	prefixPaths = append(prefixPaths, bkDir)
	for _, k := range prefixPaths {
		pathTemp := path.Join(k, fileName)
		if err := os.Remove(pathTemp); err != nil {
			fmt.Printf("[%v] path does not exist", pathTemp)
		}
	}
}

func GetFile(filename string, dirList []string, size int64) (*os.File, func() error, error) {
	idx := rand.Intn(len(dirList))
	filePath := path.Join(dirList[idx], filename)
	fInfo, err := os.Stat(filePath)
	if err != nil || fInfo.Size() != size {
		if idx == 0 {
			filePath = path.Join(dirList[1], filename)
		} else {
			filePath = path.Join(dirList[0], filename)
		}
	}
	file, err := os.Open(filePath)
	return file, file.Close, err
}

func MultiCloser(closers ...func() error) {
	for i := 0; i < len(closers); i++ {
		closers[i]()
	}
}

func ReadFile(fileName string, size1, size2 int64) (io.Reader, []func() error, error) {
	closers := make([]func() error, 2)
	firstFile, closer, err := GetFile(fileName, env.HDS_FIRST_DIRS, size1)
	if err != nil {
		return nil, closers, err
	}
	closers[0] = closer
	lastFile, closer, err := GetFile(fileName, env.HDS_LAST_DIRS, size2)
	if err != nil {
		return nil, closers, err
	}
	closers[1] = closer
	return io.MultiReader(firstFile, lastFile), closers, nil
}
