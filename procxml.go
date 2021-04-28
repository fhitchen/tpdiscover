package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Endurox struct {
	XMLName   xml.Name `xml:"endurox"`
	Comment      string   `xml:",comment"`
	Appconfig struct {
		Comment        string `xml:",comment"`
		Sanity         string `xml:"sanity"`
		Brrefresh      string `xml:"brrefresh"`
		RestartMin     string `xml:"restart_min"`
		RestartStep    string `xml:"restart_step"`
		RestartMax     string `xml:"restart_max"`
		RestartToCheck string `xml:"restart_to_check"`
		GatherPqStats  string `xml:"gather_pq_stats"`
	} `xml:"appconfig"`
	Defaults struct {
		Comment     string `xml:",comment"`
		Min      string `xml:"min"`
		Max      string `xml:"max"`
		Autokill string `xml:"autokill"`
		StartMax string `xml:"start_max"`
		Pingtime string `xml:"pingtime"`
		PingMax  string `xml:"ping_max"`
		EndMax   string `xml:"end_max"`
		Killtime string `xml:"killtime"`
	} `xml:"defaults"`
	Servers struct {
		Comment   string `xml:",comment"`
		Server []struct {
			Comment   string `xml:",comment"`
			Name   string `xml:"name,attr"`
			Min    string `xml:"min"`
			Max    string `xml:"max"`
			Srvid  string `xml:"srvid"`
			Sysopt string `xml:"sysopt"`
			Cctag  string `xml:"cctag"`
			Appopt string `xml:"appopt"`
		} `xml:"server"`
	} `xml:"servers"`
	Clients struct {
		Comment   string `xml:",comment"`
		Client []struct {
			Comment    string `xml:",comment"`
			Cmdline string `xml:"cmdline,attr"`
			Exec    []struct {
				Comment      string `xml:",comment"`
				Tag       string `xml:"tag,attr"`
				Subsect   string `xml:"subsect,attr"`
				Autostart string `xml:"autostart,attr"`
				Log       string `xml:"log,attr"`
			} `xml:"exec"`
		} `xml:"client"`
	} `xml:"clients"`
}


func ReadNdrxconfig() (e Endurox) {

	v := Endurox{}

//	buff, err := ioutil.ReadFile("ndrxconfig.xml")
	buff, err := ioutil.ReadFile("x.xml")
	if err != nil {
		fmt.Print(err)
	}

	err = xml.Unmarshal(buff, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("The lot: %#v\n", v)
	fmt.Printf("Count: %d\n", len(v.Servers.Server))
	fmt.Printf("Count: %#v\n", v.Servers.Server[0])

	for _, e := range v.Servers.Server {
		fmt.Printf("Name: %q Sysopt: %q\n", e.Name, e.Sysopt)
		fmt.Printf("Appopt: %q\n", e.Appopt)
	}


	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return v

}

