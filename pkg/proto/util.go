package proto

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	S3Prefix    = "s3://"
	GSPrefix    = "gs://"
	LocalPrefix = "local://"
)

type BackendType int

const (
	LocalType BackendType = iota
	S3Type
	GSType
)

func (t BackendType) String() string {
	switch t {
	case LocalType:
		return "local"
	case S3Type:
		return "s3"
	case GSType:
		return "gs"
	default:
		return "unknown"
	}
}

func ParseType(uri string) BackendType {
	if strings.HasPrefix(uri, LocalPrefix) {
		return LocalType
	} else if strings.HasPrefix(uri, S3Prefix) {
		return S3Type
	} else if strings.HasPrefix(uri, GSPrefix) {
		return GSType
	}

	return -1
}

func (b *Backend) Type() BackendType {
	t := ParseType(b.Uri())
	return t
}

func (b *Backend) Uri() string {
	if s3 := b.GetS3(); s3 != nil {
		return S3Prefix + s3.GetBucket() + "/" + s3.GetPath()
	} else if gs := b.GetGs(); gs != nil {
		return GSPrefix + gs.GetBucket() + "/" + gs.GetPath()
	} else if local := b.GetLocal(); local != nil {
		return LocalPrefix + local.GetPath()
	}

	return "nil path"
}

func (b *Backend) SetUri(uri string) error {
	t := ParseType(uri)
	switch t {
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
		u, err := parseUri(uri)
		if err != nil {
			return err
		}

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

	case GSType:
		u, err := parseUri(uri)
		if err != nil {
			return err
		}

		if b.GetGs() == nil {
			b.Storage = &Backend_Gs{
				Gs: &GS{
					Bucket: u.Host,
					Path:   u.Path,
				},
			}
		} else {
			b.GetGs().Bucket = u.Host
			b.GetGs().Path = u.Path
		}
	default:
		return fmt.Errorf("unknow storage backend type")
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
	case GSType:
		gs := *b.GetGs()
		cp.Storage = &Backend_Gs{&gs}
	}
	return cp
}

func parseUri(uri string) (*url.URL, error) {
	t := ParseType(uri)
	u, err := url.Parse(uri)
	if err != nil || strings.ToLower(u.Scheme) != t.String() {
		return nil, fmt.Errorf("invalid uri: %s", uri)
	}
	u.Host = strings.TrimRight(u.Host, "/ ")
	u.Path = strings.TrimLeft(u.Path, "/ ")
	return u, nil
}
