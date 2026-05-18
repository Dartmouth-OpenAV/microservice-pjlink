package main

import (
	"errors"
	"fmt"

	// "strconv"
	"strings"
	// "time"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

var connectionHashes = make(map[string]string)
var connectionHashesMutex sync.Mutex

func setPower(socketKey string, state string) (string, error) {
	function := "setPower"

	value := "notok"
	err := error(nil)
	maxRetries := 20
	for maxRetries > 0 {
		value, err = setPowerDo(socketKey, state)
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - fq3sdvc retrying power operation")
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setPowerDo(socketKey string, state string) (string, error) {
	function := "setPowerDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	if state == `"on"` {
		state = "1"
	} else if state == `"off"` {
		state = "0"
	} else {
		errMsg := fmt.Sprintf(function + " - 34vewvkj unrecognized state value: " + state)
		framework.AddToErrors(socketKey, errMsg)
		return state, errors.New(errMsg)
	}

	baseStr := "%1POWR "
	commandStr := baseStr + state
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - 98jhjg error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)

	framework.Log(function + " - resp: " + string(resp) + "\n")

	if resp != "%1POWR=OK" {
		errMsg := function + " - 654htd unrecognized response to power command: " + resp
		if resp != "%1POWR=ERR3" {
			framework.AddToErrors(socketKey, errMsg)
		}
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

func setVideoRoute(socketKey string, output string, state string) (string, error) {
	function := "setVideoRoute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - 4675dgh465 the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := "notok"
	err := error(nil)
	maxRetries := 20
	for maxRetries > 0 {
		framework.Log(fmt.Sprintf("%s r34qdsgrcb calling setVideoRouteDo, maxRetries = %d", function, maxRetries))

		value, err = setVideoRouteDo(socketKey, state)
		if value != "ok" { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setVideoRouteDo(socketKey string, state string) (string, error) {
	function := "setVideoRouteDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	baseStr := "%1INPT "
	stateNoQuotes := strings.Trim(state, "\"")
	commandStr := baseStr + stateNoQuotes
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - got resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - 45hfgnb error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)

	framework.Log(function + " - resp: " + string(resp) + "\n")

	if resp == "%1INPT=ERR2" {
		errMsg := function + " - 765gjh non-existent input: " + state + " for projector: " + socketKey + ", you can list possible inputs at: GET /inputs/1"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	} else if resp == "%1INPT=ERR3" {
		framework.Log(fmt.Sprintf("%s bgfbfgd543 got unavailable time error", function))
		return "unavailable", nil
	} else if resp != "%1INPT=OK" {
		errMsg := function + " - 3qr4gsrcv unrecognized response to input command: " + resp
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey) // close the socket so our retry starts from scratch
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

func setAudioAndVideoMute(socketKey string, output string, state string) (string, error) {
	function := "setAudioAndVideoMute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - 65fjv the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setMuteDo(socketKey, state, "both")
		if value != "ok" { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setAudioMute(socketKey string, output string, state string) (string, error) {
	function := "setAudioMute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - cadsvw34 the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setMuteDo(socketKey, state, "audio")
		if value != "ok" { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setVideoMute(socketKey string, output string, state string) (string, error) {
	function := "setVideoMute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - serbs56 the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setMuteDo(socketKey, state, "video")
		if value != "ok" { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setMuteDo(socketKey string, state string, which string) (string, error) {
	function := "setMuteDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	if which == "video" {
		// Allows backwards compatibility for "on" or "off"
		if state == `"on"` || state == `"true"` {
			state = "11"
		} else if state == `"off"` || state == `"false"` {
			state = "10"
		} else {
			errMsg := fmt.Sprintf(function + " - svfdsq43 unrecognized state value: " + state)
			framework.AddToErrors(socketKey, errMsg)
			return state, errors.New(errMsg)
		}
	} else if which == "audio" {
		if state == `"on"` || state == `"true"` {
			state = "21"
		} else if state == `"off"` || state == `"false"` {
			state = "20"
		} else {
			errMsg := fmt.Sprintf(function + " - 345wg unrecognized state value: " + state)
			framework.AddToErrors(socketKey, errMsg)
			return state, errors.New(errMsg)
		}
	} else if which == "both" {
		if state == `"on"` || state == `"true"` {
			state = "31"
		} else if state == `"off"` || state == `"false"` {
			state = "30"
		} else {
			errMsg := fmt.Sprintf(function + " - 89dsi unrecognized state value: " + state)
			framework.AddToErrors(socketKey, errMsg)
			return state, errors.New(errMsg)
		}
	}

	baseStr := "%1AVMT "
	commandStr := baseStr + state
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)

	framework.Log(function + " - resp: " + string(resp) + "\n")

	if resp != "%1AVMT=OK" {
		errMsg := "uninitialized"
		if resp == "%1AVMT=ERR2" { // we tried to invoke audio or video mute separately, and this device can't do that
			errMsg = function + " - PJLink device " + socketKey + " refused to set audio or video mute separately.  Consider using the audioandvideomute endpoint."
		} else {
			errMsg = function + " - vdfssdv unrecognized response to mute command: " + resp
		}
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

func setVolume(socketKey string, state string) (string, error) {
	function := "setVolume"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	if state == `"up"` {
		state = "1"
	} else if state == `"down"` {
		state = "0"
	} else {
		errMsg := fmt.Sprintf(function + " -  45dasfq unrecognized state value: " + state)
		framework.AddToErrors(socketKey, errMsg)
		return state, errors.New(errMsg)
	}

	baseStr := "%2SVOL "
	commandStr := baseStr + state
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - dfewear error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)

	framework.Log(function + " - resp: " + string(resp) + "\n")

	if resp != "%2SVOL=OK" {
		errMsg := function + " - efvqrrq343 unrecognized response to volume command: " + resp
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return
	return "ok", nil
}

func getPower(socketKey string) (string, error) {
	// function := "getPower"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getPowerDo(socketKey)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getPowerDo(socketKey string) (string, error) {
	function := "getPowerDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	baseStr := "%1POWR ?"
	commandStr := baseStr
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " -  afred42  error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)
	framework.Log(function + " - resp: " + string(resp) + "\n")

	value := `"unknown"`
	if resp == "%1POWR=0" || resp == "%1POWR=2" { // 2 indicates cooling off
		value = `"off"`
	} else if resp == "%1POWR=1" || resp == "%1POWR=3" { // 3 indicates warming up
		value = `"on"`
	} else {
		framework.AddToErrors(socketKey, socketKey+" - 2q34awev can't interpret response")
	}

	// framework.Log( deviceStates[socketKey]["Refresh power"] ) ;

	// If we got here, the response was good, so successful return with the state indication
	return value, nil
}

func getAudioAndVideoMute(socketKey string, output string) (string, error) {
	function := "getAudioMute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - c34qev the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getMuteDo(socketKey, "both")
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getAudioMute(socketKey string, output string) (string, error) {
	function := "getAudioMute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - 3e332 the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getMuteDo(socketKey, "audio")
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getVideoMute(socketKey string, output string) (string, error) {
	function := "getVideoMute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - q34ewrdsf the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getMuteDo(socketKey, "video")
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getMuteDo(socketKey string, which string) (string, error) {
	function := "getVideoMuteDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	baseStr := "%1AVMT ?"
	commandStr := baseStr
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - ase3232 error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)
	framework.Log(function + " - resp: " + string(resp) + "\n")

	value := `"unknown"`
	if which == "video" {
		if resp == "%1AVMT=10" || resp == "%1AVMT=30" {
			value = `"off"`
		} else if resp == "%1AVMT=11" || resp == "%1AVMT=31" {
			value = `"on"`
		} else if resp == "%1AVMT=ERR3" { // the dreaded PJLink "unavailable time"
			// when the projector is off
			value = `"off"`
		} else {
			framework.AddToErrors(socketKey, socketKey+" - 233rd32 can't interpret response")
		}
	} else if which == "audio" || which == "both" {
		if resp == "%1AVMT=20" || resp == "%1AVMT=30" {
			value = `"off"`
		} else if resp == "%1AVMT=21" || resp == "%1AVMT=31" {
			value = `"on"`
		} else if resp == "%1AVMT=ERR3" { // the dreaded PJLink "unavailable time"
			// when the projector is off
			value = `"off"`
		} else {
			framework.AddToErrors(socketKey, socketKey+" - fdsvvfrs45 can't interpret response")
		}
	}

	// framework.Log( deviceStates[socketKey]["Refresh power"] ) ;

	// If we got here, the response was good, so successful return with the state indication
	return value, nil
}

func getInputs(socketKey string, output string) (string, error) {
	function := "getInputs"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - dfvrsd34 the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getInputsDo(socketKey)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getInputsDo(socketKey string) (string, error) {
	function := "getInputsDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	baseStr := "%1INST ?"
	commandStr := baseStr
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - aefvea error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)
	framework.Log(function + " - resp: " + string(resp) + "\n")

	value := `"unknown"`
	if resp == "%1INST=ERR3" { // the dreaded PJLink "unavailable time"
		framework.AddToErrors(socketKey, socketKey+" - 324edsafa unavailable time happened")
	} else {
		respSplit := strings.Split(resp, "=")
		if len(respSplit) >= 2 {
			value = `"` + strings.Join(respSplit[1:], ", ") + `"`
		}
	}

	// If we got here, the response was good, so successful return with the state indication
	return value, nil
}

func getVideoRoute(socketKey string, output string) (string, error) {
	function := "getVideoRoute"

	if output != "1" {
		errMsg := fmt.Sprintf(function + " - fdsvvdf4 the only possible videoroute output of projectors is 1")
		framework.AddToErrors(socketKey, errMsg)
	}

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getVideoRouteDo(socketKey)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getVideoRouteDo(socketKey string) (string, error) {
	function := "getVideoRouteDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	baseStr := "%1INPT ?"
	commandStr := baseStr
	// resp, err := SendCommand(socketKey, commandStr, baseStr)
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))
	//framework.Log(fmt.Sprintf(function + " - ot resp: " + string(resp) + "\n")
	if sent != true {
		errMsg := fmt.Sprintf(function + " - vffsvd756 error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)
	framework.Log(function + " - resp: " + string(resp) + "\n")

	value := `"unknown"`
	if resp != "%1INPT=ERR3" && resp != "" { // projector might be off or just didn't respond
		respSplit := strings.Split(resp, "=")
		if len(respSplit) != 2 {
			framework.AddToErrors(socketKey, socketKey+" - jmj67765 can't interpret response #eiuthbriut")
		} else {
			val, err := strconv.Atoi(respSplit[1])
			_ = val // appease
			if err != nil {
				framework.AddToErrors(socketKey, socketKey+" - 345wgg can't interpret response #3iutherb")
			} else {
				value = `"` + respSplit[1] + `"`
			}
		}
	}

	// framework.Log( deviceStates[socketKey]["Refresh power"] ) ;

	// If we got here, the response was good, so successful return with the state indication
	return value, nil
}

func getLampHours(socketKey string) (string, error) {
	//function := "getLampHours"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getLampHoursDo(socketKey)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}
	return value, err
}

func getLampHoursDo(socketKey string) (string, error) {
	function := "getLampHoursDo"

	PJLinkEstablishConnectionIfNeeded(socketKey)

	baseStr := "%1LAMP ?"
	commandStr := baseStr
	sent := framework.WriteLineToSocket(socketKey, PJLinkizeCommand(socketKey, commandStr))

	if sent != true {
		errMsg := fmt.Sprintf(function + " - d8ruvq2a error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp := framework.ReadLineFromSocket(socketKey)
	framework.Log(function + "a2cv7z94 - resp: " + string(resp) + "\n")

	value := `"unknown"`
	if resp != "" {
		respSplit := strings.Split(resp, "=")
		if len(respSplit) != 2 {
			framework.AddToErrors(socketKey, fmt.Sprintf(socketKey+" - 9tvtwy0i can't interpret response: %s", resp))
			if respSplit[1] == "ERR3" {
				framework.Log(fmt.Sprintf("%s - pbqtb3ug got unavailable time error", function))
				return "unavailable", nil
			}
		} else {
			value = `"` + respSplit[1] + `"`
		}
	}
	// If we got here, the response was good, so successful return with the state indication
	return value, nil
}

func PJLinkEstablishConnectionIfNeeded(socketKey string) {
	function := "PJLinkEstablishConnectionIfNeeded"

	connectionEstablished := framework.CheckConnectionsMapExists(socketKey) // PJLink is always TCP

	if !connectionEstablished {
		framework.Log(fmt.Sprintf(function + " - 4f3evebv connection needs to be established\n"))

		resp := framework.ReadLineFromSocket(socketKey)
		framework.Log(fmt.Sprintf(function + " - resp: " + string(resp) + "\n"))

		// the first thing we should ever get from the PJLink protocol is
		//   PJLINK 0           if not authentication is set
		//   PJLINK 1 xxxxxxxx  where xx is a random hash specific to this connection, if authentication is set on the device

		responseSplit := strings.Split(resp, " ")

		if len(responseSplit) < 2 {
			framework.AddToErrors(socketKey, socketKey+" - afs453e it doesn't look like I'm talking to a device that implemented the PJLink protocol properly. I got an unknown response with "+resp)
		} else {
			if responseSplit[0] != "PJLINK" {
				framework.AddToErrors(socketKey, socketKey+" - vawerw345 it doesn't look like I'm talking to a device that implemented the PJLink protocol properly. I got an unknown response with "+resp)
			}

			connectionHash := ""
			if responseSplit[1] != "0" && responseSplit[1] != "1" {
				framework.AddToErrors(socketKey, socketKey+" - 324q2 I'm not sure wheter to treat this connection as authenticated or not")
			}

			if responseSplit[1] == "1" {
				connectionHash = responseSplit[2]
			}

			passwordIfAny := ""
			if strings.Count(socketKey, "@") == 1 {
				credentials := strings.Split(socketKey, "@")[0]
				if strings.Count(credentials, ":") == 1 {
					passwordIfAny = strings.Split(credentials, ":")[1]
				}
				// framework.Log("passwordIfAny: " + passwordIfAny + "\n\n")
			}
			connectionHashesMutex.Lock()
			connectionHashes[socketKey] = GetMD5Hash(connectionHash + passwordIfAny)
			framework.Log(socketKey + " - connectionHash[" + connectionHash + "], passwordIfAny[" + passwordIfAny + "] hash to: " + connectionHashes[socketKey])
			connectionHashesMutex.Unlock()
		}
	}
}

func PJLinkizeCommand(socketKey string, PJLinkCommand string) string {
	connectionHashesMutex.Lock()
	defer connectionHashesMutex.Unlock()
	return connectionHashes[socketKey] + PJLinkCommand + "\r"
}

// from https://stackoverflow.com/a/25286918
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func healthCheck(socketKey string) (string, error) {
	resp, err := getLampHours(socketKey)
	returnStr := "true"
	if err != nil && strings.Contains(err.Error(), "error sending command") {
		returnStr = "no connection"
	} else if err != nil {
		returnStr = "false"
	} else {
		returnStr = returnStr + " lamp hours: " + strings.Trim(resp, `"`)
	}
	return `"` + returnStr + `"`, nil
}
