package util

import (
	"bufio"
	"io"
	"k8s.io/klog"
	"os"
)

func ReadFile(filePath string) (result []string, err error) {
	file, err := os.Open("log/1.67_apiserver.log_3") ; if err != nil {
		klog.Warningf("Open file: %+v, err: %+v\n", filePath, err)
		return nil, err
	}
	defer func() {
		if err := file.Close() ; err != nil {
			klog.Warningf("Close file: %+v err: %+v", filePath, err)
		}
	}()

	result = make([]string, 0)
	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		result = append(result, string(a))
	}
	return result, nil
}

func ReadFileWithFilter(filePath string, filter func(string) bool, extractor func(string) string) (result []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		klog.Warningf("Open file: %+v, err: %+v\n", filePath, err)
		return nil, err
	}
	defer func() {
		if err := file.Close() ; err != nil {
			klog.Warningf("Close file: %+v err: %+v", filePath, err)
		}
	}()

	result = make([]string, 0)
	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if filter(string(a)) {
			result = append(result, extractor(string(a)))
		}
	}
	return result, nil
}

func ReadWriteFileWithFilter(sourcePath string, targetPath string, filter func(string) bool) (count int, err error) {
	count = 0
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		klog.Warningf("Open file: %+v, err: %+v\n", sourcePath, err)
		return count, err
	}

	_ = os.Remove(targetPath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		klog.Warningf("Open file: %+v, err: %+v\n", targetPath, err)
		return count, err
	}

	defer func() {
		if err := sourceFile.Close() ; err != nil {
			klog.Warningf("Close file: %+v err: %+v", sourcePath, err)
			return
		}
	}()
	defer func() {
		if err := targetFile.Close() ; err != nil {
			klog.Warningf("Close file: %+v err: %+v", targetPath, err)
			return
		}
	}()

	br := bufio.NewReader(sourceFile)
	bw := bufio.NewWriter(targetFile)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if filter(string(a)) {
			_, err = bw.WriteString(string(a)+"\n")
			count++
		}
	}
	return count, err
}