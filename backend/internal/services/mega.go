package services

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	mega "github.com/t3rm1n4l/go-mega"
)

type MegaService struct {
	client *mega.Mega
	fs     *mega.MegaFS
	mu     sync.Mutex
}

var megaSvc *MegaService

func Mega() *MegaService {
	if megaSvc == nil {
		megaSvc = &MegaService{}
		megaSvc.mustLogin()
	}
	return megaSvc
}

func (m *MegaService) mustLogin() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		return
	}

	email := os.Getenv("MEGA_EMAIL")
	pass := os.Getenv("MEGA_PASSWORD")
	if email == "" || pass == "" {
		panic("MEGA_EMAIL and MEGA_PASSWORD must be set")
	}

	m.client = mega.New()
	if err := m.client.Login(email, pass); err != nil {
		panic("MEGA login failed: " + err.Error())
	}
	m.fs = m.client.FS
}

// Root returns the filesystem root node.
func (m *MegaService) Root() (*mega.Node, error) {
	m.mustLogin()
	return m.fs.GetRoot(), nil
}

// CreateFolder creates a child directory under parent.
func (m *MegaService) CreateFolder(parent *mega.Node, name string) (*mega.Node, error) {
	m.mustLogin()
	return m.client.CreateDir(name, parent)
}

// FindNodeByHandle resolves a node by its hash/handle (as used by the library).
func (m *MegaService) FindNodeByHandle(handle string) (*mega.Node, error) {
	m.mustLogin()
	n := m.fs.HashLookup(handle)
	if n == nil {
		return nil, errors.New("node not found")
	}
	return n, nil
}

// FindNodeByPath resolves an absolute path like "/photos/2025/trip.jpg".
func (m *MegaService) FindNodeByPath(path string) (*mega.Node, error) {
	m.mustLogin()
	root := m.fs.GetRoot()
	segs := splitPath(path)
	if len(segs) == 0 {
		return root, nil
	}
	nodes, err := m.fs.PathLookup(root, segs)
	if err != nil || len(nodes) == 0 {
		return nil, errors.New("path not found")
	}
	return nodes[len(nodes)-1], nil
}

func splitPath(p string) []string {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return nil
	}
	parts := strings.Split(p, "/")
	out := make([]string, 0, len(parts))
	for _, s := range parts {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// UploadBytes streams a byte slice as a file to MEGA.
func (m *MegaService) UploadBytes(parent *mega.Node, name string, data []byte) (*mega.Node, error) {
	m.mustLogin()

	u, err := m.client.NewUpload(parent, name, int64(len(data)))
	if err != nil {
		return nil, err
	}

	for id := 0; id < u.Chunks(); id++ {
		pos64, sz, err := u.ChunkLocation(id)
		if err != nil {
			return nil, err
		}
		start := int(pos64)
		end := start + sz
		if start < 0 || end > len(data) || start > end {
			return nil, errors.New("invalid chunk bounds")
		}
		if err := u.UploadChunk(id, data[start:end]); err != nil {
			return nil, err
		}
	}

	return u.Finish()
}

// DeleteNode moves a node to trash (destroy=false).
func (m *MegaService) DeleteNode(node *mega.Node) error {
	m.mustLogin()
	return m.client.Delete(node, false)
}

// Download streams a file node into w using chunked download.
func (m *MegaService) Download(node *mega.Node, w io.Writer) error {
	m.mustLogin()
	if node.GetType() != mega.FILE {
		return errors.New("node is not a file")
	}

	dl, err := m.client.NewDownload(node)
	if err != nil {
		return err
	}

	for id := 0; id < dl.Chunks(); id++ {
		chunk, err := dl.DownloadChunk(id)
		if err != nil {
			return err
		}
		if _, err := w.Write(chunk); err != nil {
			return err
		}
	}

	return dl.Finish()
}
