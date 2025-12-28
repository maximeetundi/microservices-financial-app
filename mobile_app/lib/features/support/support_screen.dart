import 'package:flutter/material.dart';
import '../../core/services/support_api_service.dart';

class SupportScreen extends StatefulWidget {
  const SupportScreen({super.key});

  @override
  State<SupportScreen> createState() => _SupportScreenState();
}

class _SupportScreenState extends State<SupportScreen> {
  String? selectedAgentType;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A),
      body: SafeArea(
        child: Column(
          children: [
            // Header
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    Colors.purple.withOpacity(0.2),
                    Colors.blue.withOpacity(0.1),
                  ],
                ),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      IconButton(
                        onPressed: () => Navigator.pop(context),
                        icon: const Icon(Icons.arrow_back, color: Colors.white),
                      ),
                      const SizedBox(width: 8),
                      const Text(
                        'Centre d\'Assistance',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'Nous sommes là pour vous aider 24/7',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.7),
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ),
            
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Stats Row
                    Row(
                      children: [
                        _buildStatCard('⚡', 'Temps de réponse', '~2 min', Colors.green),
                        const SizedBox(width: 12),
                        _buildStatCard('🕐', 'Disponibilité', '24h/24', Colors.blue),
                      ],
                    ),
                    const SizedBox(height: 32),

                    // Agent Selection
                    const Text(
                      'Comment souhaitez-vous être aidé ?',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 16),

                    // AI Agent Card
                    _buildAgentCard(
                      emoji: '🤖',
                      title: 'Assistant IA',
                      description: 'Réponses instantanées 24/7',
                      features: ['Réponse instantanée', 'Disponible 24h/24', 'Escalade vers humain possible'],
                      gradient: [Colors.blue.withOpacity(0.3), Colors.purple.withOpacity(0.2)],
                      isSelected: selectedAgentType == 'ai',
                      onTap: () => setState(() => selectedAgentType = 'ai'),
                    ),
                    const SizedBox(height: 16),

                    // Human Agent Card
                    _buildAgentCard(
                      emoji: '👤',
                      title: 'Conseiller Humain',
                      description: 'Pour les demandes complexes',
                      features: ['Analyse personnalisée', 'Peut prendre des actions', 'Attente: 2-5 min'],
                      gradient: [Colors.teal.withOpacity(0.3), Colors.green.withOpacity(0.2)],
                      isSelected: selectedAgentType == 'human',
                      onTap: () => setState(() => selectedAgentType = 'human'),
                    ),
                    const SizedBox(height: 32),

                    // Start Conversation Button
                    if (selectedAgentType != null) ...[
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: () {
                            Navigator.push(
                              context,
                              MaterialPageRoute(
                                builder: (context) => ChatScreen(agentType: selectedAgentType!),
                              ),
                            );
                          },
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.blue,
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                          child: const Text(
                            'Démarrer la conversation',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatCard(String emoji, String label, String value, Color color) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: color.withOpacity(0.2),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Center(child: Text(emoji, style: const TextStyle(fontSize: 20))),
            ),
            const SizedBox(width: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.6),
                    fontSize: 12,
                  ),
                ),
                Text(
                  value,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAgentCard({
    required String emoji,
    required String title,
    required String description,
    required List<String> features,
    required List<Color> gradient,
    required bool isSelected,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: gradient,
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
          borderRadius: BorderRadius.circular(20),
          border: Border.all(
            color: isSelected ? Colors.blue : Colors.white.withOpacity(0.1),
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  width: 60,
                  height: 60,
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Center(child: Text(emoji, style: const TextStyle(fontSize: 32))),
                ),
                const SizedBox(width: 16),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Text(
                      description,
                      style: TextStyle(
                        color: Colors.white.withOpacity(0.7),
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
                const Spacer(),
                if (isSelected)
                  Container(
                    width: 24,
                    height: 24,
                    decoration: const BoxDecoration(
                      color: Colors.blue,
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(Icons.check, color: Colors.white, size: 16),
                  ),
              ],
            ),
            const SizedBox(height: 16),
            ...features.map((feature) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: Row(
                children: [
                  Icon(
                    Icons.check_circle,
                    color: Colors.green.withOpacity(0.8),
                    size: 20,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    feature,
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.8),
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            )),
          ],
        ),
      ),
    );
  }
}

// Chat Screen
class ChatScreen extends StatefulWidget {
  final String agentType;
  final String? ticketId; // Optional: resume existing ticket

