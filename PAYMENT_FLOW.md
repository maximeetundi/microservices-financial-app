# Modèle Standard d'Intégration des Paiements (Microservices)

Ce document décrit le modèle d'architecture à suivre pour tout nouveau microservice (ex: *Association Service*, *Donation Service*) souhaitant initier des paiements ou des transferts de fonds via le système centralisé.

## 1. Vue d'ensemble du Flux

L'architecture repose sur un modèle événementiel (Event-Driven) via Kafka. Le service initiateur (ex: `ticket-service`) ne communique jamais directement avec le `wallet-service` pour les débits. Il délègue cette responsabilité au `transfer-service`.

### Le Flux :
1. **Service Initiateur** (ex: Ticket, Association) : Valide la logique métier (ex: stock ticket disponible, PIN valide) et publie un événement `PaymentRequestEvent`.
2. **Transfer Service** : Consomme l'événement, gère la conversion de devises (via Exchange Service si nécessaire) et ordonne le débit/crédit via l'API interne du Wallet Service.
3. **Notification Service** : Informe l'utilisateur du succès ou de l'échec.
4. **Service Initiateur** : Écoute les événements de statut (`PaymentStatusEvent`) pour mettre à jour l'état de la commande (ex: passer le ticket de "PENDING" à "PAID" ou "CANCELLED").

## 2. Structure de l'Événement (PaymentRequestEvent)

Tout service initiateur doit publier un événement `messaging.EventPaymentRequest` ("payment.request") avec la structure suivante :

```go
type PaymentRequestEvent struct {
    RequestID         string  `json:"request_id"`          // UUID unique de la demande
    ReferenceID       string  `json:"reference_id"`        // ID de l'objet métier (ex: TicketID, DonationID)
    Type              string  `json:"type"`                // Type de flux (ex: "ticket_purchase", "donation")
    UserID            string  `json:"user_id"`             // ID de l'utilisateur payeur
    FromWalletID      string  `json:"from_wallet_id"`      // ID du wallet à débiter (Fourni par le frontend)
    DestinationUserID string  `json:"destination_user_id"` // ID de l'organisateur/bénéficiaire
    DebitAmount       float64 `json:"debit_amount"`        // Montant à débiter (Prix)
    Currency          string  `json:"currency"`            // Devise du montant
}
```

> **Note Importante:** Le champ `DebitAmount` est obligatoire. `Amount` n'est pas utilisé pour le traitement (car ambigu), mais `notification-service` a été mis à jour pour lire `DebitAmount` ou `Amount` par compatibilité.

## 3. Responsabilités du Service Initiateur

Si vous développez un nouveau service (ex: `association-service`), voici ce qu'il doit faire :

### A. Frontend (Vue.js)
1. Récupérer les wallets de l'utilisateur (`walletAPI.getWallets`).
2. Demander à l'utilisateur de choisir un wallet (`selectedWalletId`).
3. Valider le PIN via le composant global ou `usePin`.
4. Envoyer `wallet_id` et `pin` (crypté/validé) à votre API backend.

### B. Backend (Go)
1. **Validation** : Vérifier que `req.WalletID` est présent.
2. **Vérification PIN** : Appeler `userClient.VerifyPin(userID, pin, token)` (via `auth-service`).
3. **Création de l'objet** : Créer l'objet métier en base de données avec un statut `PENDING` (ex: `AssociationMembership { Status: "PENDING" }`).
4. **Récupération du Bénéficiaire** : S'assurer de récupérer l'ID du bénéficiaire (ex: `association.OwnerID`) pour remplir `DestinationUserID`.
5. **Publication Kafka** :
   ```go
   event := messaging.PaymentRequestEvent{
       RequestID:         uuid.New().String(),
       ReferenceID:       membership.ID, // Pour pouvoir le retrouver
       Type:              "association_membership",
       UserID:            userID,
       FromWalletID:      req.WalletID,
       DestinationUserID: association.OwnerID, // Le bénéficiaire (CRITIQUE)
       DebitAmount:       membership.Price,
       Currency:          membership.Currency,
   }
   kafkaClient.Publish(ctx, messaging.TopicPaymentEvents, event)
   ```
6. **Réponse HTTP** : Renvoyer 201 Created (ou 202 Accepted) au frontend immédiatement. Ne pas attendre le débit.

### C. Gestion des Retours (Consumer)
Le service doit aussi consommer `payment.status.success` et `payment.status.failed` pour mettre à jour l'objet.

```go
// Dans votre consumer.go
switch event.Type {
case messaging.EventPaymentSuccess:
    // Mettre à jour l'objet PENDING -> PAID
    s.repo.UpdateStatus(event.ReferenceID, "PAID")
case messaging.EventPaymentFailed:
    // Mettre à jour l'objet PENDING -> CANCELLED/FAILED
    s.repo.UpdateStatus(event.ReferenceID, "FAILED")
    // Optionnel : Libérer le stock
}
```

## 4. Ce qui est déjà géré par le Core (Ne pas refaire)
- **Conversion de devise** : Automatiquement géré par `transfer-service`. Si l'utilisateur paie en USD pour une association en XOF, le système convertit.
- **Notifications** : `notification-service` envoie automatiquement "Paiement validé" ou "Paiement échoué".
- **Logs et Traces** : Centralisés.

---
*Ce modèle garantit que tous les services financiers de l'application (Tickets, Dons, Tontines, Associations) fonctionnent de manière uniforme et sécurisée.*
