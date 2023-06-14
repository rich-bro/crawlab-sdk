package entity

import "io"

type OssTask struct {
	Type         int       `json:"type"` //1 by file,2 by io.read
	FilePath     string    `json:"file_path"`
	OssPath      string    `json:"oss_path"`
	FileIOReader io.Reader `json:"file_io_reader"`
}
