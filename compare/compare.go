package compare

import (
    "os"
    "crypto/md5"
)

type Comparator struct {
    Dir1 string
    Dir2 string

    FilesInDir1 map[string]os.FileInfo
    FilesInDir2 map[string]os.FileInfo

    UniqueFilesInDir1 map[string]os.FileInfo
    CommonFiles map[string]os.FileInfo
    UniqueFilesInDir2 map[string]os.FileInfo

    CommonFilesSize int64
}

func Init(dir1 string, dir2 string) *Comparator, error{
    _, err = os.Stat(dir1)
    if err != nil {
        logrus.Fatalf("Dir: %s does not exist...\n", dir1)
        return nil, err
    }

    _, err = os.Stat(dir2)
    if err != nil {
        logrus.Fatalf("Dir: %s does not exist...\n", dir1)
        return nil, err
    }

    return &Comparator {
        Dir1: dir1, 
        Dir2: dir2,
    }, nil
}

func (c *Comparator)Compare() error {
    c.FilesInDir1, err := walk(c.Dir1)
    if err != nil {
        logrus.Fatalf("Fail to walk dir: %s...", c.Dir1)
    }
    c.FilesInDir2, err := walk(c.Dir2)
    if err != nil {
        logrus.Fatalf("Fail to walk dir: %s...", c.Dir2)
    }

    c.UniqueFilesInDir1 = c.UniqueFilesInDir1
    c.UniqueFilesInDir2 = c.UniqueFilesInDir2

    int64 commonSize = 0

    for hash, fileInfo := range Comparator.FilesInDir1 {
        _, ok := Comparator.FilesInDir2[hash]
        if ok {
            c.CommonFiles[hash] = fileInfo
            commonSize += fileInfo.Size()
            delete(c.UniqueFilesInDir1, hash)
            delete(c.UniqueFilesInDir2, hash)
        }
    }

    return nil
}

func walk(dir string) map[string]os.FileInfo, error {
    var commonFIles map[string]os.FileInfo

    err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
            // fail to get file info
            if f == nil {return err}

            // current file is a regular file
            if f.Mode().IsRegular() {
                src, err := os.Open(path)
                if err != nil {
                    logrus.Warnf("Fail to open file: %s\n", path)
                    return err
                }
                defer src.Close()

                m := md5.New()
                _, err := io.Copy(m, src)
                if err != nil {
                    logrus.Warn("Fail to copy from src to md5 instance...")
                }
                hashValue := fmt.Sprintf("%x", m.Sum(nil))
                commonFIles[hashValue] = f
            }

            return nil
        })

    if err != nil {
        logrus.Warnf("Fail to walk dir: %s...", dir)
        return err
    }

    return nil
}


















