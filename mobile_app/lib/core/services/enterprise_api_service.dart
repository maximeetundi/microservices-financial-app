import 'package:dio/dio.dart';
import 'base_api_service.dart';

class EnterpriseApiService extends BaseApiService {
  static const String _basePath = '/enterprise-service/api/v1';
  
  // ==================== ENTERPRISE ====================
  
  /// Get all enterprises the user has access to
  Future<dynamic> getEnterprises() async {
    return get('$_basePath/enterprises');
  }
  
  /// Get single enterprise by ID
  Future<dynamic> getEnterprise(String id) async {
    return get('$_basePath/enterprises/$id');
  }
  
  /// Create new enterprise
  Future<dynamic> createEnterprise(Map<String, dynamic> data) async {
    return post('$_basePath/enterprises', data: data);
  }
  
  /// Update enterprise
  Future<dynamic> updateEnterprise(String id, Map<String, dynamic> data) async {
    return put('$_basePath/enterprises/$id', data: data);
  }
  
  // ==================== EMPLOYEES ====================
  
  /// Get all employees for an enterprise
  Future<dynamic> getEmployees(String enterpriseId) async {
    return get('$_basePath/enterprises/$enterpriseId/employees');
  }
  
  /// Get current user's employee record
  Future<dynamic> getMyEmployee(String enterpriseId) async {
    return get('$_basePath/enterprises/$enterpriseId/employees/me');
  }
  
  /// Invite new employee
  Future<dynamic> inviteEmployee(String enterpriseId, Map<String, dynamic> data) async {
    return post('$_basePath/enterprises/$enterpriseId/employees/invite', data: data);
  }
  
  /// Update employee (role, salary, etc.)
  Future<dynamic> updateEmployee(String enterpriseId, String employeeId, Map<String, dynamic> data) async {
    return put('$_basePath/enterprises/$enterpriseId/employees/$employeeId', data: data);
  }
  
  /// Accept employee invitation
  Future<dynamic> acceptInvitation(String pin, String token) async {
    return post('$_basePath/employees/accept', data: {'pin': pin, 'token': token});
  }
  
  // ==================== WALLETS ====================
  
  /// Get enterprise wallets (uses wallet service)
  Future<dynamic> getEnterpriseWallets() async {
    return get('/wallet-service/api/v1/wallets');
  }
  
  /// Create enterprise wallet
  Future<dynamic> createEnterpriseWallet(Map<String, dynamic> data) async {
    return post('/wallet-service/api/v1/wallets', data: data);
  }
  
  /// Get wallet transactions
  Future<dynamic> getWalletTransactions(String walletId) async {
    return get('/wallet-service/api/v1/wallets/$walletId/transactions');
  }
  
  // ==================== MULTI-ADMIN APPROVAL ====================
  
  /// Get pending approvals for an enterprise
  Future<dynamic> getPendingApprovals(String enterpriseId) async {
    return get('$_basePath/enterprises/$enterpriseId/approvals');
  }
  
  /// Get single approval by ID
  Future<dynamic> getApprovalById(String approvalId) async {
    return get('$_basePath/approvals/$approvalId');
  }
  
  /// Initiate action (creates approval request)
  Future<dynamic> initiateAction(String enterpriseId, Map<String, dynamic> data) async {
    return post('$_basePath/enterprises/$enterpriseId/actions', data: data);
  }
  
  /// Approve action (requires PIN)
  Future<dynamic> approveAction(String approvalId, String encryptedPin) async {
    return post('$_basePath/approvals/$approvalId/approve', data: {'pin': encryptedPin});
  }
  
  /// Reject action (requires PIN)
  Future<dynamic> rejectAction(String approvalId, String encryptedPin, {String? reason}) async {
    return post('$_basePath/approvals/$approvalId/reject', data: {
      'pin': encryptedPin,
      if (reason != null) 'reason': reason,
    });
  }
  
  // ==================== SERVICES ====================
  
  /// Get enterprise services (from enterprise.service_groups)
  Future<dynamic> getServices(String enterpriseId) async {
    final enterprise = await getEnterprise(enterpriseId);
    return enterprise['service_groups'] ?? [];
  }
  
  /// Update enterprise services
  Future<dynamic> updateServices(String enterpriseId, List<Map<String, dynamic>> serviceGroups) async {
    return put('$_basePath/enterprises/$enterpriseId', data: {
      'service_groups': serviceGroups,
    });
  }
  
  // ==================== PAYROLL ====================
  
  /// Get payroll preview
  Future<dynamic> getPayrollPreview(String enterpriseId) async {
    return post('$_basePath/enterprises/$enterpriseId/payroll/preview');
  }
  
  /// Run payroll
  Future<dynamic> runPayroll(String enterpriseId, Map<String, dynamic> data) async {
    return post('$_basePath/enterprises/$enterpriseId/payroll/run', data: data);
  }
  
  /// Get payroll history
  Future<dynamic> getPayrollHistory(String enterpriseId) async {
    return get('$_basePath/enterprises/$enterpriseId/payroll/history');
  }
  
  // ==================== QR CODES ====================
  
  /// Generate enterprise QR code
  Future<dynamic> getEnterpriseQR(String enterpriseId) async {
    return get('$_basePath/enterprises/$enterpriseId/qrcode');
  }
  
  /// Generate service QR code
  Future<dynamic> getServiceQR(String enterpriseId, String serviceId) async {
    return get('$_basePath/enterprises/$enterpriseId/services/$serviceId/qrcode');
  }
  
  // ==================== BILLING ====================
  
  /// Create invoice
  Future<dynamic> createInvoice(Map<String, dynamic> data) async {
    return post('$_basePath/invoices', data: data);
  }
  
  /// Get enterprise invoices
  Future<dynamic> getInvoices(String enterpriseId, {String? status}) async {
    final query = status != null ? '?status=$status' : '';
    return get('$_basePath/enterprises/$enterpriseId/invoices$query');
  }
}
