package sand

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	ClientId = ""
	ApiKey   = ""
)

type Droplet struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	ImageId          int       `json:"image_id"`
	SizeId           int       `json:"size_id"`
	RegionId         int       `json:"region_id"`
	BackupsActive    bool      `json:"backups_active"`
	IPAddress        string    `json:"ip_address"`
	PrivateIPAddress string    `json:"private_ip_address"`
	Locked           bool      `json:"locked"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type DropletCreation struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	ImageId int    `json:"image_id"`
	SizeId  int    `json:"size_id"`
	EventId int    `json:"event_id"`
}

type Domain struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	TTL               int    `json:"ttl"`
	LiveZoneFile      string `json:"live_zone_file"`
	Error             string `json:"error"`
	ZoneFileWithError string `json:"zone_file_with_error"`
}

type Record struct {
	Id         int    `json:"id"`
	DomainId   int    `json:"domain_id"`
	RecordType string `json:"record_type"`
	Name       string `json:"name"`
	Data       string `json:"data"`
	Priority   int    `json:"priority"`
	Port       string `json:"port"`
	Weight     string `json:"weight"`
}

type Key struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"ssh_pub_key"`
}

type Image struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Distribution string `json:"distribution"`
}

type Region struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Size struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type EventId int

type Event struct {
	Id         int    `json:"id"`
	Status     string `json:"action_status"`
	DropletId  int    `json:"droplet_id"`
	EventType  int    `json:"event_type_id"`
	Percentage string `json:"percentage"`
}

type Response interface {
	GetStatus() string
	GetMessage() string
}

type DropletsResponse struct {
	Status   string     `json:"status"`
	Message  string     `json:"message"`
	Droplets []*Droplet `json:"droplets"`
}

func (d *DropletsResponse) GetStatus() string  { return d.Status }
func (d *DropletsResponse) GetMessage() string { return d.Message }

type DropletResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Droplet *Droplet `json:"droplet"`
}

func (d *DropletResponse) GetStatus() string  { return d.Status }
func (d *DropletResponse) GetMessage() string { return d.Message }

type DropletCreationResponse struct {
	Status          string           `json:"status"`
	Message         string           `json:"message"`
	DropletCreation *DropletCreation `json:"droplet"`
}

func (d *DropletCreationResponse) GetStatus() string  { return d.Status }
func (d *DropletCreationResponse) GetMessage() string { return d.Message }

type DomainsResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Domains []*Domain `json:"domains"`
}

func (d *DomainsResponse) GetStatus() string  { return d.Status }
func (d *DomainsResponse) GetMessage() string { return d.Message }

type DomainResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Domain  *Domain `json:"domain"`
}

func (d *DomainResponse) GetStatus() string  { return d.Status }
func (d *DomainResponse) GetMessage() string { return d.Message }

type RecordsResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Records []*Record `json:"records"`
}

func (r *RecordsResponse) GetStatus() string  { return r.Status }
func (r *RecordsResponse) GetMessage() string { return r.Message }

type RecordResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Record  *Record `json:"record"`
}

func (r *RecordResponse) GetStatus() string  { return r.Status }
func (r *RecordResponse) GetMessage() string { return r.Message }

type KeysResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Keys    []*Key `json:"ssh_keys"`
}

func (k *KeysResponse) GetStatus() string  { return k.Status }
func (k *KeysResponse) GetMessage() string { return k.Message }

type KeyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Key     *Key   `json:"ssh_key"`
}

func (k *KeyResponse) GetStatus() string  { return k.Status }
func (k *KeyResponse) GetMessage() string { return k.Message }

type ImagesResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Images  []*Image `json:"images"`
}

func (i *ImagesResponse) GetStatus() string  { return i.Status }
func (i *ImagesResponse) GetMessage() string { return i.Message }

type ImageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Image   *Image `json:"image"`
}

func (i *ImageResponse) GetStatus() string  { return i.Status }
func (i *ImageResponse) GetMessage() string { return i.Message }

type SizesResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Sizes   []*Size `json:"sizes"`
}

