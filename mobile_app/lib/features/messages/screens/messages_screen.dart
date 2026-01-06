import 'dart:io';
import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:image_picker/image_picker.dart';
import 'package:file_picker/file_picker.dart';
import 'package:dio/dio.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../core/services/association_api_service.dart';
import '../../core/services/api_client.dart';
import '../../../core/services/contacts_sync_service.dart';
import '../../../core/services/websocket_service.dart';
import 'widgets/pay_contribution_sheet.dart';

class MessagesScreen extends StatefulWidget {
  const MessagesScreen({super.key});

  @override
  State<MessagesScreen> createState() => _MessagesScreenState();
}

class _MessagesScreenState extends State<MessagesScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final AssociationApiService _api = AssociationApiService(ApiClient().dio);
  final ContactsSyncService _contactsService = ContactsSyncService();
  final WebSocketService _wsService = WebSocketService();
  
  List<Map<String, dynamic>> _userConversations = [];
  List<Map<String, dynamic>> _associationChats = [];
  bool _loading = true;
  Timer? _chatActivityTimer; // Added

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _initContacts();
    _initWebSocket();
    _loadData();
    _setChatActivity(true); // Added
    // Heartbeat every 30 seconds
    _chatActivityTimer = Timer.periodic(const Duration(seconds: 30), (_) { // Added
      _setChatActivity(true); // Added
    }); // Added
  }

  Future<void> _setChatActivity(bool active) async { // Added
    try { // Added
      await ApiClient().dio.post('/auth-service/api/v1/users/chat-activity', data: {'active': active}); // Added
    } catch (e) { // Added
      // Ignore errors // Added
    } // Added
  } // Added

  Future<void> _initWebSocket() async {
    // Get current user ID from secure storage
    final prefs = await SharedPreferences.getInstance();
    if (userJson != null) {
      try {
        final user = jsonDecode(userJson);
        final userId = user['id']?.toString();
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
        // Reload conversations to get the new message
        _loadData();
        break;
      case WsMessageType.presence:
        // Handle presence updates
        break;
    }
  }

  Future<void> _initContacts() async {
    // Load from cache first for instant display
    await _contactsService.loadFromCache();
    // Then sync in background
    await _contactsService.syncContacts();
    // Refresh UI if mounted
    if (mounted) setState(() {});
  }

  /// Get display name for a conversation
  /// Uses contact name if phone/email is in device contacts, otherwise fallback
  String _getConversationDisplayName(Map<String, dynamic> conv) {
    // Try to get participant info from conversation
    final participants = conv['participants'] as List?;
    if (participants != null && participants.isNotEmpty) {
      for (final p in participants) {
        if (p is Map) {
          final phone = p['phone']?.toString();
          final email = p['email']?.toString();
          // Use contacts service to get name if available
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
    // Fallback to conversation name or default
    return conv['name']?.toString() ?? 'Conversation';
  }

  @override
  void dispose() {
    _tabController.dispose();
    _wsService.disconnect();
    _chatActivityTimer?.cancel(); // Added
    _setChatActivity(false); // Added
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() => _loading = true);
    try {
      // Load associations for chat tab
      final response = await _api.getMyAssociations();
      setState(() {
        _associationChats = (response.data is List ? response.data : [])
            .map((a) => Map<String, dynamic>.from(a))
            .toList();
      });
    } catch (e) {
      debugPrint('Failed to load data: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  void _openChat(Map<String, dynamic> item, bool isAssociation) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => ChatScreen(
          chatId: item['id'],
          chatName: item['name'] ?? 'Chat',
          isAssociation: isAssociation,
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
              _buildTabBar(),
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

  Widget _buildTabBar() {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.1),
        border: Border(bottom: BorderSide(color: Colors.white.withOpacity(0.2))),
      ),
      child: TabBar(
        controller: _tabController,
        indicatorColor: const Color(0xFF25D366),
        indicatorWeight: 3,
        labelColor: Colors.white,
        unselectedLabelColor: Colors.white.withOpacity(0.6),
        tabs: const [
          Tab(text: 'Utilisateurs'),
          Tab(text: 'Associations'),
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
      child: TabBarView(
        controller: _tabController,
        children: [
          _buildUserConversations(),
          _buildAssociationChats(),
        ],
      ),
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
          onTap: () => _openChat(conv, false),
          leading: CircleAvatar(
            backgroundColor: const Color(0xFF25D366),
            child: Text(
              displayName[0].toUpperCase(),
              style: const TextStyle(color: Colors.white),
            ),
          ),
          title: Text(
            displayName,
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
          subtitle: Text(
            conv['lastMessage'] ?? '',
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          trailing: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                conv['time'] ?? '',
                style: TextStyle(color: Colors.grey[600], fontSize: 12),
              ),
              if ((conv['unread'] ?? 0) > 0)
                Container(
                  margin: const EdgeInsets.only(top: 4),
                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                  decoration: const BoxDecoration(
                    color: Color(0xFF25D366),
                    shape: BoxShape.circle,
                  ),
                  child: Text(
                    '${conv['unread']}',
                    style: const TextStyle(color: Colors.white, fontSize: 10),
                  ),
                ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildAssociationChats() {
    if (_associationChats.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.group_outlined, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              'Aucune association',
              style: TextStyle(color: Colors.grey[600], fontSize: 16),
            ),
          ],
        ),
      );
    }

    return ListView.separated(
      itemCount: _associationChats.length,
      separatorBuilder: (context, index) => const Divider(height: 1),
      itemBuilder: (context, index) {
        final assoc = _associationChats[index];
        return ListTile(
          onTap: () => _openChat(assoc, true),
          leading: const CircleAvatar(
            backgroundColor: Color(0xFF5b6ecd),
            child: Icon(Icons.people, color: Colors.white),
          ),
          title: Text(
            assoc['name'] ?? 'Association',
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
          subtitle: Text(
            '${assoc['total_members'] ?? 0} membres',
            style: TextStyle(color: Colors.grey[600], fontSize: 13),
          ),
        );
      },
    );
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
                    setState(() {
                      searchResults = [];
                    });
                    return;
                  }

                  setState(() => isSearching = true);

                  try {
                    final response = await _api.dio.get(
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
      // TODO: Get current user name from secure storage or state
      final response = await _api.dio.post(
        '/messaging-service/api/v1/conversations',
        data: {
          'participant_id': user['id'],
          'participant_name': user['name'] ?? 'Utilisateur',
          'participant_email': user['email'] ?? '',
          'participant_phone': user['phone'] ?? '',
          'my_name': 'Moi', // TODO: Get from user profile
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

      _openChat(conversation, false);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e')),
      );
    }
  }
}

// Chat Screen
class ChatScreen extends StatefulWidget {
  final String chatId;
  final String chatName;
  final bool isAssociation;

  const ChatScreen({
    super.key,
    required this.chatId,
    required this.chatName,
    required this.isAssociation,
  });

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  final AssociationApiService _api = AssociationApiService(ApiClient().dio);
  
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
    if (!widget.isAssociation) {
      setState(() => _loading = false);
      return;
    }

    try {
      final response = await _api.getAssociationMessages(widget.chatId);
      setState(() {
        _messages = (response.data is List ? response.data : [])
            .map((m) => Map<String, dynamic>.from(m))
            .toList()
            .reversed
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

    if (!widget.isAssociation) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Messagerie utilisateur bientôt disponible')),
      );
      return;
    }

    try {
      await _api.sendAssociationMessage(widget.chatId, content);
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
          CircleAvatar(
            backgroundColor: widget.isAssociation ? const Color(0xFF5b6ecd) : const Color(0xFF25D366),
            child: Icon(
              widget.isAssociation ? Icons.people : Icons.person,
              color: Colors.white,
            ),
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
      decoration: BoxDecoration(
        color: const Color(0xFFECE5DD),
        image: DecorationImage(
          image: const AssetImage('assets/images/chat_bg.png'),
          fit: BoxFit.cover,
          colorFilter: ColorFilter.mode(
            Colors.white.withOpacity(0.9),
            BlendMode.lighten,
          ),
        ),
      ),
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
          color: isMine ? const Color(0xFF25D366) : Colors.white,
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
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              content,
              style: TextStyle(
                color: isMine ? Colors.white : Colors.black87,
                fontSize: 15,
              ),
              softWrap: true,
            ),
            const SizedBox(height: 4),
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  timestamp,
                  style: TextStyle(
                    color: isMine ? Colors.white70 : Colors.grey[600],
                    fontSize: 11,
                  ),
                ),
                if (isMine) ...[
                  const SizedBox(width: 4),
                  Icon(
                    Icons.done_all,
                    size: 14,
                    color: Colors.white70,
                  ),
                ],
              ],
            ),
          ],
        ),
      ),
    );
  }
}
