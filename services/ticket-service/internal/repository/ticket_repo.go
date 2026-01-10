package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TicketRepository struct {
	collection *mongo.Collection
	eventColl  *mongo.Collection
	tierColl   *mongo.Collection
}

func NewTicketRepository(db *mongo.Database) *TicketRepository {
	return &TicketRepository{
		collection: db.Collection("tickets"),
		eventColl:  db.Collection("events"),
		tierColl:   db.Collection("ticket_tiers"),
	}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	if ticket.ID == "" {
		ticket.ID = uuid.New().String()
	}
	ticket.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), ticket)
	return err
}

func (r *TicketRepository) GetByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.Background(), filter).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	
	// Populate Event details
	r.populateEventDetails(&ticket)
	return &ticket, nil
}

func (r *TicketRepository) GetByCode(code string) (*models.Ticket, error) {
	var ticket models.Ticket
	filter := bson.M{"ticket_code": code}
	err := r.collection.FindOne(context.Background(), filter).Decode(&ticket)
	if err != nil {
		return nil, err
	}

	r.populateEventDetails(&ticket)
	return &ticket, nil
}

func (r *TicketRepository) populateEventDetails(ticket *models.Ticket) {
	var event models.Event
	err := r.eventColl.FindOne(context.Background(), bson.M{"_id": ticket.EventID}).Decode(&event)
	if err == nil {
		ticket.EventTitle = event.Title
		ticket.EventDate = &event.StartDate
		ticket.EventLocation = event.Location
		ticket.EventCreatorID = event.CreatorID
	}
}

func (r *TicketRepository) GetByBuyer(buyerID string, limit, offset int) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	filter := bson.M{"buyer_id": buyerID, "status": "paid"}
	
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var ticket models.Ticket
		if err := cursor.Decode(&ticket); err != nil {
			return nil, err
		}
		// Optimization: Could use aggregation $lookup, but doing separate queries for now for simplicity
		// or maybe not populate event details for list view if not strictly needed, 
		// but frontend likely expects it.
		r.populateEventDetails(&ticket)
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

func (r *TicketRepository) GetByTransactionID(transactionID string) (*models.Ticket, error) {
	var ticket models.Ticket
	filter := bson.M{"transaction_id": transactionID}
	err := r.collection.FindOne(context.Background(), filter).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	r.populateEventDetails(&ticket)
	return &ticket, nil
}

func (r *TicketRepository) GetByEvent(eventID string, limit, offset int) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	filter := bson.M{"event_id": eventID}
	
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var ticket models.Ticket
		if err := cursor.Decode(&ticket); err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

func (r *TicketRepository) UpdateStatus(id, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *TicketRepository) UpdateTransactionID(id, transactionID string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"transaction_id": transactionID}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *TicketRepository) UpdateStatusByTransactionID(transactionID, status string) error {
	filter := bson.M{"transaction_id": transactionID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *TicketRepository) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *TicketRepository) MarkAsUsed(id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":  "used",
			"used_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *TicketRepository) GetEventStats(eventID string) (*models.TicketStats, error) {
	stats := &models.TicketStats{}

	// Aggregation for overall stats
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"event_id": eventID}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"total_tickets": bson.M{"$sum": 1},
			"sold_tickets": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$in": bson.A{"$status", bson.A{"paid", "used"}}},
						1,
						0,
					},
				},
			},
			"used_tickets": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$status", "used"}},
						1,
						0,
					},
				},
			},
			"total_revenue": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$in": bson.A{"$status", bson.A{"paid", "used"}}},
						"$price",
						0,
					},
				},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		var result struct {
			TotalTickets int     `bson:"total_tickets"`
			SoldTickets  int     `bson:"sold_tickets"`
			UsedTickets  int     `bson:"used_tickets"`
			TotalRevenue float64 `bson:"total_revenue"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.TotalTickets = result.TotalTickets
			stats.SoldTickets = result.SoldTickets
			stats.UsedTickets = result.UsedTickets
			stats.TotalRevenue = result.TotalRevenue
		}
	}

	// Get tier breakdown
	// We need to fetch tiers from Tiers collection and map sold/revenue
	tiersCursor, err := r.tierColl.Find(context.Background(), 
		bson.M{"event_id": eventID},
		options.Find().SetSort(bson.M{"sort_order": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer tiersCursor.Close(context.Background())

	for tiersCursor.Next(context.Background()) {
		var tier models.TicketTier
		if err := tiersCursor.Decode(&tier); err != nil {
			continue
		}
		
		// Ensure we use the actual sold/revenue from tickets collection logic if needed, 
		// but tier object also has 'sold' field maintained.
		// For consistency, let's trust the tier.Sold field for count, but calculate revenue ourselves 
		// or just use tier.Price * tier.Sold
		
		tierStats := models.TierStats{
			TierID:   tier.ID,
			TierName: tier.Name,
			TierIcon: tier.Icon,
			Price:    tier.Price,
			Quantity: tier.Quantity,
			Sold:     tier.Sold,
			Used:     0, // Need to query ticket collection for this specific tier used count?
			Revenue:  float64(tier.Sold) * tier.Price,
		}
		
		// Optional: querying used count for this tier
		usedCount, _ := r.collection.CountDocuments(context.Background(), bson.M{
			"event_id": eventID,
			"tier_id": tier.ID,
			"status": "used",
		})
		tierStats.Used = int(usedCount)

		stats.TierBreakdown = append(stats.TierBreakdown, tierStats)
	}

	return stats, nil
}
