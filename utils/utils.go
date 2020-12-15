// Copyright 2018 Open Networking Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//Package utils gives type conversion helpers
package utils

import (
	"encoding/hex"
	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
	"strings"
)

var OnuSNhex bool

//OnuSnHexEncode converts ONU sn from human readable format like 'SCOM00001B6D' to hex like '53434F4D00001B6D'
func OnuSnHexEncode(onuSn string) ( string) {
        if len(onuSn) != 12 {
                logger.Warn("ignoring serial to hex encode, 12 chars are expected for onuSn but %d are provided", len(onuSn))
                return onuSn
        }

        output := hex.EncodeToString([]byte(onuSn[0:4])) + onuSn[4:12]
        return strings.ToUpper(output)
}

func GetOnuSN(onuSN string) (string) {
	if OnuSNhex  {
                return OnuSnHexEncode(onuSN)
        } else {
                return onuSN
        }
}
