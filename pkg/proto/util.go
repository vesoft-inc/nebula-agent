package proto

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	S3Prefix    = "s3://"
	LocalPrefix = "local://"
	HdfsPrefix  = "hdfs://"
)

type BackendType int

const (
	InvalidType BackendType = iota
	LocalType
	HdfsType
	S3Type
)

func (t BackendType) String() string {
	switch t {
	case LocalType:
		return "Local"
	case HdfsType:
		return "HDFS"
	case S3Type:
		return "S3-compatible"
	case InvalidType:
		return "Invalid"
	default:
		return "Unknown"
	}
}

func ParseType(uri string) BackendType {
	if strings.HasPrefix(uri, LocalPrefix) {
		return LocalType
	} else if strings.HasPrefix(uri, S3Prefix) {
		return S3Type
	}

	return InvalidType
}

func (b *Backend) Type() BackendType {
	t := ParseType(b.Uri())
	switch t {
	case LocalType:
		if b.GetLocal() == nil {
			panic(fmt.Errorf("uri %s is local, but local pointer is nil", b.Uri()))
		}
	case S3Type:
		if b.GetS3() == nil {
			panic(fmt.Errorf("uri %s is s3, but s3 pointer is nil", b.Uri()))
		}
	}
	return t
}

func (b *Backend) Uri() string {
	if hdfs := b.GetHdfs(); hdfs != nil {
		return hdfs.GetRemote()
	} else if s3 := b.GetS3(); s3 != nil {
		return S3Prefix + s3.GetBucket() + "/" + s3.GetPath()
	} else if local := b.GetLocal(); local != nil {
		return LocalPrefix + local.GetPath()
	}

	return "nil path"
}

func (b *Backend) SetUri(uri string) error {
	if ParseType(uri) == InvalidType {
		return fmt.Errorf("invalid uri: %s", uri)
	}
	if b.Type() != InvalidType && b.Type() != ParseType(uri) {
		return fmt.Errorf("not same type uri, backend: %s, paramter: %s", b.Uri(), uri)
	}

	switch ParseType(uri) {
	case LocalType:
		if b.GetLocal() == nil {
			b.Storage = &Backend_Local{
				Local: &Local{
					Path: strings.TrimPrefix(uri, LocalPrefix),
				},
			}
		} else {
			b.GetLocal().Path = strings.TrimPrefix(uri, LocalPrefix)
		}

	case S3Type:
		u, err := url.Parse(uri)
		if err != nil || strings.ToLower(u.Scheme) != "s3" {
			return fmt.Errorf("bad s3 format uri: %s", uri)
		}
		u.Host = strings.TrimRight(u.Host, "/ ")
		u.Path = strings.TrimLeft(u.Path, "/ ")

		if b.GetS3() == nil {
			b.Storage = &Backend_S3{
				S3: &S3{
					Bucket: u.Host,
					Path:   u.Path,
				},
			}
		} else {
			b.GetS3().Bucket = u.Host
			b.GetS3().Path = u.Path
		}
	case HdfsType:
		if b.GetHdfs() == nil {
			b.Storage = &Backend_Hdfs{
				Hdfs: &HDFS{
					Remote: uri,
				},
			}
		} else {
			b.GetHdfs().Remote = uri
		}
	}

	return nil
}

func (b *Backend) DeepCopy() *Backend {
	cp := &Backend{}
	switch b.Type() {
	case LocalType:
		l := *b.GetLocal()
		cp.Storage = &Backend_Local{&l}
	case S3Type:
		s3 := *b.GetS3()
		cp.Storage = &Backend_S3{&s3}
	}
	return cp
}
