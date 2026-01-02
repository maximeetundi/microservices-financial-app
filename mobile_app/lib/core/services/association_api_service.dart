import 'package:dio/dio.dart';
import '../constants/api_constants.dart';

class AssociationApiService {
  final Dio _dio;

  AssociationApiService(this._dio);

  // ========== Associations ==========
  
  Future<Response> createAssociation(Map<String, dynamic> data) async {
    return await _dio.post('${ApiConstants.baseUrl}/association-service/api/v1/associations', data: data);
  }

  Future<Response> getAssociations() async {
    return await _dio.get('${ApiConstants.baseUrl}/association-service/api/v1/associations');
  }

  Future<Response> getAssociation(String id) async {
    return await _dio.get('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id');
  }

  Future<Response> updateAssociation(String id, Map<String, dynamic> data) async {
    return await _dio.put('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id', data: data);
  }

  Future<Response> deleteAssociation(String id) async {
    return await _dio.delete('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id');
  }

  // ========== Members ==========

  Future<Response> joinAssociation(String id, String message) async {
    return await _dio.post('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/join', 
      data: {'message': message});
  }

  Future<Response> leaveAssociation(String id) async {
    return await _dio.post('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/leave');
  }

  Future<Response> getMembers(String id) async {
    return await _dio.get('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/members');
  }

  Future<Response> updateMemberRole(String associationId, String userId, String role) async {
    return await _dio.put(
      '${ApiConstants.baseUrl}/association-service/api/v1/associations/$associationId/members/$userId/role',
      data: {'role': role}
    );
  }

  // ========== Meetings ==========

  Future<Response> createMeeting(String id, Map<String, dynamic> data) async {
    return await _dio.post('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/meetings', data: data);
  }

  Future<Response> getMeetings(String id) async {
    return await _dio.get('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/meetings');
  }

  Future<Response> recordAttendance(String meetingId, Map<String, bool> attendance) async {
    return await _dio.post(
      '${ApiConstants.baseUrl}/association-service/api/v1/meetings/$meetingId/attendance',
      data: {'attendance': attendance}
    );
  }

  // ========== Treasury ==========

  Future<Response> recordContribution(String id, Map<String, dynamic> data) async {
    return await _dio.post('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/contributions', data: data);
  }

  Future<Response> getTreasury(String id) async {
    return await _dio.get('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/treasury');
  }

  // ========== Loans ==========

  Future<Response> requestLoan(String id, Map<String, dynamic> data) async {
    return await _dio.post('${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/loans', data: data);
  }

  Future<Response> approveLoan(String loanId, bool approve, String reason) async {
    return await _dio.put(
      '${ApiConstants.baseUrl}/association-service/api/v1/loans/$loanId/approve',
      data: {'approve': approve, 'reason': reason}
    );
  }

  Future<Response> repayLoan(String loanId, double amount) async {
    return await _dio.post(
      '${ApiConstants.baseUrl}/association-service/api/v1/loans/$loanId/repay',
      data: {'amount': amount}
    );
  }

  Future<Response> distributeFunds(String id, double amount, List<String> memberIds) async {
    return await _dio.post(
      '${ApiConstants.baseUrl}/association-service/api/v1/associations/$id/distribute',
      data: {'amount': amount, 'member_ids': memberIds}
    );
  }
}
