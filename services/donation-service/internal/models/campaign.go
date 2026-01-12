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
	CollectedAmount float64            `bson:"collected_amount" json:"collected_amount"`
	Currency        string             `bson:"currency" json:"currency"`
	Status          CampaignStatus     `bson:"status" json:"status"`
	AllowAnonymous  bool               `bson:"allow_anonymous" json:"allow_anonymous"`
	AllowRecurring  bool               `bson:"allow_recurring" json:"allow_recurring"`
	Tags            []string           `bson:"tags" json:"tags"`
	
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