  const ChatScreen({super.key, required this.agentType, this.ticketId});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  final List<ChatMessage> _messages = [];
  bool _isTyping = false;
  String? _ticketId;
  final SupportApiService _supportApi = SupportApiService();

  @override
  void initState() {
    super.initState();
    _ticketId = widget.ticketId;
    _initializeConversation();
  }
  
  Future<void> _initializeConversation() async {
    if (_ticketId != null) {
      // Load existing ticket messages
      await _loadMessages();
    } else {
      // Create new ticket
      await _createTicket();
    }
  }
  
  Future<void> _createTicket() async {
    try {
      final ticket = await _supportApi.createTicket(
        subject: 'Conversation ${widget.agentType == 'ai' ? 'IA' : 'Conseiller'}',
        category: 'general',
        description: 'Nouvelle conversation de support',
        priority: 'normal',
        agentType: widget.agentType,
      );
      setState(() {
        _ticketId = ticket['id']?.toString();
      });
    } catch (e) {
      // Continue with local simulation if API fails
      debugPrint('Failed to create ticket: $e');
    }
    _addWelcomeMessage();
  }
  
  Future<void> _loadMessages() async {
    if (_ticketId == null) return;
    try {
      final messages = await _supportApi.getMessages(_ticketId!);
      setState(() {
        _messages.addAll(messages.map((msg) => ChatMessage(
          id: msg['id']?.toString() ?? '',
          content: msg['content'] ?? '',
          isUser: msg['sender_type'] == 'user',
          timestamp: DateTime.tryParse(msg['created_at'] ?? '') ?? DateTime.now(),
        )));
      });
    } catch (e) {
      debugPrint('Failed to load messages: $e');
    }
    if (_messages.isEmpty) {
      _addWelcomeMessage();
    }
  }

  void _addWelcomeMessage() {
    _messages.add(ChatMessage(
      id: '1',
      content: widget.agentType == 'ai'
          ? 'Bonjour ! 👋 Je suis l\'assistant virtuel Zekora.\n\nComment puis-je vous aider ?\n\n• 💳 Cartes bancaires\n• 💸 Transferts\n• ₿ Cryptomonnaies\n• 📊 Frais\n• 🔐 Sécurité'
          : 'Bonjour ! Un conseiller va prendre en charge votre demande sous peu.\n\n⏱️ Temps d\'attente estimé : 2-5 minutes.',
      isUser: false,
      timestamp: DateTime.now(),
    ));
  }

  void _sendMessage() async {
    if (_messageController.text.trim().isEmpty) return;

    final content = _messageController.text.trim();
    _messageController.clear();

    setState(() {
      _messages.add(ChatMessage(
        id: DateTime.now().millisecondsSinceEpoch.toString(),
        content: content,
        isUser: true,
        timestamp: DateTime.now(),
      ));
      _isTyping = true;
    });

    _scrollToBottom();

    // Simulate AI response
    await Future.delayed(const Duration(milliseconds: 1500));

    if (widget.agentType == 'ai') {
      setState(() {
        _isTyping = false;
        _messages.add(ChatMessage(
          id: DateTime.now().millisecondsSinceEpoch.toString(),
          content: _generateAIResponse(content),
          isUser: false,
          timestamp: DateTime.now(),
        ));
      });
    }

    _scrollToBottom();
  }

