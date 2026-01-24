import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:flutter/foundation.dart';

/// WebSocket message types
class WsMessageType {
  static const String newMessage = 'new_message';
  static const String typing = 'typing';
  static const String read = 'read';
  static const String presence = 'presence';
}

/// WebSocket message model
class WsMessage {
  final String type;
  final String? conversationId;
  final String? senderId;
  final dynamic content;

  WsMessage({
    required this.type,
    this.conversationId,
    this.senderId,
    this.content,
  });

  factory WsMessage.fromJson(Map<String, dynamic> json) {
    return WsMessage(
      type: json['type'] ?? '',
      conversationId: json['conversation_id'],
      senderId: json['sender_id'],
      content: json['content'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'type': type,
      if (conversationId != null) 'conversation_id': conversationId,
      if (senderId != null) 'sender_id': senderId,
      if (content != null) 'content': content,
    };
  }
}

/// WebSocket service for real-time messaging
class WebSocketService {
  static final WebSocketService _instance = WebSocketService._internal();
  factory WebSocketService() => _instance;
  WebSocketService._internal();

  WebSocketChannel? _channel;
  StreamSubscription? _subscription;
  String? _currentUserId;
  
  bool _isConnected = false;
  bool _isReconnecting = false;
  Timer? _reconnectTimer;
  Timer? _pingTimer;
  
  // Callbacks
  final List<void Function(WsMessage)> _onMessageCallbacks = [];
  final List<void Function()> _onConnectedCallbacks = [];
  final List<void Function()> _onDisconnectedCallbacks = [];

  static const String _wsUrl = 'wss://api.app.tech-afm.com/messaging-service/ws/chat';

  bool get isConnected => _isConnected;

  /// Connect to WebSocket server
  Future<void> connect(String userId) async {
    if (_isConnected || _isReconnecting) return;
    
    _currentUserId = userId;
    
    try {
      debugPrint('WebSocket: Connecting for user $userId...');
      
      final uri = Uri.parse('$_wsUrl?user_id=$userId');
      _channel = WebSocketChannel.connect(uri);
      
      await _channel!.ready;
      
      _isConnected = true;
      debugPrint('WebSocket: Connected');
      
      // Notify connected callbacks
      for (final cb in _onConnectedCallbacks) {
        cb();
      }
      
      // Start ping timer (keep-alive)
      _startPingTimer();
      
      // Listen for messages
      _subscription = _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnect,
      );
      
    } catch (e) {
      debugPrint('WebSocket: Connection error: $e');
      _scheduleReconnect();
    }
  }

  /// Disconnect from WebSocket server
  void disconnect() {
    debugPrint('WebSocket: Disconnecting...');
    _pingTimer?.cancel();
    _reconnectTimer?.cancel();
    _subscription?.cancel();
    _channel?.sink.close();
    _channel = null;
    _isConnected = false;
    _isReconnecting = false;
    
    // Notify disconnected callbacks
    for (final cb in _onDisconnectedCallbacks) {
      cb();
    }
  }

  /// Send a message through WebSocket
  void send(WsMessage message) {
    if (!_isConnected || _channel == null) {
      debugPrint('WebSocket: Not connected, cannot send');
      return;
    }
    
    _channel!.sink.add(jsonEncode(message.toJson()));
  }

  /// Send typing indicator
  void sendTyping(String conversationId) {
    send(WsMessage(
      type: WsMessageType.typing,
      conversationId: conversationId,
    ));
  }

  /// Send read receipt
  void sendRead(String conversationId, String messageId) {
    send(WsMessage(
      type: WsMessageType.read,
      conversationId: conversationId,
      content: {'message_id': messageId},
    ));
  }

  /// Register message callback
  void onMessage(void Function(WsMessage) callback) {
    _onMessageCallbacks.add(callback);
  }

  /// Unregister message callback
  void offMessage(void Function(WsMessage) callback) {
    _onMessageCallbacks.remove(callback);
  }

  /// Register connected callback
  void onConnected(void Function() callback) {
    _onConnectedCallbacks.add(callback);
  }

  /// Register disconnected callback
  void onDisconnected(void Function() callback) {
    _onDisconnectedCallbacks.add(callback);
  }

  void _handleMessage(dynamic data) {
    try {
      final json = jsonDecode(data as String);
      final message = WsMessage.fromJson(json);
      
      debugPrint('WebSocket: Received ${message.type}');
      
      // Notify all callbacks
      for (final cb in _onMessageCallbacks) {
        cb(message);
      }
    } catch (e) {
      debugPrint('WebSocket: Message parse error: $e');
    }
  }

  void _handleError(dynamic error) {
    debugPrint('WebSocket: Error: $error');
    _scheduleReconnect();
  }

  void _handleDisconnect() {
    debugPrint('WebSocket: Disconnected');
    _isConnected = false;
    _pingTimer?.cancel();
    
    // Notify disconnected callbacks
    for (final cb in _onDisconnectedCallbacks) {
      cb();
    }
    
    _scheduleReconnect();
  }

  void _scheduleReconnect() {
    if (_isReconnecting || _currentUserId == null) return;
    
    _isReconnecting = true;
    debugPrint('WebSocket: Scheduling reconnect in 3 seconds...');
    
    _reconnectTimer = Timer(const Duration(seconds: 3), () {
      _isReconnecting = false;
      if (_currentUserId != null) {
        connect(_currentUserId!);
      }
    });
  }

  void _startPingTimer() {
    _pingTimer?.cancel();
    _pingTimer = Timer.periodic(const Duration(seconds: 30), (_) {
      if (_isConnected) {
        // Send ping (the server will respond with pong)
        debugPrint('WebSocket: Ping');
      }
    });
  }
}
