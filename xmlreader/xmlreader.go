package xmlreader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

func readXMLFile(fn string) (Ldml, error) {
	res := Ldml{}

	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return res, fmt.Errorf("failed to read XML file : %v", err)
	}

	err = xml.Unmarshal(bytes, &res)
	if err != nil {
		return res, fmt.Errorf("failed to peocess XML file : %v", err)
	}

	return res, nil
}
