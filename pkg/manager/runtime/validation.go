package runtime

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/asaskevich/govalidator"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

func ValidateName(name string) error {
	if !govalidator.StringLength(name, NameMinLength, NameMaxLength) {
		return fmt.Errorf("the length of the name should be 1 to 255")
	}
	return nil
}

func ValidateURL(url string) error {
	if !govalidator.IsURL(url) {
		return fmt.Errorf("url format error")
	}
	return nil
}

func ValidateLabelString(labelString string) error {
	mapLabel, err := url.ParseQuery(labelString)
	if err != nil {
		return err
	}
	err = ValidateLabelMapFmt(mapLabel)
	if err != nil {
		return err
	}
	err = ValidateLabelMapContent(mapLabel)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSelectorString(selectorString string) error {
	selectorMap, err := url.ParseQuery(selectorString)
	if err != nil {
		return err
	}
	err = ValidateSelectorMapFmt(selectorMap)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSelectorMapFmt(selectorMap map[string][]string) error {
	for sKey, sValues := range selectorMap {
		err := ValidateLabelKey(sKey)
		if err != nil {
			return err
		}
		for _, sValue := range sValues {
			err := ValidateLabelValue(sValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateLabelMapFmt(labelMap map[string][]string) error {
	for mKey, mValue := range labelMap {
		if len(mValue) != 1 {
			return fmt.Errorf("label format error ")
		}
		err := ValidateLabelKey(mKey)
		if err != nil {
			return err
		}
		err = ValidateLabelValue(mValue[0])
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateLabelMapContent(labelMap map[string][]string) error {
	runtimeValue, ok := labelMap[LabelRuntime]
	if !ok {
		return fmt.Errorf("label [runtime] is required")
	}

	if i := utils.FindString(VmBaseRuntime, runtimeValue[0]); i != -1 {
		if _, ok := labelMap[LabelZone]; !ok {
			return fmt.Errorf("vm-based runtime need label [zone]")
		}
		return nil
	}

	if i := utils.FindString(CmBaseRuntime, runtimeValue[0]); i != -1 {
		return nil
	}
	return fmt.Errorf("runtime not support")
}

var LabelNameRegexp = regexp.MustCompile(LabelKeyFmt)

func ValidateLabelKey(labelName string) error {
	if !govalidator.StringLength(labelName, LabelKeyMinLength, LabelKeyMaxLength) {
		return fmt.Errorf("the length of the label key should be 1 to 50")
	}
	if !LabelNameRegexp.Match([]byte(labelName)) {
		return fmt.Errorf("label key format error %v", labelName)
	}
	return nil
}

func ValidateLabelValue(labelValue string) error {
	if !govalidator.StringLength(labelValue, LabelValueMinLength, LabelValueMaxLength) {
		return fmt.Errorf("the length of the label value should be 1 to 255")
	}
	return nil
}

func validateCreateRuntimeRequest(req *pb.CreateRuntimeRequest) error {
	err := ValidateName(req.Name.GetValue())
	if err != nil {
		return err
	}
	err = ValidateLabelString(req.Labels.GetValue())
	if err != nil {
		return err
	}
	err = ValidateURL(req.RuntimeUrl.GetValue())
	if err != nil {
		return err
	}
	return nil
}

func validateModifyRuntimeRequest(req *pb.ModifyRuntimeRequest) error {
	err := ValidateName(req.Name.GetValue())
	if err != nil {
		return err
	}
	err = ValidateLabelString(req.Labels.GetValue())
	if err != nil {
		return err
	}
	return nil
}

func validateDeleteRuntimeRequest(req *pb.DeleteRuntimeRequest) error {
	return nil
}

func validateDescribeRuntimeRequest(req *pb.DescribeRuntimesRequest) error {
	err := ValidateSelectorString(req.Selector.GetValue())
	if err != nil {
		return err
	}
	return nil
}
