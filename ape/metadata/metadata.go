package metadata

import (
	"encoding/json"
	"strings"
)

type SentinelMetadata struct {
	Name  string
	Image string
	Attributes
}
type Attributes struct {
	Background string
	BaseArmor  string
	Body       string
	Face       string
	Head       string
}

func (smd *SentinelMetadata) UnmarshalJSON(data []byte) error {
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
			smd.Background = strings.ToLower(attribute["value"])
			continue
		}

		if attribute["trait_type"] == "Base Armor" {
			smd.BaseArmor = strings.ToLower(attribute["value"])
		}

		if attribute["trait_type"] == "Body" {
			smd.Body = strings.ToLower(attribute["value"])
		}

		if attribute["trait_type"] == "Face" {
			smd.Face = strings.ToLower(attribute["value"])
		}

		if attribute["trait_type"] == "Head" {
			smd.Head = strings.ToLower(attribute["value"])
		}

	}
	smd.Name = s.Name
	smd.Image = s.Image

	return nil

}
