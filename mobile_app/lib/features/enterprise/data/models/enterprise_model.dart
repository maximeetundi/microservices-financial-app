/// Enterprise model representing a business entity
class Enterprise {
  final String id;
  final String name;
  final String? description;
  final String type;
  final String status;
  final String ownerId;
  final String? logo;
  final String? defaultWalletId;
  final List<String> walletIds;
  final List<ServiceGroup> serviceGroups;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Enterprise({
    required this.id,
    required this.name,
    this.description,
    required this.type,
    required this.status,
    required this.ownerId,
    this.logo,
    this.defaultWalletId,
    this.walletIds = const [],
    this.serviceGroups = const [],
    this.createdAt,
    this.updatedAt,
  });

  factory Enterprise.fromJson(Map<String, dynamic> json) {
    return Enterprise(
      id: json['id'] ?? json['_id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      type: json['type'] ?? 'SERVICE',
      status: json['status'] ?? 'ACTIVE',
      ownerId: json['owner_id'] ?? '',
      logo: json['logo'],
      defaultWalletId: json['default_wallet_id'],
      walletIds: List<String>.from(json['wallet_ids'] ?? []),
      serviceGroups: (json['service_groups'] as List<dynamic>?)
          ?.map((e) => ServiceGroup.fromJson(e))
          .toList() ?? [],
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at']) 
          : null,
      updatedAt: json['updated_at'] != null 
          ? DateTime.tryParse(json['updated_at']) 
          : null,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'description': description,
    'type': type,
    'status': status,
    'owner_id': ownerId,
    'logo': logo,
    'default_wallet_id': defaultWalletId,
    'wallet_ids': walletIds,
    'service_groups': serviceGroups.map((e) => e.toJson()).toList(),
  };

  bool get isOwner => true; // Will be set based on current user
}

/// Enterprise types
class EnterpriseType {
  static const String service = 'SERVICE';
  static const String school = 'SCHOOL';
  static const String transport = 'TRANSPORT';
  static const String utility = 'UTILITY';
}

/// Service group within an enterprise
class ServiceGroup {
  final String id;
  final String name;
  final String? walletId;
  final String currency;
  final bool isPrivate;
  final List<Service> services;

  ServiceGroup({
    required this.id,
    required this.name,
    this.walletId,
    required this.currency,
    this.isPrivate = false,
    this.services = const [],
  });

  factory ServiceGroup.fromJson(Map<String, dynamic> json) {
    return ServiceGroup(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      walletId: json['wallet_id'],
      currency: json['currency'] ?? 'XOF',
      isPrivate: json['is_private'] ?? false,
      services: (json['services'] as List<dynamic>?)
          ?.map((e) => Service.fromJson(e))
          .toList() ?? [],
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'wallet_id': walletId,
    'currency': currency,
    'is_private': isPrivate,
    'services': services.map((e) => e.toJson()).toList(),
  };
}

/// Individual service offered by enterprise
class Service {
  final String id;
  final String name;
  final String? description;
  final double price;
  final String? reference;

  Service({
    required this.id,
    required this.name,
    this.description,
    required this.price,
    this.reference,
  });

  factory Service.fromJson(Map<String, dynamic> json) {
    return Service(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      price: (json['price'] ?? 0).toDouble(),
      reference: json['reference'],
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'description': description,
    'price': price,
    'reference': reference,
  };
}
