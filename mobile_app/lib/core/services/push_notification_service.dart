import 'dart:convert';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../api/api_client.dart';

/// Service for handling push notifications
class PushNotificationService {
  static final PushNotificationService _instance = PushNotificationService._internal();
  factory PushNotificationService() => _instance;
  PushNotificationService._internal();

  final FlutterLocalNotificationsPlugin _localNotifications = FlutterLocalNotificationsPlugin();
  final FlutterSecureStorage _storage = const FlutterSecureStorage();
  final ApiClient _apiClient = ApiClient();

  bool _isInitialized = false;
  
  // Callback when notification is tapped
  void Function(String? payload)? onNotificationTap;

  /// Initialize notification service
  Future<void> initialize() async {
    if (_isInitialized) return;

    // Skip on web - notifications not supported
    if (kIsWeb) {
      debugPrint('PushNotificationService: Web platform - skipping');
      _isInitialized = true;
      return;
    }

    // Initialize local notifications
    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );
    
    const initSettings = InitializationSettings(
      android: androidSettings,
      iOS: iosSettings,
    );

    await _localNotifications.initialize(
      initSettings,
      onDidReceiveNotificationResponse: _onNotificationTap,
    );

    // Request permissions on iOS
    if (Platform.isIOS) {
      await _localNotifications
          .resolvePlatformSpecificImplementation<IOSFlutterLocalNotificationsPlugin>()
          ?.requestPermissions(alert: true, badge: true, sound: true);
    }

    // Android permissions are requested via AndroidManifest.xml

    _isInitialized = true;
    debugPrint('PushNotificationService initialized');
  }

  /// Handle notification tap
  void _onNotificationTap(NotificationResponse response) {
    final payload = response.payload;
    debugPrint('Notification tapped with payload: $payload');
    onNotificationTap?.call(payload);
  }

  /// Show a local notification
  Future<void> showNotification({
    required int id,
    required String title,
    required String body,
    String? payload,
    String? channelId,
    String? channelName,
  }) async {
    final androidDetails = AndroidNotificationDetails(
      channelId ?? 'crypto_bank_channel',
      channelName ?? 'Zekora Notifications',
      channelDescription: 'Notifications from Zekora',
      importance: Importance.high,
      priority: Priority.high,
      showWhen: true,
      enableVibration: true,
      playSound: true,
      icon: '@mipmap/ic_launcher',
      largeIcon: const DrawableResourceAndroidBitmap('@mipmap/ic_launcher'),
      styleInformation: BigTextStyleInformation(body),
    );

    const iosDetails = DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );

    final details = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );

    await _localNotifications.show(id, title, body, details, payload: payload);
  }

  /// Show transfer notification
  Future<void> showTransferNotification({
    required String senderName,
    required double amount,
    required String currency,
    String? transferId,
  }) async {
    await showNotification(
      id: DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title: 'Transfert reçu 💸',
      body: 'Vous avez reçu ${amount.toStringAsFixed(2)} $currency de $senderName',
      payload: jsonEncode({
        'type': 'transfer',
        'transfer_id': transferId,
      }),
      channelId: 'transfers',
      channelName: 'Transferts',
    );
  }

  /// Show card payment notification
  Future<void> showCardPaymentNotification({
    required String merchantName,
    required double amount,
    required String currency,
  }) async {
    await showNotification(
      id: DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title: 'Paiement carte 💳',
      body: 'Paiement de ${amount.toStringAsFixed(2)} $currency chez $merchantName',
      payload: jsonEncode({'type': 'card_payment'}),
      channelId: 'cards',
      channelName: 'Carte bancaire',
    );
  }

  /// Show security alert notification
  Future<void> showSecurityNotification({
    required String message,
    String? eventType,
  }) async {
    await showNotification(
      id: DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title: 'Alerte sécurité 🔒',
      body: message,
      payload: jsonEncode({'type': 'security', 'event': eventType}),
      channelId: 'security',
      channelName: 'Sécurité',
    );
  }

  /// Show merchant payment received notification
  Future<void> showMerchantPaymentNotification({
    required double amount,
    required String currency,
    String? paymentId,
  }) async {
    await showNotification(
      id: DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title: 'Paiement reçu 🏪',
      body: 'Vous avez reçu ${amount.toStringAsFixed(2)} $currency',
      payload: jsonEncode({
        'type': 'merchant_payment',
        'payment_id': paymentId,
      }),
      channelId: 'merchant',
      channelName: 'Paiements marchands',
    );
  }

  /// Show exchange notification
  Future<void> showExchangeNotification({
    required double fromAmount,
    required String fromCurrency,
    required double toAmount,
    required String toCurrency,
    String? exchangeId,
  }) async {
    await showNotification(
      id: DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title: 'Échange réussi 💱',
      body: 'Vous avez échangé ${fromAmount.toStringAsFixed(2)} $fromCurrency contre ${toAmount.toStringAsFixed(2)} $toCurrency',
      payload: jsonEncode({
        'type': 'exchange',
        'exchange_id': exchangeId,
      }),
      channelId: 'exchange',
      channelName: 'Échanges',
    );
  }

  /// Register device token with backend
  Future<void> registerDeviceToken(String token) async {
    try {
      await _apiClient.post(
        '/notification-service/api/v1/devices/register',
        data: {
          'token': token,
          'platform': Platform.isIOS ? 'ios' : 'android',
          'device_id': await _getDeviceId(),
        },
      );
      await _storage.write(key: 'push_token', value: token);
      debugPrint('Device token registered successfully');
    } catch (e) {
      debugPrint('Failed to register device token: $e');
    }
  }

  /// Unregister device token from backend
  Future<void> unregisterDeviceToken() async {
    try {
      final token = await _storage.read(key: 'push_token');
      if (token != null) {
        await _apiClient.delete(
          '/notification-service/api/v1/devices/unregister',
          data: {'token': token},
        );
        await _storage.delete(key: 'push_token');
      }
    } catch (e) {
      debugPrint('Failed to unregister device token: $e');
    }
  }

  /// Get stored device ID or generate new one
  Future<String> _getDeviceId() async {
    String? deviceId = await _storage.read(key: 'device_id');
    if (deviceId == null) {
      deviceId = 'device_${DateTime.now().millisecondsSinceEpoch}';
      await _storage.write(key: 'device_id', value: deviceId);
    }
    return deviceId;
  }

  /// Cancel all notifications
  Future<void> cancelAllNotifications() async {
    await _localNotifications.cancelAll();
  }

  /// Cancel specific notification
  Future<void> cancelNotification(int id) async {
    await _localNotifications.cancel(id);
  }

  /// Get pending notifications
  Future<List<PendingNotificationRequest>> getPendingNotifications() async {
    return _localNotifications.pendingNotificationRequests();
  }

  /// Schedule a notification for later
  Future<void> scheduleNotification({
    required int id,
    required String title,
    required String body,
    required DateTime scheduledDate,
    String? payload,
  }) async {
    final androidDetails = AndroidNotificationDetails(
      'scheduled_notifications',
      'Scheduled Notifications',
      channelDescription: 'Scheduled notifications',
      importance: Importance.high,
      priority: Priority.high,
    );

    const iosDetails = DarwinNotificationDetails();

    final details = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );

    // Note: For proper timezone scheduling, use flutter_timezone package
    await _localNotifications.show(
      id,
      title,
      body,
      details,
      payload: payload,
    );
  }
}
