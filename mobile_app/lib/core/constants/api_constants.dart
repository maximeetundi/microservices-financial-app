/// API Constants for the mobile application
class ApiConstants {
  // Production API Base URL (Kong API Gateway)
  static const String baseUrl = 'https://api.app.maximeetundi.store';
  
  // Development API Base URL (local)
  static const String devBaseUrl = 'http://localhost:8080';
  
  // Timeout durations
  static const int connectionTimeout = 30000; // 30 seconds
  static const int receiveTimeout = 30000; // 30 seconds
  
  // WebSocket URL
  static const String wsBaseUrl = 'wss://api.app.maximeetundi.store';
  static const String devWsBaseUrl = 'ws://localhost:8080';
}
