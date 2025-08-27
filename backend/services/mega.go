package services

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"

	mega "github.com/t3rm1n4l/go-mega"
)

type MegaService struct {
	client *mega.Mega
	fs     *mega.FS
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
	m.client = mega.New()
	if err := m.client.Login(email, pass); err != nil {
		panic("MEGA login failed: " + err.Error())
	}
	m.fs = m.client.FS
}

// CreateFolder under parent node, returns node handle
func (m *MegaService) CreateFolder(parent *mega.Node, name string) (*mega.Node, error) {
	m.mustLogin()
	return m.fs.CreateDir(parent, name)
}

func (m *MegaService) Root() (*mega.Node, error) {
	m.mustLogin()
	return m.fs.GetRoot()
}

func (m *MegaService) FindNodeByHandle(handle string) (*mega.Node, error) {
	m.mustLogin()
	return m.fs.Lookup(handle)
}

func (m *MegaService) UploadBytes(parent *mega.Node, name string, data []byte) (*mega.Node, error) {
	m.mustLogin()
	reader := bytes.NewReader(data)
	return m.client.Upload(parent, reader, name, int64(len(data)))
}

func (m *MegaService) DeleteNode(node *mega.Node) error {
	m.mustLogin()
	return m.fs.Rm(node)
}

func (m *MegaService) Download(node *mega.Node, w io.Writer) error {
	m.mustLogin()
	if node.GetType() != mega.NODE_FILE {
		return errors.New("node is not a file")
	}
	_, err := m.client.Download(node, w)
	return err
}
