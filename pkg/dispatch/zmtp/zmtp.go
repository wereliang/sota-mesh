// zmtp 2.0 pure golang implement, clone and modify from https://github.com/prepor/go-zmtp
// zmtp 2.0 refer https://rfc.zeromq.org/spec/15/
package zmtp

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type SocketType uint8

type ReadChannel <-chan [][]byte
type AcceptChannel <-chan Accept

const (
	PAIR   SocketType = 0x00
	PUB    SocketType = 0x01
	SUB    SocketType = 0x02
	REQ    SocketType = 0x03
	REP    SocketType = 0x04
	DEALER SocketType = 0x05
	ROUTER SocketType = 0x06
	PULL   SocketType = 0x07
	PUSH   SocketType = 0x08
)

const (
	identitySize = 40
	finalShort   = 0x00
	moreShort    = 0x01
	finalLong    = 0x02
	moreLong     = 0x03
)

type Socket struct {
	t               SocketType
	conn            *net.TCPConn
	identity        []byte
	otherIdentity   []byte
	otherSocketType SocketType
	readChannel     chan [][]byte
	writeChannel    chan [][]byte
	waiter          sync.WaitGroup
	errorMutex      sync.Mutex
	error           error

	writeBuf *bufio.Writer
	readBuf  *bufio.Reader
}

type Listener struct {
	listener      *net.TCPListener
	waiter        sync.WaitGroup
	errorMutex    sync.Mutex
	error         error
	acceptChannel chan Accept
}

type Accept struct {
	Err    error
	Socket *Socket
	Addr   net.Addr
}

type SocketConfig struct {
	Type     SocketType
	Endpoint string
}

func Connect(config *SocketConfig) (*Socket, error) {
	addr, err := net.ResolveTCPAddr("tcp", config.Endpoint)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return socketFromConnection(conn, config)
}

func Listen(config *SocketConfig) (*Listener, error) {
	return newListener(config)
}

func (self *Listener) Accept() AcceptChannel {
	return self.acceptChannel
}

func (self *Socket) Read() ReadChannel {
	return self.readChannel
}

func (self *Socket) Close() (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			return
		}
	}()
	self.conn.Close()
	self.waiter.Wait()
	return nil
}

func (self *Listener) Close() (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			return
		}
	}()
	self.listener.Close()
	self.waiter.Wait()
	return nil
}

func (self *Socket) Error() error {
	self.errorMutex.Lock()
	err := self.error
	self.errorMutex.Unlock()
	return err
}

func (self *Listener) Error() error {
	self.errorMutex.Lock()
	err := self.error
	self.errorMutex.Unlock()
	return err
}

func (self *Socket) Send(frames ...interface{}) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			return
		}
	}()
	var msg [][]byte
	for _, v := range frames {
		switch t := v.(type) {
		case []byte:
			msg = append(msg, t)
		case [][]byte:
			msg = append(msg, t...)
		case string:
			msg = append(msg, []byte(t))
		case []string:
			strings := make([][]byte, len(t))
			for i, s := range t {
				strings[i] = []byte(s)
			}
			msg = append(msg, strings...)
		default:
			msg = append(msg, []byte(fmt.Sprintf("%v", t)))
		}
	}
	self.writeChannel <- msg
	return nil
}

func (self *Socket) prepare() error {
	if err := self.sendGreeting(); err != nil {
		return fmt.Errorf("zmqgo/zmtp: Got error while sending greeting: %v", err)
	}
	if err := self.recvGreeting(); err != nil {
		return fmt.Errorf("zmqgo/zmtp: Got error while receiving greeting: %v", err)
	}
	return nil
}

const (
	signaturePrefix = 0xFF
	signatureSuffix = 0x7F
	revision        = 0x01
)

var byteOrder = binary.BigEndian

type greeting struct {
	SignaturePrefix byte
	_               [8]byte
	SignatureSuffix byte
	Revision        byte
	SocketType      byte
	FinalShort      byte
	IdentitySize    byte
}

func (self *Socket) sendGreeting() error {
	var identity [identitySize]byte
	rand.Read(identity[:])
	greeting := greeting{
		SignaturePrefix: signaturePrefix,
		SignatureSuffix: signatureSuffix,
		Revision:        revision,
		SocketType:      byte(self.t),
		FinalShort:      finalShort,
		IdentitySize:    identitySize,
	}
	if err := binary.Write(self.writeBuf, byteOrder, &greeting); err != nil {
		return err
	}
	if _, err := self.writeBuf.Write(identity[:]); err != nil {
		return err
	}
	self.identity = identity[:]
	self.writeBuf.Flush()
	return nil
}

func (self *Socket) recvGreeting() error {
	var greeting greeting

	if err := binary.Read(self.readBuf, byteOrder, &greeting); err != nil {
		return fmt.Errorf("Error while reading: %v", err)
	}

	if greeting.SignaturePrefix != signaturePrefix {
		return fmt.Errorf("Signature prefix received does not correspond with expected signature. Received: %#v. Expected: %#v.", greeting.SignaturePrefix, signaturePrefix)
	}

	if greeting.SignatureSuffix != signatureSuffix {
		return fmt.Errorf("Signature prefix received does not correspond with expected signature. Received: %#v. Expected: %#v.", greeting.SignatureSuffix, signatureSuffix)
	}

	otherIdentity := make([]byte, greeting.IdentitySize)
	if _, err := io.ReadFull(self.readBuf, otherIdentity); err != nil {
		return err
	}
	self.otherIdentity = otherIdentity
	self.otherSocketType = SocketType(greeting.SocketType)

	if !IsSocketTypesCompatible(self.otherSocketType, self.t) {
		return fmt.Errorf("Socket type %v is not compatible with %v", self.otherSocketType, self.t)
	}
	return nil
}

