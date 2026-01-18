/// Action approval model for multi-admin workflow
class ActionApproval {
  final String id;
  final String enterpriseId;
  final String actionType;
  final String actionName;
  final String description;
  final double? amount;
  final String? currency;
  final String status;
  final String initiatorUserId;
  final String? initiatorName;
  final int requiredApprovals;
  final List<ApprovalVote> approvals;
  final Map<String, dynamic> payload;
  final DateTime? expiresAt;
  final DateTime? createdAt;

  ActionApproval({
    required this.id,
    required this.enterpriseId,
    required this.actionType,
    required this.actionName,
    required this.description,
    this.amount,
    this.currency,
    required this.status,
    required this.initiatorUserId,
    this.initiatorName,
    this.requiredApprovals = 1,
    this.approvals = const [],
    this.payload = const {},
    this.expiresAt,
    this.createdAt,
  });

  factory ActionApproval.fromJson(Map<String, dynamic> json) {
    return ActionApproval(
      id: json['id'] ?? json['_id'] ?? '',
      enterpriseId: json['enterprise_id'] ?? '',
      actionType: json['action_type'] ?? '',
      actionName: json['action_name'] ?? '',
      description: json['description'] ?? '',
      amount: json['amount']?.toDouble(),
      currency: json['currency'],
      status: json['status'] ?? 'PENDING',
      initiatorUserId: json['initiator_user_id'] ?? '',
      initiatorName: json['initiator_name'],
      requiredApprovals: json['required_approvals'] ?? 1,
      approvals: (json['approvals'] as List<dynamic>?)
          ?.map((e) => ApprovalVote.fromJson(e))
          .toList() ?? [],
      payload: Map<String, dynamic>.from(json['payload'] ?? {}),
      expiresAt: json['expires_at'] != null 
          ? DateTime.tryParse(json['expires_at']) 
          : null,
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at']) 
          : null,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'enterprise_id': enterpriseId,
    'action_type': actionType,
    'action_name': actionName,
    'description': description,
    'amount': amount,
    'currency': currency,
    'status': status,
    'initiator_user_id': initiatorUserId,
    'initiator_name': initiatorName,
    'required_approvals': requiredApprovals,
    'payload': payload,
  };

  int get approvedCount => approvals.where((a) => a.decision == 'APPROVED').length;
  int get rejectedCount => approvals.where((a) => a.decision == 'REJECTED').length;
  double get progress => requiredApprovals > 0 ? approvedCount / requiredApprovals : 0;
  
  bool get isPending => status == ApprovalStatus.pending;
  bool get isApproved => status == ApprovalStatus.approved;
  bool get isRejected => status == ApprovalStatus.rejected;
  bool get isExpired => status == ApprovalStatus.expired;
  
  bool hasUserVoted(String userId) {
    return approvals.any((a) => a.adminUserId == userId);
  }
  
  String get statusLabel {
    switch (status) {
      case ApprovalStatus.pending: return 'En attente';
      case ApprovalStatus.approved: return 'Approuvé';
      case ApprovalStatus.rejected: return 'Rejeté';
      case ApprovalStatus.executed: return 'Exécuté';
      case ApprovalStatus.expired: return 'Expiré';
      default: return status;
    }
  }
}

/// Individual vote on an approval
class ApprovalVote {
  final String adminUserId;
  final String? adminName;
  final String decision;
  final String? reason;
  final DateTime? votedAt;

  ApprovalVote({
    required this.adminUserId,
    this.adminName,
    required this.decision,
    this.reason,
    this.votedAt,
  });

  factory ApprovalVote.fromJson(Map<String, dynamic> json) {
    return ApprovalVote(
      adminUserId: json['admin_user_id'] ?? '',
      adminName: json['admin_name'],
      decision: json['decision'] ?? '',
      reason: json['reason'],
      votedAt: json['voted_at'] != null 
          ? DateTime.tryParse(json['voted_at']) 
          : null,
    );
  }

  bool get isApproved => decision == 'APPROVED';
  bool get isRejected => decision == 'REJECTED';
}

/// Approval status constants
class ApprovalStatus {
  static const String pending = 'PENDING';
  static const String approved = 'APPROVED';
  static const String rejected = 'REJECTED';
  static const String executed = 'EXECUTED';
  static const String expired = 'EXPIRED';
}

/// Action types
class ActionType {
  static const String transaction = 'TRANSACTION';
  static const String payroll = 'PAYROLL';
  static const String employeeUpdate = 'EMPLOYEE_UPDATE';
  static const String settings = 'SETTINGS';
}
