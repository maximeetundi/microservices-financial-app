import 'package:dio/dio.dart';
import 'base_api_service.dart';

class EnterpriseApiService extends BaseApiService {
  // Endpoints
  static const String _basePath = '/enterprise-service/api/v1';
  
  Future<dynamic> createEnterprise(Map<String, dynamic> data) async {
    return post('$_basePath/enterprises', data: data);
  }

  Future<dynamic> getEnterprise(String id) async {
    return get('$_basePath/enterprises/$id');
  }

  Future<dynamic> getEnterprises() async {
    return get('$_basePath/enterprises');
  }

  // Employee
  Future<dynamic> inviteEmployee(Map<String, dynamic> data) async {
    return post('$_basePath/employees/invite', data: data);
  }

  Future<dynamic> acceptInvitation(String pin, String token) async {
    return post('$_basePath/employees/accept', data: {'pin': pin, 'token': token});
  }

  // Payroll
  Future<dynamic> getPayrollPreview(String entId) async {
    return post('$_basePath/enterprises/$entId/payroll/preview');
  }
}