const readChannelBuffer = 10
const writeChannelBuffer = 10

func (self *Socket) startChannels(config *SocketConfig) {
	self.writeChannel = make(chan [][]byte, writeChannelBuffer)
	self.readChannel = make(chan [][]byte, readChannelBuffer)

	go self.startReadChannel()
	go self.startWriteChannel()
}

func (self *Socket) startWriteChannel() {
	self.waiter.Add(1)
	defer self.waiter.Done()
	var err error
	for v := range self.writeChannel {
		if len(v) != 0 {
			err = self.send(v)
			if err != nil {
				self.error = err
				return
			}
		}
	}
}

func (self *Socket) startReadChannel() {
	self.waiter.Add(1)
	defer close(self.readChannel)
	defer close(self.writeChannel)
	defer self.waiter.Done()
	for {
		msg := make([][]byte, 0, 3)
		for {
			frame, isLast, err := self.readFrame()
			if err != nil {
				self.error = err
				return
			}
			if isLast {
				msg = append(msg, frame)
				self.readChannel <- msg
				break
			} else {
				msg = append(msg, frame)
			}
		}

	}
}

func (self *Socket) send(msg [][]byte) error {
	var err error
	var last bool
	len := len(msg)
	for i, v := range msg {
		last = i == len-1
		if err = self.sendFrame(v, last); err != nil {
			return err
		}
	}
	self.writeBuf.Flush()
	return nil
}

func (self *Socket) sendFrame(body []byte, isLast bool) error {
	length := len(body)
	isLong := length > 255
	var mark uint8
	if isLast {
		if isLong {
			mark = finalLong
		} else {
			mark = finalShort
		}
	} else {
		if isLong {
			mark = moreLong
		} else {
			mark = moreShort
		}
	}

	if err := binary.Write(self.writeBuf, byteOrder, mark); err != nil {
		return err
	}
	if isLong {
		if err := binary.Write(self.writeBuf, byteOrder, int64(length)); err != nil {
			return err
		}
	} else {
		if err := binary.Write(self.writeBuf, byteOrder, uint8(length)); err != nil {
			return err
		}
	}
	if _, err := self.writeBuf.Write(body); err != nil {
		return err
	}
	return nil
}

func (self *Socket) readFrame() ([]byte, bool, error) {
	var header [2]byte
	var longLength [4]byte
	isLast := true
	if _, err := io.ReadFull(self.readBuf, header[:]); err != nil {
		return nil, true, err
	}
	if header[0] == moreLong || header[0] == moreShort {
		isLast = false
	}
	bodyLength := uint64(0)
	if header[0] == finalLong || header[0] == moreLong {
		longLength[0] = header[1]
		if _, err := io.ReadFull(self.readBuf, longLength[1:]); err != nil {
			return nil, true, err
		}
	} else {
		bodyLength = uint64(header[1])
	}
	body := make([]byte, bodyLength)
	if _, err := io.ReadFull(self.readBuf, body); err != nil {
		return nil, true, err
	}
	return body, isLast, nil
}

func (self *Socket) setError(err error) {
	self.errorMutex.Lock()
	self.error = err
	self.errorMutex.Unlock()
}

func (self *Listener) setError(err error) {
	self.errorMutex.Lock()
	self.error = err
	self.errorMutex.Unlock()
}

func newListener(config *SocketConfig) (*Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp", config.Endpoint)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	v := &Listener{
		listener:      listener,
		acceptChannel: make(chan Accept),
	}
	go v.start(config)
	return v, nil
}

func (self *Listener) start(config *SocketConfig) {
	defer close(self.acceptChannel)
	defer self.waiter.Done()
	self.waiter.Add(1)
	for {
		conn, err := self.listener.AcceptTCP()
		if err != nil {
			self.setError(err)
			return
		}
		socket, err := socketFromConnection(conn, config)
		if err != nil {
			self.acceptChannel <- Accept{Err: err}
		} else {
			self.acceptChannel <- Accept{Socket: socket, Addr: conn.RemoteAddr()}
		}
	}
}

func socketFromConnection(conn *net.TCPConn, config *SocketConfig) (*Socket, error) {
	s := &Socket{
		t:        config.Type,
		conn:     conn,
		writeBuf: bufio.NewWriter(conn),
		readBuf:  bufio.NewReader(conn),
	}
	err := s.prepare()
	if err != nil {
		return nil, err
	}
	s.startChannels(config)
	return s, nil
}

func socketFromReadWriter(rw io.ReadWriter, config *SocketConfig) (*Socket, error) {
	s := &Socket{
		t: config.Type,
		// conn:     conn,
		writeBuf: bufio.NewWriter(rw),
		readBuf:  bufio.NewReader(rw),
	}
	return s, nil
}

func (self SocketType) String() string {
	switch self {
	default:
		return "UNKNOWN"
	case PAIR:
		return "PAIR"
	case PUB:
		return "PUB"
	case SUB:
		return "SUB"
	case REQ:
		return "REQ"
	case REP:
		return "REP"
	case DEALER:
		return "DEALER"
	case ROUTER:
		return "ROUTER"
	case PULL:
		return "PULL"
	case PUSH:
		return "PUSH"
	}
}

func IsSocketTypesCompatible(one SocketType, another SocketType) bool {
	switch one {
	default:
		return false
	case PAIR:
		return another == PAIR
	case PUB:
		return another == SUB
	case SUB:
		return another == PUB
	case REQ:
		return another == REP || another == ROUTER
	case REP:
		return another == REQ || another == DEALER
	case DEALER:
		return another == REP || another == DEALER || another == ROUTER
	case ROUTER:
		return another == REQ || another == DEALER || another == ROUTER
	case PULL:
		return another == PUSH
	case PUSH:
		return another == PULL
	}
}
