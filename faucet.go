package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whub/faucet/cmd"
	"github.com/whub/faucet/fancy"
	"github.com/whub/faucet/sand"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

type Config struct {
	ClientId string `json:"clientId"`
	ApiKey   string `json:"apiKey"`
}

func loadConfig() {
	f, err := os.Open("faucet.json")
	if err != nil {
		fancy.Println(fancy.Red, err)
		os.Exit(1)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	config := Config{}
	err = dec.Decode(&config)
	if err != nil {
		fancy.Println(fancy.Red, err)
		os.Exit(1)
	}
	sand.ClientId = config.ClientId
	sand.ApiKey = config.ApiKey
}

func main() {
	loadConfig()

	root := cmd.Root(os.Args[0])

	droplets := root.Parent("droplets", "manage droplets")
	droplets.Command("list", "list droplets", "", dropletsList)
	droplets.Command("show", "show details for a droplet", "<droplet id>", dropletsShow)
	droplets.Command("new", "create a new droplet", "", dropletsNew)
	droplets.Command("ssh", "ssh into a droplet", "<droplet id>", dropletsSSH)
	droplets.Command("scp", "scp a file to a droplet", "<file> <droplet id>", dropletsSCP)
	droplets.Command("open", "open the droplet's ip address in a browser", "<droplet id>", dropletsOpen)
	droplets.Command("shutdown", "cleanly shutdown a droplet", "<droplet id>", dropletsShutdown)
	droplets.Command("reboot", "cleanly reboot a droplet", "<droplet id>", dropletsReboot)
	droplets.Command("poweroff", "power off a droplet", "<droplet id>", dropletsPoweroff)
	droplets.Command("poweron", "power on a droplet", "<droplet id>", dropletsPoweron)
	droplets.Command("powercycle", "power off then power on a droplet", "<droplet id>", dropletsPowercycle)
	droplets.Command("resize", "change the size of a droplet", "<droplet id> <size id>", dropletsResize)
	droplets.Command("snapshot", "take a snapshot of a droplet", "<droplet id> <name>", dropletsSnapshot)
	droplets.Command("restore", "revert a droplet back to a snapshot", "<droplet id> <image id>", dropletsRestore)
	droplets.Command("rebuild", "reinstall an image to a droplet", "<droplet id> <image id>", dropletsRebuild)
	droplets.Command("rename", "change the name of a droplet", "<droplet id> <name>", dropletsRename)
	droplets.Command("resetpass", "reset the root password of a droplet", "<droplet id>", dropletsResetpass)
	droplets.Command("destroy", "destroy a droplet", "<droplet id> <scrub data?>", dropletsDestroy)

	domains := root.Parent("domains", "manage domains")
	domains.Command("list", "list domains", "", domainsList)
	domains.Command("show", "show details of a domain", "<domain id>", domainsShow)
	domains.Command("new", "create a new domain", "", domainsNew)
	domains.Command("destroy", "destroy a domain", "<domain id>", domainsDestroy)

	records := domains.Parent("records", "manage records")
	records.Command("list", "list records", "<domain id>", recordsList)
	records.Command("show", "show details for a record", "<domain id> <record id>", recordsShow)
	records.Command("new", "create a new record", "<domain id>", recordsNew)
	records.Command("edit", "edit a record", "<domain id> <record id>", recordsEdit)
	records.Command("destroy", "destroy a record", "<domain id> <record id>", recordsDestroy)

	keys := root.Parent("keys", "manage ssh keys")
	keys.Command("list", "list keys", "", keysList)
	keys.Command("show", "show details of a key", "<key id>", keysShow)
	keys.Command("add", "add ~/.ssh/id_rsa.pub to the key list", "<name>", keysAdd)
	keys.Command("update", "change a key to match ~/.ssh/id_rsa.pub", "<key id>", keysUpdate)
	keys.Command("delete", "delete a key", "<key id>", keysDelete)

	images := root.Parent("images", "manage images")
	images.Command("list", "list images", "", imagesList)
	images.Command("show", "show details of an image", "<image id>", imagesShow)
	images.Command("transfer", "transfer an image to a region", "<image id> <region id>", imagesTransfer)
	images.Command("destroy", "destroy an image", "<image id>", imagesDestroy)

	root.Command("regions", "list available regions", "", regions)
	root.Command("sizes", "list available sizes", "", sizes)
	root.Command("event", "show progress of an event", "<event id>", event)
	root.Command("help", "show usage for a specific command", "<command>", help)

	err := root.Dispatch(os.Args, 1)
	if err != nil {
		fancy.Println(fancy.Red, err)
	}
}

func dropletsList(args []string) error {
	if len(args) != 0 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching droplets... ")
	droplets, err := sand.GetDroplets()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	if len(droplets) == 0 {
		fmt.Println("No droplets.")
	}
	for _, d := range droplets {
		DropletPrint(d)
	}
	return nil
}

func dropletsShow(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching droplet... ")
	d, err := sand.GetDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	DropletPrint(d)
	return nil
}

func dropletsNew(args []string) error {
	fmt.Print("name: ")
	var name string
	_, err := fmt.Scanln(&name)
	if err != nil {
		return err
	}
	err = sizes([]string{})
	if err != nil {
		return err
	}
	fmt.Print("size id: ")
	var sizeId string
	_, err = fmt.Scanln(&sizeId)
	if err != nil {
		return err
	}
	err = imagesList([]string{})
	if err != nil {
		return err
	}
	fmt.Print("image id: ")
	var imageId string
	_, err = fmt.Scanln(&imageId)
	if err != nil {
		return err
	}
	err = regions([]string{})
	if err != nil {
		return err
	}
	fmt.Print("region id: ")
	var regionId string
	_, err = fmt.Scanln(&regionId)
	if err != nil {
		return err
	}
	err = keysList([]string{})
	if err != nil {
		return err
	}
	fmt.Print("key ids (comma separated): ")
	var keyIds string
	_, err = fmt.Scanln(&keyIds)
	if err != nil {
		return err
	}
	fmt.Print("creating droplet... ")
	d, err := sand.CreateDroplet(name, sizeId, imageId, regionId, keyIds)
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	DropletCreationPrint(d)
	return nil
}

func dropletsSSH(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching droplet... ")
	d, err := sand.GetDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	fmt.Println("running ssh...")
	command := exec.Command("ssh", "root@"+d.IPAddress)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func dropletsSCP(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching droplet... ")
	d, err := sand.GetDroplet(args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	fmt.Println("running scp...")
	command := exec.Command("scp", args[0], "root@"+d.IPAddress+":")
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func dropletsOpen(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching droplet... ")
	d, err := sand.GetDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	fmt.Println("opening...")
	command := exec.Command("open", "http://"+d.IPAddress)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func dropletsShutdown(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing shutdown command... ")
	e, err := sand.ShutdownDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsReboot(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing reboot command... ")
	e, err := sand.RebootDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsPoweroff(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing poweroff command... ")
	e, err := sand.PoweroffDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsPoweron(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing poweron command... ")
	e, err := sand.PoweronDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsPowercycle(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing powercycle command... ")
	e, err := sand.PowercycleDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsResize(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing resize command... ")
	e, err := sand.ResizeDroplet(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsSnapshot(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing snapshot command... ")
	e, err := sand.SnapshotDroplet(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsRestore(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing restore command... ")
	e, err := sand.RestoreDroplet(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsRebuild(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing rebuild command... ")
	e, err := sand.RebuildDroplet(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsRename(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing rename command... ")
	e, err := sand.RenameDroplet(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsResetpass(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing resetpass command... ")
	e, err := sand.ResetpassDroplet(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func dropletsDestroy(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	scrub, err := strconv.ParseBool(args[1])
	if err != nil {
		return err
	}
	fmt.Print("issuing destroy command... ")
	e, err := sand.DestroyDroplet(args[0], scrub)
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func domainsList(args []string) error {
	if len(args) != 0 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching domains... ")
	domains, err := sand.GetDomains()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	for _, d := range domains {
		DomainPrint(d)
	}
	return nil
}

func domainsShow(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching domain... ")
	d, err := sand.GetDomain(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	DomainPrint(d)
	return nil
}

func domainsNew(args []string) error {
	return errors.New("not implemented")
}

func domainsDestroy(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("destroying the domain... ")
	err := sand.DestroyDomain(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	return nil
}

func recordsList(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching records... ")
	records, err := sand.GetRecords(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	for _, r := range records {
		RecordPrint(r)
	}
	return nil
}

func recordsShow(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching record... ")
	r, err := sand.GetRecord(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	RecordPrint(r)
	return nil
}

func recordsNew(args []string) error {
	return errors.New("not implemented")
}

func recordsEdit(args []string) error {
	return errors.New("not implemented")
}

func recordsDestroy(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("destroying record... ")
	err := sand.DestroyRecord(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	return nil
}

func keysList(args []string) error {
	if len(args) != 0 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching ssh keys... ")
	keys, err := sand.GetKeys()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	for _, k := range keys {
		KeyPrint(k)
	}
	return nil
}

func keysShow(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching ssh key... ")
	k, err := sand.GetKey(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	KeyPrint(k)
	return nil
}

func keysAdd(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("looking for local key... ")
	keyStr, err := readPublicKey()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	fmt.Print("uploading key... ")
	k, err := sand.AddKey(args[0], keyStr)
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	KeyPrint(k)
	return nil
}

func keysUpdate(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("looking for local key... ")
	keyStr, err := readPublicKey()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	fmt.Print("updating remote key to match... ")
	k, err := sand.UpdateKey(args[0], keyStr)
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	KeyPrint(k)
	return nil
}

func keysDelete(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("deleting key... ")
	err := sand.DeleteKey(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	return nil
}

func imagesList(args []string) error {
	if len(args) != 0 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching images... ")
	images, err := sand.GetImages()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	for _, image := range images {
		ImagePrint(image)
	}
	return nil
}

func imagesShow(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching image... ")
	image, err := sand.GetImage(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	ImagePrint(image)
	return nil
}

func imagesTransfer(args []string) error {
	if len(args) != 2 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("issuing transfer command... ")
	e, err := sand.TransferImage(args[0], args[1])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventIdPrint(e)
	return nil
}

func imagesDestroy(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("destroying image... ")
	err := sand.DestroyImage(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	return nil
}

func regions(args []string) error {
	if len(args) != 0 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching regions... ")
	regions, err := sand.GetRegions()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	for _, r := range regions {
		RegionPrint(r)
	}
	return nil
}

func sizes(args []string) error {
	if len(args) != 0 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching sizes... ")
	sizes, err := sand.GetSizes()
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	for _, s := range sizes {
		SizePrint(s)
	}
	return nil
}

func event(args []string) error {
	if len(args) != 1 {
		return cmd.ErrInvalidArgs
	}
	fmt.Print("fetching event status... ")
	e, err := sand.GetEvent(args[0])
	if err != nil {
		return err
	}
	fancy.Println(fancy.Green, "OK")
	EventPrint(e)
	return nil
}

func help(args []string) error {
	return errors.New("not implemented")
}

func readPublicKey() (string, error) {
	b, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub")
	return string(b), err
}

func DropletPrint(d *sand.Droplet) {
	fancy.Print(fancy.Blue, d.Name)
	fmt.Printf(` {
  Id: %d
  ImageId: %d
  SizeId: %d
  RegionId: %d
  BackupsActive: %t
  IPAddress: %s
  PrivateIPAddress: %s
  Locked: %t
  Status: %s
  CreatedAt: %v
}
`, d.Id, d.ImageId, d.SizeId, d.RegionId, d.BackupsActive,
		d.IPAddress, d.PrivateIPAddress, d.Locked, d.Status, d.CreatedAt)
}

func DropletCreationPrint(d *sand.DropletCreation) {
	fancy.Print(fancy.Blue, d.Name)
	fmt.Printf(` {
  Id: %d
  ImageId: %d
  SizeId: %d
  EventId: %d
}
`, d.Id, d.ImageId, d.SizeId, d.EventId)
}

func DomainPrint(d *sand.Domain) {
	fancy.Print(fancy.Blue, d.Name)
	fmt.Printf(` {
  Id: %d
  TTL: %d
  LiveZoneFile:
%s
  Error: %s
  ZoneFileWithError:
%s
}
`, d.Id, d.TTL, d.LiveZoneFile, d.Error, d.ZoneFileWithError)
}

func RecordPrint(r *sand.Record) {
	fancy.Print(fancy.Blue, r.Name)
	fmt.Printf(` {
  Id: %d
  DomainId: %d
  RecordType: %s
  Data: %s
  Priority: %d
  Port: %s
  Weight: %s
}
`, r.Id, r.DomainId, r.RecordType, r.Data, r.Priority, r.Port, r.Weight)
}

func EventIdPrint(e *sand.EventId) {
	fmt.Print("Event Id: ")
	fancy.Println(fancy.Blue, int(*e))
}

func EventPrint(e *sand.Event) {
	fancy.Println(fancy.Blue, e.Id)
	fmt.Printf(` {
  Status: %s
  DropletId: %d
  EventType: %d
  Percentage: %s
}
`, e.Status, e.DropletId, e.EventType, e.Percentage)
}

func KeyPrint(k *sand.Key) {
	fancy.Print(fancy.Blue, k.Name)
	if k.PublicKey != "" {
		fmt.Printf(` {
  Id: %d
  PublicKey: %s
}
`, k.Id, k.PublicKey)
	} else {
		fmt.Printf(` {
  Id: %d
}
`, k.Id)
	}
}

func ImagePrint(i *sand.Image) {
	fancy.Print(fancy.Blue, i.Name)
	fmt.Printf(` {
  Id: %d
  Distribution: %s
}
`, i.Id, i.Distribution)
}

func RegionPrint(r *sand.Region) {
	fancy.Print(fancy.Blue, r.Name)
	fmt.Printf(` {
  Id: %d
}
`, r.Id)
}

func SizePrint(s *sand.Size) {
	fancy.Print(fancy.Blue, s.Name)
	fmt.Printf(` {
  Id: %d
}
`, s.Id)
}
