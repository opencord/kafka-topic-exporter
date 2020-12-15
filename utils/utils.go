//Package utils gives type conversion helpers
package utils

import (
	"encoding/hex"
	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
)

var OnuSNhex bool

//OnuSnHexEncode converts ONU sn from human readable format like 'SCOM00001B6D' to hex like '53434F4D00001B6D'
func OnuSnHexEncode(onuSn string) ( string) {
        if len(onuSn) != 12 {
                logger.Warn("ignoring serial to hex encode, 12 chars are expected for onuSn but %d are provided", len(onuSn))
                return onuSn
        }

        output := hex.EncodeToString([]byte(onuSn[0:4])) + onuSn[4:12]
        return string(output)
}

func GetOnuSN(onuSN string) (string) {
	if OnuSNhex  {
                return OnuSnHexEncode(onuSN)
        } else {
                return onuSN
        }
}
