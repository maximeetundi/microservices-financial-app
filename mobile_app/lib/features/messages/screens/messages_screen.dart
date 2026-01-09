import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../core/services/api_client.dart';
import '../../../core/services/contacts_sync_service.dart';
import '../../../core/services/websocket_service.dart';

class MessagesScreen extends StatefulWidget {
  const MessagesScreen({super.key});

  @override
  State<MessagesScreen> createState() => _MessagesScreenState();
}

class _MessagesScreenState extends State<MessagesScreen> {
  final ContactsSyncService _contactsService = ContactsSyncService();
  final WebSocketService _wsService = WebSocketService();
  
  List<Map<String, dynamic>> _userConversations = [];
  bool _loading = true;
  Timer? _chatActivityTimer;
  String _currentUserId = '';

  @override
  void initState() {
    super.initState();
    _initContacts();
    _initWebSocket();
    _loadData();
    _setChatActivity(true);
    _chatActivityTimer = Timer.periodic(const Duration(seconds: 30), (_) {
      _setChatActivity(true);
    });
  }

  Future<void> _setChatActivity(bool active) async {
    try {
      await ApiClient().dio.post('/auth-service/api/v1/users/chat-activity', data: {'active': active});
    } catch (e) {
      // Ignore errors
    }
  }

  Future<void> _initWebSocket() async {
    final prefs = await SharedPreferences.getInstance();
    final userJson = prefs.getString('user');
    if (userJson != null) {
      try {
        final user = jsonDecode(userJson);
        final userId = user['id']?.toString();
        _currentUserId = userId ?? '';
        if (userId != null && userId.isNotEmpty) {
          _wsService.connect(userId);
          _wsService.onMessage(_handleWebSocketMessage);
        }
      } catch (e) {
        debugPrint('Failed to parse user for WebSocket: $e');
      }
    }
  }

  void _handleWebSocketMessage(WsMessage msg) {
    switch (msg.type) {
      case WsMessageType.newMessage:
        _loadData();
        break;
      case WsMessageType.presence:
        break;
    }
  }

  Future<void> _initContacts() async {
    await _contactsService.loadFromCache();
    await _contactsService.syncContacts();
    if (mounted) setState(() {});
  }

  String _getConversationDisplayName(Map<String, dynamic> conv) {
    final participants = conv['participants'] as List?;
    if (participants != null && participants.isNotEmpty) {
      for (final p in participants) {
        if (p is Map) {
          final participantId = p['user_id']?.toString();
          if (participantId != _currentUserId) {
            final phone = p['phone']?.toString();
            final email = p['email']?.toString();
            final contactName = _contactsService.getDisplayName(
              phone: phone,
              email: email,
              fallbackName: p['name']?.toString(),
            );
            if (contactName != 'Utilisateur') {
              return contactName;
            }
          }
        }
      }
    }
    return conv['name']?.toString() ?? 'Conversation';
  }