  String _generateAIResponse(String message) {
    final lower = message.toLowerCase();
    
    if (lower.contains('solde') || lower.contains('balance')) {
      return 'Pour consulter votre solde, rendez-vous sur l\'écran d\'accueil. Votre solde total s\'affiche en haut de l\'écran.';
    }
    if (lower.contains('frais') || lower.contains('commission')) {
      return '📊 Nos frais :\n• Transferts SEPA : Gratuit\n• Crypto-Crypto : 0.5%\n• Fiat-Crypto : 0.75%\n\nNous sommes jusqu\'à 8x moins chers que les banques !';
    }
    if (lower.contains('carte')) {
      return 'Pour commander une carte :\n1. Allez dans le menu "Cartes"\n2. Cliquez sur "Commander"\n3. Choisissez virtuelle ou physique\n\nVotre carte virtuelle est instantanée !';
    }
    if (lower.contains('humain') || lower.contains('agent')) {
      return 'Je comprends que vous souhaitez parler à un conseiller. Appuyez sur le bouton "Parler à un humain" en haut de l\'écran.';
    }
    
    return 'Je comprends votre demande. Pourriez-vous me donner plus de détails ?\n\nOu utilisez le bouton "Parler à un humain" pour une assistance personnalisée.';
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
        title: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: widget.agentType == 'ai' 
                    ? Colors.blue.withOpacity(0.2) 
                    : Colors.green.withOpacity(0.2),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Center(
                child: Text(
                  widget.agentType == 'ai' ? '🤖' : '👤',
                  style: const TextStyle(fontSize: 20),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.agentType == 'ai' ? 'Assistant IA' : 'Conseiller',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Row(
                  children: [
                    Container(
                      width: 8,
                      height: 8,
                      decoration: const BoxDecoration(
                        color: Colors.green,
                        shape: BoxShape.circle,
                      ),
                    ),
                    const SizedBox(width: 4),
                    Text(
                      'En ligne',
                      style: TextStyle(
                        color: Colors.green.withOpacity(0.8),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ],
        ),
        actions: [
          if (widget.agentType == 'ai')
            TextButton.icon(
              onPressed: () {
                // Escalate to human
                setState(() {
                  _messages.add(ChatMessage(
                    id: 'system',
                    content: '🔔 Votre demande a été transférée à un conseiller humain.',
                    isUser: false,
                    timestamp: DateTime.now(),
                    isSystem: true,
                  ));
                });
              },
              icon: const Text('👤', style: TextStyle(fontSize: 16)),
              label: const Text('Humain', style: TextStyle(color: Colors.orange)),
            ),
        ],
      ),
      body: Column(
        children: [
          // Messages
          Expanded(
            child: ListView.builder(
              controller: _scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: _messages.length + (_isTyping ? 1 : 0),
              itemBuilder: (context, index) {
                if (index == _messages.length && _isTyping) {
                  return _buildTypingIndicator();
                }
                return _buildMessageBubble(_messages[index]);
              },
            ),
          ),

          // Input
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              border: Border(
                top: BorderSide(color: Colors.white.withOpacity(0.1)),
              ),
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _messageController,
                    style: const TextStyle(color: Colors.white),
                    decoration: InputDecoration(
                      hintText: 'Écrivez votre message...',
                      hintStyle: TextStyle(color: Colors.white.withOpacity(0.4)),
                      filled: true,
                      fillColor: Colors.white.withOpacity(0.1),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                    ),
                    onSubmitted: (_) => _sendMessage(),
                  ),
                ),
                const SizedBox(width: 12),
                Container(
                  decoration: BoxDecoration(
                    gradient: const LinearGradient(
                      colors: [Colors.blue, Colors.purple],
                    ),
                    borderRadius: BorderRadius.circular(24),
                  ),
                  child: IconButton(
                    onPressed: _sendMessage,
                    icon: const Icon(Icons.send, color: Colors.white),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMessageBubble(ChatMessage message) {
    if (message.isSystem) {
      return Center(
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 8),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          decoration: BoxDecoration(
            color: Colors.orange.withOpacity(0.2),
            borderRadius: BorderRadius.circular(20),
          ),
          child: Text(
            message.content,
            style: TextStyle(color: Colors.orange.withOpacity(0.9), fontSize: 13),
          ),
        ),
      );
    }

    return Align(
      alignment: message.isUser ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        constraints: BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.75),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          gradient: message.isUser
              ? const LinearGradient(colors: [Colors.blue, Colors.purple])
              : null,
          color: message.isUser ? null : Colors.white.withOpacity(0.1),
          borderRadius: BorderRadius.only(
            topLeft: const Radius.circular(20),
            topRight: const Radius.circular(20),
            bottomLeft: Radius.circular(message.isUser ? 20 : 4),
            bottomRight: Radius.circular(message.isUser ? 4 : 20),
          ),
        ),
        child: Text(
          message.content,
          style: const TextStyle(color: Colors.white, fontSize: 15),
        ),
      ),
    );
  }

  Widget _buildTypingIndicator() {
    return Align(
      alignment: Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.1),
          borderRadius: const BorderRadius.only(
            topLeft: Radius.circular(20),
            topRight: Radius.circular(20),
            bottomLeft: Radius.circular(4),
            bottomRight: Radius.circular(20),
          ),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: List.generate(3, (index) => Container(
            width: 8,
            height: 8,
            margin: const EdgeInsets.symmetric(horizontal: 2),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.5),
              shape: BoxShape.circle,
            ),
          )),
        ),
      ),
    );
  }
}

class ChatMessage {
  final String id;
  final String content;
  final bool isUser;
  final DateTime timestamp;
  final bool isSystem;

  ChatMessage({
    required this.id,
    required this.content,
    required this.isUser,
    required this.timestamp,
    this.isSystem = false,
  });
}
