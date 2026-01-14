
class Subscription {
  final String id;
  final String enterpriseId;
  final String enterpriseName;
  final String serviceId;
  final String serviceName;
  final String status;
  final double amount;
  final String billingFrequency;
  final DateTime nextBillingAt;
  final String externalId;
  final String clientName;
  
  // School Details (Optional)
  final String? studentName;
  final String? className;

  Subscription({
    required this.id,
    required this.enterpriseId,
    required this.enterpriseName,
    required this.serviceId,
    required this.serviceName,
    required this.status,
    required this.amount,
    required this.billingFrequency,
    required this.nextBillingAt,
    required this.externalId,
    required this.clientName,
    this.studentName,
    this.className,
  });

  factory Subscription.fromJson(Map<String, dynamic> json) {
    return Subscription(
      id: json['id'] ?? '',
      enterpriseId: json['enterprise_id'] ?? '',
      enterpriseName: json['enterprise_name'] ?? 'Unknown Enterprise',
      serviceId: json['service_id'] ?? '',
      serviceName: json['service_name'] ?? 'Unknown Service',
      status: json['status'] ?? 'ACTIVE',
      amount: (json['amount'] ?? 0).toDouble(),
      billingFrequency: json['billing_frequency'] ?? 'MONTHLY',
      nextBillingAt: DateTime.tryParse(json['next_billing_at'] ?? '') ?? DateTime.now(),
      externalId: json['external_id'] ?? '',
      clientName: json['client_name'] ?? '',
      studentName: json['school_details']?['student_name'],
      className: json['school_details']?['class_id'],
    );
  }
}