func (s *SizesResponse) GetStatus() string  { return s.Status }
func (s *SizesResponse) GetMessage() string { return s.Message }

type RegionsResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Regions []*Region `json:"regions"`
}

func (r *RegionsResponse) GetStatus() string  { return r.Status }
func (r *RegionsResponse) GetMessage() string { return r.Message }

type EventIdResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	EventId *EventId `json:"event_id"`
}

func (e *EventIdResponse) GetStatus() string  { return e.Status }
func (e *EventIdResponse) GetMessage() string { return e.Message }

type EventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Event   *Event `json:"event"`
}

func (e *EventResponse) GetStatus() string  { return e.Status }
func (e *EventResponse) GetMessage() string { return e.Message }

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (s *StatusResponse) GetStatus() string  { return s.Status }
func (s *StatusResponse) GetMessage() string { return s.Message }

func GetDroplets() ([]*Droplet, error) {
	r := &DropletsResponse{}
	err := get("/droplets/", url.Values{}, r)
	return r.Droplets, err
}

func GetDroplet(id string) (*Droplet, error) {
	r := &DropletResponse{}
	err := get(fmt.Sprintf("/droplets/%s", id), url.Values{}, r)
	return r.Droplet, err
}

func CreateDroplet(name, sizeId, imageId, regionId, keyIds string) (*DropletCreation, error) {
	r := &DropletCreationResponse{}
	q := url.Values{}
	q.Set("name", name)
	q.Set("size_id", sizeId)
	q.Set("image_id", imageId)
	q.Set("region_id", regionId)
	q.Set("ssh_key_ids", keyIds)
	err := get(fmt.Sprintf("/droplets/new"), q, r)
	return r.DropletCreation, err
}

func ShutdownDroplet(id string) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/shutdown/", id), url.Values{}, r)
	return r.EventId, err
}

func RebootDroplet(id string) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/reboot/", id), url.Values{}, r)
	return r.EventId, err
}

func PoweroffDroplet(id string) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/power_off/", id), url.Values{}, r)
	return r.EventId, err
}

func PoweronDroplet(id string) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/power_on/", id), url.Values{}, r)
	return r.EventId, err
}

func PowercycleDroplet(id string) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/power_cycle/", id), url.Values{}, r)
	return r.EventId, err
}

func ResizeDroplet(dropletId, sizeId string) (*EventId, error) {
	r := &EventIdResponse{}
	q := url.Values{}
	q.Set("size_id", sizeId)
	err := get(fmt.Sprintf("/droplets/%s/resize/", dropletId), q, r)
	return r.EventId, err
}

func SnapshotDroplet(id, name string) (*EventId, error) {
	r := &EventIdResponse{}
	q := url.Values{}
	q.Set("name", name)
	err := get(fmt.Sprintf("/droplets/%s/snapshot/", id), q, r)
	return r.EventId, err
}

func RestoreDroplet(dropletId, imageId string) (*EventId, error) {
	r := &EventIdResponse{}
	q := url.Values{}
	q.Set("image_id", imageId)
	err := get(fmt.Sprintf("/droplets/%s/restore/", dropletId), q, r)
	return r.EventId, err
}

func RebuildDroplet(dropletId, imageId string) (*EventId, error) {
	r := &EventIdResponse{}
	q := url.Values{}
	q.Set("image_id", imageId)
	err := get(fmt.Sprintf("/droplets/%s/rebuild/", dropletId), q, r)
	return r.EventId, err
}

func RenameDroplet(id, name string) (*EventId, error) {
	r := &EventIdResponse{}
	q := url.Values{}
	q.Set("name", name)
	err := get(fmt.Sprintf("/droplets/%s/rename/", id), q, r)
	return r.EventId, err
}

func ResetpassDroplet(id string) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/password_reset/", id), url.Values{}, r)
	return r.EventId, err
}

func DestroyDroplet(id string, scrub bool) (*EventId, error) {
	r := &EventIdResponse{}
	err := get(fmt.Sprintf("/droplets/%s/destroy/", id), url.Values{}, r)
	return r.EventId, err
}

