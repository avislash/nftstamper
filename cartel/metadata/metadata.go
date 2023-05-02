package metadata

import (
	"encoding/json"
)

type HoundMetadata struct {
	Name       string
	Image      string
	Background string
	Face       string
	Form       string
	Mouth      string
	Torso      string
}

func (hmd *HoundMetadata) UnmarshalJSON(data []byte) error {
	s := struct {
		Name       string            `json:"name"`
		Image      string            `json:"image"`
		Attributes []json.RawMessage `json:"attributes"`
	}{}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	for _, _attribute := range s.Attributes {
		attribute := make(map[string]string)
		err := json.Unmarshal(_attribute, &attribute)
		if err != nil {
			return err
		}

		if attribute["trait_type"] == "Background" {
			hmd.Background = attribute["value"]
			continue
		}

		if attribute["trait_type"] == "Face" {
			hmd.Face = attribute["value"]
		}

		if attribute["trait_type"] == "Forms" {
			hmd.Form = attribute["value"]
		}

		if attribute["trait_type"] == "Mouth" {
			hmd.Mouth = attribute["value"]
		}

		if attribute["trait_type"] == "Torso" {
			hmd.Torso = attribute["value"]
		}
	}
	hmd.Name = s.Name
	hmd.Image = s.Image

	return nil
}
