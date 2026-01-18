/// Employee model representing a member of an enterprise
class Employee {
  final String id;
  final String enterpriseId;
  final String userId;
  final String role;
  final String status;
  final String? firstName;
  final String? lastName;
  final String? email;
  final String? phone;
  final String? position;
  final String? department;
  final double? salary;
  final String? salaryCurrency;
  final String? salaryFrequency;
  final List<String> permissions;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Employee({
    required this.id,
    required this.enterpriseId,
    required this.userId,
    required this.role,
    required this.status,
    this.firstName,
    this.lastName,
    this.email,
    this.phone,
    this.position,
    this.department,
    this.salary,
    this.salaryCurrency,
    this.salaryFrequency,
    this.permissions = const [],
    this.createdAt,
    this.updatedAt,
  });

  factory Employee.fromJson(Map<String, dynamic> json) {
    return Employee(
      id: json['id'] ?? json['_id'] ?? '',
      enterpriseId: json['enterprise_id'] ?? '',
      userId: json['user_id'] ?? '',
      role: json['role'] ?? 'EMPLOYEE',
      status: json['status'] ?? 'PENDING',
      firstName: json['first_name'],
      lastName: json['last_name'],
      email: json['email'],
      phone: json['phone_number'],
      position: json['job_title'] ?? json['position'],
      department: json['department'],
      salary: json['salary']?.toDouble(),
      salaryCurrency: json['salary_currency'],
      salaryFrequency: json['salary_frequency'],
      permissions: List<String>.from(json['permissions'] ?? []),
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
    'enterprise_id': enterpriseId,
    'user_id': userId,
    'role': role,
    'status': status,
    'first_name': firstName,
    'last_name': lastName,
    'email': email,
    'phone_number': phone,
    'job_title': position,
    'department': department,
    'salary': salary,
    'salary_currency': salaryCurrency,
    'salary_frequency': salaryFrequency,
    'permissions': permissions,
  };

  String get fullName => '$firstName $lastName'.trim();
  
  bool get isAdmin => role == EmployeeRole.admin || role == EmployeeRole.owner;
  bool get isOwner => role == EmployeeRole.owner;
  bool get isActive => status == EmployeeStatus.active;
  bool get isPending => status == EmployeeStatus.pending || status == EmployeeStatus.pendingInvite;
  
  String get statusLabel {
    switch (status) {
      case EmployeeStatus.active: return 'Actif';
      case EmployeeStatus.pending: return 'En attente';
      case EmployeeStatus.pendingInvite: return 'Invitation envoyée';
      case EmployeeStatus.suspended: return 'Suspendu';
      case EmployeeStatus.terminated: return 'Terminé';
      default: return status;
    }
  }
  
  String get roleLabel {
    switch (role) {
      case EmployeeRole.owner: return 'Propriétaire';
      case EmployeeRole.admin: return 'Administrateur';
      case EmployeeRole.employee: return 'Employé';
      default: return role;
    }
  }
}

/// Employee roles
class EmployeeRole {
  static const String owner = 'OWNER';
  static const String admin = 'ADMIN';
  static const String employee = 'EMPLOYEE';
}

/// Employee status
class EmployeeStatus {
  static const String active = 'ACTIVE';
  static const String pending = 'PENDING';
  static const String pendingInvite = 'PENDING_INVITE';
  static const String suspended = 'SUSPENDED';
  static const String terminated = 'TERMINATED';
}
