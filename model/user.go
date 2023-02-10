package model

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"errors"
	"net"
	"sort"
	"time"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"

	"wgnetwork/pkg/otp"
	"wgnetwork/pkg/rand"
)

type userSession struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// User model.
type User struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`

	IsManager bool        `json:"is_manager"`
	TFASecret string      `json:"tfa_secret"`
	Session   userSession `json:"session"`

	Devices []net.IP `json:"devices"`
}

// NewUser constructor
func NewUser(name string) User {
	user := User{
		UUID: uuid.New().String(),
		Name: name,

		Devices: []net.IP{},
	}

	return user
}

// LoadUser constructor
func LoadUser(tx *bolt.Tx, uuid string) (User, error) {
	bname := []byte("users")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return User{}, errors.New("not found")
	}

	key := []byte(uuid)
	v := bucket.Get(key)
	if v == nil {
		return User{}, errors.New("not found")
	}

	u := User{}
	err := json.Unmarshal(v, &u)
	if err != nil {
		return User{}, err
	}

	for i := range u.Devices {
		u.Devices[i] = u.Devices[i].To4()
	}

	return u, nil
}

// SessionUser constructor
func SessionUser(tx *bolt.Tx, secret, s string) (User, error) {
	envelope, err := ParseEnvelope(s)
	if err != nil {
		return User{}, err
	}

	if !envelope.CheckSignature(secret) {
		return User{}, errors.New("broken session")
	}

	u, err := LoadUser(tx, envelope.UUID())
	if err != nil {
		return User{}, errors.New("broken session")
	}

	if !u.IsManager {
		return User{}, errors.New("user isn't manager")
	}

	if u.Session.Token != envelope.Payload() {
		return User{}, errors.New("broken session")
	}

	return u, nil
}

// SetManager privileges.
func (u *User) SetManager(tfaIssuer string) (string, string, error) {
	b, err := rand.CRandomBytes(20)
	if err != nil {
		return "", "", err
	}
	period := uint16(30) // 30 seconds
	window := int8(1)    // 1 * 30 = 30 seconds
	otp := otp.New(b, period, window)

	u.IsManager = true
	u.TFASecret = base32.StdEncoding.EncodeToString(b)

	uri, secret := otp.ProvisionURI(u.Name, tfaIssuer)

	return uri, secret, nil
}

// UnsetManager privileges.
func (u *User) UnsetManager() {
	u.IsManager = false
	u.TFASecret = ""
}

// OTPCheck code.
func (u *User) OTPCheck(code string) error {
	period := uint16(30) // 30 seconds
	window := int8(1)    // 1 * 30 = 30 seconds
	secret, err := base32.StdEncoding.DecodeString(u.TFASecret)
	if err != nil {
		return err
	}

	otp := otp.New(secret, period, window)
	ok, err := otp.Check(code)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("bad code")
	}

	return nil
}

// OTPProvisionURI for authenticator app.
func (u *User) OTPProvisionURI(tfaIssuer string) (string, error) {
	period := uint16(30) // 30 seconds
	window := int8(1)    // 1 * 30 = 30 seconds
	secret, err := base32.StdEncoding.DecodeString(u.TFASecret)
	if err != nil {
		return "", err
	}

	otp := otp.New(secret, period, window)
	uri, _ := otp.ProvisionURI(u.Name, tfaIssuer)

	return uri, nil
}

// AddDevice ip.
func (u *User) AddDevice(ip net.IP) {
	ip = ip.To4()
	if ip == nil {
		return
	}

	_, found := u.isDeviceExists(ip)
	if found {
		return
	}

	u.Devices = append(u.Devices, ip)
	u.sortDevices()
}

// RemoveDevice ip.
func (u *User) RemoveDevice(ip net.IP) {
	ip = ip.To4()
	if ip == nil {
		return
	}

	i, found := u.isDeviceExists(ip)
	if !found {
		return
	}

	u.Devices[i] = u.Devices[len(u.Devices)-1]
	u.Devices[len(u.Devices)-1] = nil
	u.Devices = u.Devices[:len(u.Devices)-1]
	u.sortDevices()
}

// CreateSession for the user
func (u *User) CreateSession(
	secret string,
	ttl time.Duration,
) (string, error) {
	if !u.IsManager {
		err := errors.New("user isn't manager")
		return "", err
	}

	token := rand.RandomString(20)
	expires := time.Now().Add(ttl)

	envelope, err := NewEnvelope(
		secret,
		u.UUID,
		token,
		expires)
	if err != nil {
		return "", err
	}

	u.Session = userSession{Token: token, Expires: expires.Unix()}

	return envelope.String(), nil
}

// ProlongSession for the user
func (u *User) ProlongSession(
	secret string,
	ttl time.Duration,
) (string, error) {
	if !u.IsManager {
		err := errors.New("user isn't manager")
		return "", err
	}

	if u.Session == (userSession{}) {
		err := errors.New("session doesn't exists")
		return "", err
	}

	token := u.Session.Token
	if len(token) != 20 {
		return "", errors.New("broken session")
	}

	expires := time.Unix(u.Session.Expires, 0)
	if expires.Before(time.Now().UTC()) {
		return "", errors.New("expired session")
	}
	expires = time.Now().Add(ttl)

	envelope, err := NewEnvelope(
		secret,
		u.UUID,
		token,
		expires)
	if err != nil {
		return "", err
	}

	u.Session = userSession{Token: token, Expires: expires.Unix()}

	return envelope.String(), nil
}

// DestroySession for the user
func (u *User) DestroySession() {
	u.Session = userSession{}
}

// Store to database.
func (u *User) Store(tx *bolt.Tx) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	bname := []byte("users")
	bucket, err := tx.CreateBucketIfNotExists(bname)
	if err != nil {
		return err
	}

	key := []byte(u.UUID)
	value, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return bucket.Put(key, value)
}

func (u *User) sortDevices() {
	sort.Slice(u.Devices, func(i, j int) bool {
		return bytes.Compare(u.Devices[i], u.Devices[j]) < 0
	})
}

func (u *User) isDeviceExists(ip net.IP) (int, bool) {
	i, found := sort.Find(len(u.Devices), func(i int) int {
		return bytes.Compare(ip, u.Devices[i])
	})

	return i, found
}

// RemoveUser from database
func RemoveUser(tx *bolt.Tx, uuid string) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	u, err := LoadUser(tx, uuid)
	if err != nil {
		return err
	}

	for i := range u.Devices {
		u.Devices[i] = u.Devices[i].To4()
		err := RemoveDevice(tx, u.Devices[i])
		if err != nil {
			return err
		}
	}

	bname := []byte("users")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return nil
	}

	key := []byte(u.UUID)
	return bucket.Delete(key)
}

// Users type
type Users []User

// LoadUsers returns all users from database.
func LoadUsers(tx *bolt.Tx) (Users, error) {
	bname := []byte("users")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return nil, nil
	}

	cnt := bucket.Stats().KeyN
	users := make(Users, cnt)

	i := 0
	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		u := User{}
		err := json.Unmarshal(v, &u)
		if err != nil {
			return nil, err
		}

		for i := range u.Devices {
			u.Devices[i] = u.Devices[i].To4()
		}

		users[i] = u
		i++
	}

	return users, nil
}
