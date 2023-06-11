package metadata

import (
	"encoding/json"
	"strings"
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
			hmd.Background = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Face" {
			hmd.Face = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Forms" {
			hmd.Form = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Mouth" {
			hmd.Mouth = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Torso" {
			hmd.Torso = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}
	}
	hmd.Name = s.Name
	hmd.Image = s.Image

	return nil
}

type MAYCMetadata struct {
	Name       string
	Image      string
	Background string
	Mouth      string
	Clothes    string
	Earring    string
	Eyes       string
	Fur        string
	Hat        string
}

func (mm *MAYCMetadata) UnmarshalJSON(data []byte) error {
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
			mm.Background = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Mouth" {
			mm.Mouth = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Clothes" {
			mm.Clothes = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Earring" {
			mm.Earring = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Eyes" {
			mm.Eyes = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Fur" {
			mm.Fur = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Hat" {
			mm.Hat = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}
	}
	mm.Name = s.Name
	mm.Image = s.Image

	return nil
}
