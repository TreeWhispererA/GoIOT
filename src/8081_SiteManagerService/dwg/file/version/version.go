package version

import (
	"errors"
	"strconv"
)

type ACadVersion uint16

const (
	INVALID_VERSION ACadVersion = 0
	AC1_2           ACadVersion = 102
	AC1_4           ACadVersion = 104
	AC1_50          ACadVersion = 150
	AC2_10          ACadVersion = 210
	AC1002          ACadVersion = 1002
	AC1003          ACadVersion = 1003
	AC1004          ACadVersion = 1004
	AC1006          ACadVersion = 1006
	AC1009          ACadVersion = 1009
	AC1012          ACadVersion = 1012
	AC1014          ACadVersion = 1014
	AC1015          ACadVersion = 1015
	AC1018          ACadVersion = 1018
	AC1021          ACadVersion = 1021
	AC1024          ACadVersion = 1024
	AC1027          ACadVersion = 1027
	AC1032          ACadVersion = 1032
)

func GetVersionFromName(name string) (ACadVersion, error) {
	if name[:2] != "AC" {
		return INVALID_VERSION, errors.New("Invalid version string")
	}
	versionNum, err := strconv.ParseInt(string(name[2:]), 10, 64)
	if err != nil {
		return INVALID_VERSION, err
	}

	return ACadVersion(versionNum), nil
}