func GetDomains() ([]*Domain, error) {
	r := &DomainsResponse{}
	err := get("/domains/", url.Values{}, r)
	return r.Domains, err
}

func GetDomain(id string) (*Domain, error) {
	r := &DomainResponse{}
	err := get(fmt.Sprintf("/domains/%s", id), url.Values{}, r)
	return r.Domain, err
}

func DestroyDomain(id string) error {
	return get(fmt.Sprintf("/domains/%s/destroy/", id), url.Values{}, &StatusResponse{})
}

func GetRecords(domainId string) ([]*Record, error) {
	r := &RecordsResponse{}
	err := get(fmt.Sprintf("/domains/%s/records/", domainId), url.Values{}, r)
	return r.Records, err
}

func GetRecord(domainId, recordId string) (*Record, error) {
	r := &RecordResponse{}
	err := get(fmt.Sprintf("/domains/%s/records/%s/", domainId, recordId), url.Values{}, r)
	return r.Record, err
}

func DestroyRecord(domainId, recordId string) error {
	return get(fmt.Sprintf("/domains/%s/records/%s/destroy", domainId, recordId), url.Values{}, &StatusResponse{})
}

func GetKeys() ([]*Key, error) {
	r := &KeysResponse{}
	err := get("/ssh_keys/", url.Values{}, r)
	return r.Keys, err
}

func GetKey(id string) (*Key, error) {
	r := &KeyResponse{}
	err := get(fmt.Sprintf("/ssh_keys/%s", id), url.Values{}, r)
	return r.Key, err
}

func AddKey(name, key string) (*Key, error) {
	r := &KeyResponse{}
	q := url.Values{}
	q.Set("name", name)
	q.Set("ssh_pub_key", key)
	err := get(fmt.Sprintf("/ssh_keys/new/"), q, r)
	return r.Key, err
}

func UpdateKey(id, key string) (*Key, error) {
	r := &KeyResponse{}
	q := url.Values{}
	q.Set("ssh_pub_key", key)
	err := get(fmt.Sprintf("/ssh_keys/%s/edit/", id), q, r)
	return r.Key, err
}

func DeleteKey(id string) error {
	return get(fmt.Sprintf("/ssh_keys/%s/destroy/", id), url.Values{}, &StatusResponse{})
}

func GetImages() ([]*Image, error) {
	r := &ImagesResponse{}
	err := get("/images/", url.Values{}, r)
	return r.Images, err
}

func GetImage(id string) (*Image, error) {
	r := &ImageResponse{}
	err := get(fmt.Sprintf("/images/%s/", id), url.Values{}, r)
	return r.Image, err
}

func TransferImage(imageId, regionId string) (*EventId, error) {
	r := &EventIdResponse{}
	q := url.Values{}
	q.Set("region_id", regionId)
	err := get(fmt.Sprintf("/images/%s/transfer/", imageId), q, r)
	return r.EventId, err
}

func DestroyImage(id string) error {
	return get(fmt.Sprintf("/images/%s/destroy/", id), url.Values{}, &StatusResponse{})
}

func GetRegions() ([]*Region, error) {
	r := &RegionsResponse{}
	err := get("/regions/", url.Values{}, r)
	return r.Regions, err
}

func GetSizes() ([]*Size, error) {
	r := &SizesResponse{}
	err := get("/sizes/", url.Values{}, r)
	return r.Sizes, err
}

func GetEvent(id string) (*Event, error) {
	r := &EventResponse{}
	err := get(fmt.Sprintf("/events/%s/", id), url.Values{}, r)
	return r.Event, err
}

func get(path string, query url.Values, response Response) error {
	u, err := url.Parse("https://api.digitalocean.com" + path)
	if err != nil {
		return err
	}
	query.Set("client_id", ClientId)
	query.Set("api_key", ApiKey)
	u.RawQuery = query.Encode()
	r, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(response)
	if err != nil {
		return err
	}
	status := response.GetStatus()
	if status != "OK" {
		return errors.New(fmt.Sprintf("%s: %s", status, response.GetMessage()))
	}
	return nil
}