  @override
  void dispose() {
    _wsService.disconnect();
    _chatActivityTimer?.cancel();
    _setChatActivity(false);
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() => _loading = true);
    try {
      final response = await ApiClient().dio.get('/messaging-service/api/v1/conversations');
      setState(() {
        _userConversations = (response.data?['conversations'] is List ? response.data['conversations'] : [])
            .map((c) => Map<String, dynamic>.from(c))
            .toList();
      });
    } catch (e) {
      debugPrint('Failed to load data: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  void _openChat(Map<String, dynamic> conversation) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => _ChatScreen(
          conversationId: conversation['id'],
          chatName: _getConversationDisplayName(conversation),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF075E54), Color(0xFF128C7E)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildHeader(),
              Expanded(child: _buildContent()),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showNewConversationDialog,
        backgroundColor: const Color(0xFF25D366),
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.all(20),
      child: const Row(
        children: [
          Text(
            'Messages',
            style: TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildContent() {
    if (_loading) {
      return const Center(
        child: CircularProgressIndicator(color: Colors.white),
      );
    }

    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(24),
          topRight: Radius.circular(24),
        ),
      ),
      child: _buildUserConversations(),
    );
  }

  Widget _buildUserConversations() {
    if (_userConversations.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.chat_bubble_outline, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              'Aucune conversation',
              style: TextStyle(color: Colors.grey[600], fontSize: 16),
            ),
          ],
        ),
      );
    }

    return ListView.separated(
      itemCount: _userConversations.length,
      separatorBuilder: (context, index) => const Divider(height: 1),
      itemBuilder: (context, index) {
        final conv = _userConversations[index];
        final displayName = _getConversationDisplayName(conv);
        return ListTile(
          onTap: () => _openChat(conv),
          leading: CircleAvatar(
            backgroundColor: const Color(0xFF25D366),
            child: Text(
              displayName.isNotEmpty ? displayName[0].toUpperCase() : 'U',
              style: const TextStyle(color: Colors.white),
            ),
          ),
          title: Text(
            displayName,
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
          subtitle: Text(
            conv['last_message']?.toString() ?? '',
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          trailing: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                _formatTime(conv['updated_at']),
                style: TextStyle(color: Colors.grey[600], fontSize: 12),
              ),
              if ((conv['unread_count'] ?? 0) > 0)
                Container(
                  margin: const EdgeInsets.only(top: 4),
                  padding: const EdgeInsets.all(6),
                  decoration: const BoxDecoration(
                    color: Color(0xFF25D366),
                    shape: BoxShape.circle,
                  ),
                  child: Text(
                    '${conv['unread_count']}',
                    style: const TextStyle(color: Colors.white, fontSize: 10),
                  ),
                ),
            ],
          ),
        );
      },
    );
  }

  String _formatTime(dynamic timestamp) {
    if (timestamp == null) return '';
    try {
      final date = DateTime.parse(timestamp.toString());
      final now = DateTime.now();
      final diff = now.difference(date);
      
      if (diff.inHours < 24) {
        return DateFormat('HH:mm').format(date);
      } else if (diff.inDays < 2) {
        return 'Hier';
      } else {
        return DateFormat('dd/MM').format(date);
      }
    } catch (e) {
      return '';
    }
  }

  void _showNewConversationDialog() {
    final searchController = TextEditingController();
    List<Map<String, dynamic>> searchResults = [];
    bool isSearching = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom,
            left: 20,
            right: 20,
            top: 20,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Nouvelle conversation',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  IconButton(
                    icon: const Icon(Icons.close),
                    onPressed: () => Navigator.pop(context),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              TextField(
                controller: searchController,
                decoration: const InputDecoration(
                  hintText: 'Email ou numéro...',
                  prefixIcon: Icon(Icons.search),
                  border: OutlineInputBorder(),
                ),
                onChanged: (query) async {
                  if (query.length < 3) {
                    setState(() => searchResults = []);
                    return;
                  }

                  setState(() => isSearching = true);

                  try {
                    final response = await ApiClient().dio.get(
                      '/auth-service/api/v1/users/search',
                      queryParameters: {'q': query},
                    );
                    setState(() {
                      searchResults = (response.data['users'] as List)
                          .map((u) => Map<String, dynamic>.from(u))
                          .toList();
                      isSearching = false;
                    });
                  } catch (e) {
                    setState(() => isSearching = false);
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('Erreur: $e')),
                    );
                  }
                },
              ),
              const SizedBox(height: 16),
              if (isSearching)
                const Center(child: CircularProgressIndicator())
              else if (searchResults.isNotEmpty)
                SizedBox(
                  height: 300,
                  child: ListView.builder(
                    itemCount: searchResults.length,
                    itemBuilder: (context, index) {
                      final user = searchResults[index];
                      return ListTile(
                        leading: CircleAvatar(
                          backgroundColor: const Color(0xFF25D366),
                          child: Text(
                            (user['name'] ?? 'U')[0].toUpperCase(),
                            style: const TextStyle(color: Colors.white),
                          ),
                        ),
                        title: Text(user['name'] ?? 'Utilisateur'),
                        subtitle: Text(user['email'] ?? user['phone'] ?? ''),
                        onTap: () async {
                          Navigator.pop(context);
                          await _createConversation(user);
                        },
                      );
                    },
                  ),
                )
              else if (searchController.text.length >= 3)
                const Center(
                  child: Padding(
                    padding: EdgeInsets.all(20),
                    child: Text('Aucun utilisateur trouvé'),
                  ),
                )
              else
                const Center(
                  child: Padding(
                    padding: EdgeInsets.all(20),
                    child: Text(
                      'Tapez au moins 3 caractères pour rechercher',
                      style: TextStyle(color: Colors.grey),
                    ),
                  ),
                ),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _createConversation(Map<String, dynamic> user) async {
    try {
      final response = await ApiClient().dio.post(
        '/messaging-service/api/v1/conversations',
        data: {
          'participant_id': user['id'],
          'participant_name': user['name'] ?? 'Utilisateur',
          'participant_email': user['email'] ?? '',
          'participant_phone': user['phone'] ?? '',
          'my_name': 'Moi',
        },
      );

      final conversation = Map<String, dynamic>.from(response.data);
      conversation['name'] = user['name'] ?? user['email'] ?? user['phone'];

      setState(() {
        _userConversations.insert(0, conversation);
      });

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Conversation créée !')),
      );

      _openChat(conversation);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e')),
      );
    }
  }
}

// Simple Chat Screen for user-to-user messaging
class _ChatScreen extends StatefulWidget {
  final String conversationId;
  final String chatName;

  const _ChatScreen({
    required this.conversationId,
    required this.chatName,
  });

