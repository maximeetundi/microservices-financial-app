package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

type Campaign struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatorID       string             `bson:"creator_id" json:"creator_id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	ImageURL        string             `bson:"image_url" json:"image_url"`
	VideoURL        string             `bson:"video_url" json:"video_url"`
	TargetAmount    float64            `bson:"target_amount" json:"target_amount"`       // 0 for open ended
	MinAmount       float64            `bson:"min_amount" json:"min_amount"`             // Minimum donation amount
	MaxAmount       float64            `bson:"max_amount" json:"max_amount"`             // Maximum donation amount
	CollectedAmount float64            `bson:"collected_amount" json:"collected_amount"`
	Currency        string             `bson:"currency" json:"currency"`
	Status          CampaignStatus     `bson:"status" json:"status"`
	AllowAnonymous  bool               `bson:"allow_anonymous" json:"allow_anonymous"`
	AllowRecurring  bool               `bson:"allow_recurring" json:"allow_recurring"`
	Tags            []string           `bson:"tags" json:"tags"`
	ShowDonors      bool               `bson:"show_donors" json:"show_donors"`
	
	QRCode          string             `bson:"qr_code,omitempty" json:"qr_code,omitempty"`
	CampaignCode    string             `bson:"campaign_code,omitempty" json:"campaign_code,omitempty"`
	
	// FormSchema for dynamic data collection from donors
	FormSchema      []FormField `bson:"form_schema,omitempty" json:"form_schema,omitempty"`

	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}

type FormField struct {
	Name     string   `bson:"name" json:"name"`
	Label    string   `bson:"label" json:"label"`
	Type     string   `bson:"type" json:"type"`
	Required bool     `bson:"required" json:"required"`
	Options  []string `bson:"options,omitempty" json:"options,omitempty"`
}
