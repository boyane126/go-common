package helpers

import (
	"archive/zip"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// HashString 计算字符串的哈希值（返回一个0-100的数值）
func HashString(input string) int {
	if input == "" {
		return 0
	}

	// 使用MD5计算哈希值
	hash := md5.Sum([]byte(input))

	// 取前4个字节转为uint32
	value := binary.BigEndian.Uint32(hash[:4])

	// 取模，得到0-100的值
	return int(value % 100)
}

// MD5 计算字符串的MD5哈希
func MD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// IsNullOrEmpty 判断字符串是否为空或null
func IsNullOrEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// zipDir 压缩目录
func ZipDir(dirPath, hashCacheFile string, buf io.Writer) error {
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和哈希缓存文件
		if info.IsDir() || filepath.Base(path) == hashCacheFile {
			return nil
		}

		// 创建相对路径
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// 添加到zip
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}
