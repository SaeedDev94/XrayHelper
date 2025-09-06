package tools

import (
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	tagTools        = "tools"
	packageListPath = "/data/system/packages.list"
)

var packageMap = make(map[string]string)

// loadPackage load and parse Android package with uid list into a map
func loadPackage() {
	if len(packageMap) > 0 {
		return
	}
	packageListFile, err := os.Open(packageListPath)
	if err != nil {
		log.HandleDebug("load package failed, " + err.Error())
		return
	}
	defer packageListFile.Close()
	packageScanner := bufio.NewScanner(packageListFile)
	packageScanner.Split(bufio.ScanLines)
	for packageScanner.Scan() {
		packageInfo := strings.Fields(packageScanner.Text())
		if len(packageInfo) >= 2 {
			packageMap[packageInfo[0]] = packageInfo[1]
		}
	}
	log.HandleDebug(packageMap)
}

func GetUid(pkgInfo string) []string {
	loadPackage()
	var (
		userId    int
		pkgUserId []string
	)
	info := strings.Split(pkgInfo, ":")
	if len(info) == 2 {
		userId, _ = strconv.Atoi(info[1])
	}
	for pkgStr, pkgIdStr := range packageMap {
		if common.WildcardMatch(pkgStr, info[0]) {
			pkgId, _ := strconv.Atoi(pkgIdStr)
			pkgUserIdStr := strconv.Itoa(userId*100000 + pkgId)
			pkgUserId = append(pkgUserId, pkgUserIdStr)
		}
	}
	return pkgUserId
}

func DisableIPV6DNS() error {
	if err := common.Ipt6.Insert("filter", "OUTPUT", 1, "-p", "udp", "--dport", "53", "-j", "REJECT"); err != nil {
		return e.New("disable dns request on ipv6 failed, ", err).WithPrefix(tagTools)
	}
	return nil
}

func EnableIPV6DNS() {
	_ = common.Ipt6.Delete("filter", "OUTPUT", "-p", "udp", "--dport", "53", "-j", "REJECT")
}