  @override
  State<_ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<_ChatScreen> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  
  List<Map<String, dynamic>> _messages = [];
  bool _loading = true;
  String _currentUserId = '';

  @override
  void initState() {
    super.initState();
    _loadCurrentUser();
    _loadMessages();
  }

  Future<void> _loadCurrentUser() async {
    final prefs = await SharedPreferences.getInstance();
    final userJson = prefs.getString('user');
    if (userJson != null) {
      try {
        final user = jsonDecode(userJson);
        setState(() {
          _currentUserId = user['id']?.toString() ?? '';
        });
      } catch (e) {
        debugPrint('Failed to parse user: $e');
      }
    }
  }

  @override
  void dispose() {
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _loadMessages() async {
    try {
      final response = await ApiClient().dio.get(
        '/messaging-service/api/v1/conversations/${widget.conversationId}/messages'
      );
      setState(() {
        _messages = (response.data?['messages'] is List ? response.data['messages'] : [])
            .map((m) => Map<String, dynamic>.from(m))
            .toList();
        _loading = false;
      });
      _scrollToBottom();
    } catch (e) {
      debugPrint('Failed to load messages: $e');
      setState(() => _loading = false);
    }
  }

  Future<void> _sendMessage() async {
    if (_messageController.text.trim().isEmpty) return;

    final content = _messageController.text;
    _messageController.clear();

    try {
      await ApiClient().dio.post(
        '/messaging-service/api/v1/conversations/${widget.conversationId}/messages',
        data: {
          'content': content,
          'message_type': 'text',
        },
      );
      _loadMessages();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e')),
      );
    }
  }

  void _scrollToBottom() {
    if (_scrollController.hasClients) {
      Future.delayed(const Duration(milliseconds: 100), () {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            colors: [Color(0xFF075E54), Color(0xFF128C7E)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildChatHeader(),
              Expanded(child: _buildMessagesArea()),
              _buildInputArea(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildChatHeader() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      color: const Color(0xFF075E54),
      child: Row(
        children: [
          IconButton(
            icon: const Icon(Icons.arrow_back, color: Colors.white),
            onPressed: () => Navigator.pop(context),
          ),
          const CircleAvatar(
            backgroundColor: Color(0xFF25D366),
            child: Icon(Icons.person, color: Colors.white),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.chatName,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                  overflow: TextOverflow.ellipsis,
                  maxLines: 1,
                ),
                Text(
                  'En ligne',
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.8),
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMessagesArea() {
    return Container(
      color: const Color(0xFFECE5DD),
      child: _loading
          ? const Center(child: CircularProgressIndicator())
          : ListView.builder(
              controller: _scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: _messages.length,
              itemBuilder: (context, index) {
                final message = _messages[index];
                final senderId = message['sender_id']?.toString() ?? '';
                final isMine = senderId == _currentUserId;
                return _buildMessageBubble(message, isMine);
              },
            ),
    );
  }

  Widget _buildInputArea() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 4,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _messageController,
              decoration: InputDecoration(
                hintText: 'Message...',
                filled: true,
                fillColor: Colors.grey[100],
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide.none,
                ),
                contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
              ),
              maxLines: null,
              textInputAction: TextInputAction.send,
              onSubmitted: (_) => _sendMessage(),
            ),
          ),
          const SizedBox(width: 8),
          CircleAvatar(
            backgroundColor: const Color(0xFF25D366),
            child: IconButton(
              icon: const Icon(Icons.send, color: Colors.white, size: 20),
              onPressed: _sendMessage,
            ),
          ),
        ],
      ),
    );
  }

  String _formatTime(dynamic timestamp) {
    if (timestamp == null) return '';
    try {
      final date = DateTime.parse(timestamp.toString());
      return DateFormat('HH:mm').format(date);
    } catch (e) {
      return '';
    }
  }

  Widget _buildMessageBubble(Map<String, dynamic> message, bool isMine) {
    final content = message['content']?.toString() ?? '';
    final timestamp = _formatTime(message['created_at']);
    
    return Align(
      alignment: isMine ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        constraints: BoxConstraints(
          maxWidth: MediaQuery.of(context).size.width * 0.75,
        ),
        decoration: BoxDecoration(
          color: isMine ? const Color(0xFFDCF8C6) : Colors.white,
          borderRadius: BorderRadius.only(
            topLeft: const Radius.circular(12),
            topRight: const Radius.circular(12),
            bottomLeft: Radius.circular(isMine ? 12 : 0),
            bottomRight: Radius.circular(isMine ? 0 : 12),
          ),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 4,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Text(
              content,
              style: const TextStyle(fontSize: 15),
            ),
            const SizedBox(height: 4),
            Text(
              timestamp,
              style: TextStyle(
                fontSize: 11,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
