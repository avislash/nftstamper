package metadata

import (
	"encoding/json"
	"strings"
)

type HoundMetadata struct {
	Name       string
	Image      string
	Background string
	Head       string
	Face       string
	Nose       string
	Mouth      string
	Form       string
	Leg        string
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

		if attribute["trait_type"] == "Head" {
			hmd.Head = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Face" {
			hmd.Face = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Nose" {
			hmd.Nose = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Mouth" {
			hmd.Mouth = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Forms" {
			hmd.Form = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Leg" {
			hmd.Leg = strings.TrimSpace(strings.ToLower(attribute["value"]))
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

		if attribute["trait_type"] == "Name" {
			mm.Name = strings.TrimSpace(strings.ToLower(attribute["value"]))
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
	mm.Image = s.Image

	return nil
}

type BAYCMetadata struct {
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

func (bm *BAYCMetadata) UnmarshalJSON(data []byte) error {
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

		if attribute["trait_type"] == "Name" {
			bm.Name = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Background" {
			bm.Background = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Mouth" {
			bm.Mouth = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Clothes" {
			bm.Clothes = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Earring" {
			bm.Earring = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Eyes" {
			bm.Eyes = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Fur" {
			bm.Fur = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}

		if attribute["trait_type"] == "Hat" {
			bm.Hat = strings.TrimSpace(strings.ToLower(attribute["value"]))
		}
	}
	bm.Image = s.Image

	return nil
}
