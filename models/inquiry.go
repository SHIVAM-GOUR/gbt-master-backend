package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"gorm.io/gorm"
)

type Inquiry struct {
	ID                     *int16     `gorm:"column:id;primaryKey;autoIncrement" json:"id" form:"-"`
	Name                   *string    `gorm:"column:name;size:100;not null" json:"name,omitempty" form:"name" binding:"required"`
	Email                  *string    `gorm:"column:email;size:255" json:"email,omitempty" form:"email"`
	Phone                  *string    `gorm:"column:phone;size:20" json:"phone,omitempty" form:"phone"`
	Industry               *string    `gorm:"column:industry;size:40" json:"industry,omitempty" form:"industry"`
	HasExistingWebsite     *bool      `gorm:"column:has_existing_website" json:"has_existing_website,omitempty" form:"has_existing_website"`
	PreferredContactMethod *string    `gorm:"column:preferred_contact_method;size:40" json:"preferred_contact_method,omitempty" form:"preferred_contact_method"`
	Status                 *string    `gorm:"column:status;size:40;default:'new'" json:"status,omitempty" form:"status"`
	CreatedAt              *time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at" form:"-"`
	UpdatedAt              *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at" form:"-"`
}

func (i *Inquiry) TableName() string {
	return "inquiry"
}

func (i *Inquiry) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.Omits = i.getOmittedList()
	return nil
}

func (i *Inquiry) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.Omits = i.getOmittedList()
	return nil
}

func (i *Inquiry) getOmittedList() []string {
	omittedList := []string{}
	if i.Name == nil {
		omittedList = append(omittedList, "name")
	}
	if i.Email == nil {
		omittedList = append(omittedList, "email")
	}
	if i.Phone == nil {
		omittedList = append(omittedList, "phone")
	}
	if i.Industry == nil {
		omittedList = append(omittedList, "industry")
	}
	if i.HasExistingWebsite == nil {
		omittedList = append(omittedList, "has_existing_website")
	}
	if i.PreferredContactMethod == nil {
		omittedList = append(omittedList, "preferred_contact_method")
	}
	if i.Status == nil {
		omittedList = append(omittedList, "status")
	}
	return omittedList
}

func (i *Inquiry) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(i)
	if err != nil {
		return err
	}

	if i.Name == nil {
		return errors.New("invalid request body")
	}

	if *i.Name == "" {
		return errors.New("invalid request body")
	}

	return nil
}
