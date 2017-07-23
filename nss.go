package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// IDStart is the userid of the first smith user
	IDStart = 10
)

// Users populates the passwd, group and nsswitch.conf
// files in the etc directory relative to the outputDir
// parameter. The users in the users parameter will be
// populated in the etc/passwd file along with some
// predetermined users that are required in the container image.
// The predetermined users are currently root, bin and daemon.
// The etc/group and etc/nsswitch.conf files written are not
// currently influenced at all by the function's parameters.
// This function will overwrite any existing passwd, group or
// nsswitch.conf files.
func Users(outputDir string, users []string) error {
	etcDir := filepath.Join(outputDir, "etc")
	if err := os.MkdirAll(etcDir, 0755); err != nil {
		return err
	}
	// add a group for each user
	if err := groups(outputDir, users); err != nil {
		return err
	}
	s := []string{
		"root:x:0:0:root:/write:",
		"bin:x:1:0:bin:/bin:",
		"daemon:x:2:0:daemon:/bin:",
	}
	for i, user := range users {
		s = append(s, fmt.Sprintf("%s:x:%d:%d:%s:/write:", user, IDStart+i, IDStart+i, user))
	}
	path := filepath.Join(etcDir, "passwd")
	if err := ioutil.WriteFile(path, []byte(strings.Join(s, "\n")), 0644); err != nil {
		return err
	}
	return nss(outputDir)
}

func groups(outputDir string, groups []string) error {
	s := []string{
		"root:x:0:",
		"daemon:x:1:",
		"bin:x:2:",
		"sys:x:3:",
		"adm:x:4:",
		"tty:x:5:",
	}
	for i, group := range groups {
		s = append(s, fmt.Sprintf("%s:x:%d:", group, IDStart+i))
	}
	path := filepath.Join(outputDir, "etc", "group")
	if err := ioutil.WriteFile(path, []byte(strings.Join(s, "\n")), 0644); err != nil {
		return err
	}
	return nil
}

func nss(outputDir string) error {
	s := []string{
		"passwd:     files",
		"shadow:     files",
		"group:      files",
		"hosts:      files dns",
	}
	path := filepath.Join(outputDir, "etc", "nsswitch.conf")
	if err := ioutil.WriteFile(path, []byte(strings.Join(s, "\n")), 0644); err != nil {
		return err
	}
	return nil
}
