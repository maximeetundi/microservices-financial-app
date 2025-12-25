import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'dart:io';
import 'package:image_picker/image_picker.dart';

import '../../../../core/widgets/glass_container.dart';
import '../../../../core/services/auth_api_service.dart';

/// KYC Verification Page matching web design
class KycPage extends StatefulWidget {
  const KycPage({super.key});

  @override
  State<KycPage> createState() => _KycPageState();
}

class _KycPageState extends State<KycPage> {
  final AuthApiService _authApi = AuthApiService();
  final ImagePicker _picker = ImagePicker();
  
  String _kycStatus = 'pending'; // pending, submitted, verified, rejected
  int _kycLevel = 0;
  bool _isLoading = true;
  bool _isUploading = false;
  String? _selectedDocType;
  File? _selectedFile;
  
  List<Map<String, dynamic>> _documents = [];
  
  final Map<String, Map<String, dynamic>> _documentStatus = {
    'identity': {'status': 'required', 'label': 'Requis'},
    'selfie': {'status': 'required', 'label': 'Requis'},
    'address': {'status': 'required', 'label': 'Requis'},
  };

  @override
  void initState() {
    super.initState();
    _loadKycStatus();
  }

  Future<void> _loadKycStatus() async {
    try {
      final profile = await _authApi.getProfile();
      final kycDocs = await _authApi.getKYCDocuments();
      
      setState(() {
        _kycStatus = profile['kyc_status'] ?? 'pending';
        _kycLevel = profile['kyc_level'] ?? 0;
        _documents = List<Map<String, dynamic>>.from(kycDocs['documents'] ?? []);
        
        // Update document statuses
        for (var doc in _documents) {
          final type = doc['type'];
          if (_documentStatus.containsKey(type)) {
            _documentStatus[type] = {
              'status': doc['status'] ?? 'pending',
              'label': _getStatusLabel(doc['status']),
            };
          }
        }
        
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      debugPrint('Failed to load KYC status: $e');
    }
  }

  String _getStatusLabel(String? status) {
    switch (status) {
      case 'approved':
        return 'Approuvé';
      case 'rejected':
        return 'Refusé';
      case 'pending':
        return 'En cours';
      default:
        return 'Requis';
    }
  }

  Future<void> _pickDocument() async {
    if (_selectedDocType == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Veuillez sélectionner un type de document')),
      );
      return;
    }

    final XFile? image = await _picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1920,
      maxHeight: 1920,
      imageQuality: 85,
    );

    if (image != null) {
      setState(() {
        _selectedFile = File(image.path);
      });
    }
  }

  Future<void> _uploadDocument() async {
    if (_selectedFile == null || _selectedDocType == null) return;

    setState(() => _isUploading = true);

    try {
      await _authApi.uploadKYCDocument(_selectedDocType!, _selectedFile!);
      
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Document envoyé avec succès!'),
          backgroundColor: Colors.green,
        ),
      );

      setState(() {
        _documentStatus[_selectedDocType!] = {'status': 'pending', 'label': 'En cours'};
        _selectedFile = null;
        _selectedDocType = null;
      });

      await _loadKycStatus();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
      );
    } finally {
      setState(() => _isUploading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      body: SafeArea(
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  _buildHeader(isDark),
                  const SizedBox(height: 24),
                  _buildStatusCard(isDark),
                  const SizedBox(height: 24),
                  _buildDocumentsSection(isDark),
                  if (_kycStatus != 'verified') ...[
                    const SizedBox(height: 24),
                    _buildUploadSection(isDark),
                  ],
                  const SizedBox(height: 24),
                  _buildInfoBox(isDark),
                ],
              ),
      ),
    );
  }

  Widget _buildHeader(bool isDark) {
    return Row(
      children: [
        GlassContainer(
          padding: EdgeInsets.zero,
          width: 40,
          height: 40,
          borderRadius: 12,
          child: IconButton(
            icon: Icon(Icons.arrow_back_ios_new, size: 20,
                color: isDark ? Colors.white : const Color(0xFF1E293B)),
            onPressed: () => context.go('/more/settings'),
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '📋 Vérification KYC',
                style: GoogleFonts.inter(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                'Validez votre identité pour débloquer toutes les fonctionnalités',
                style: GoogleFonts.inter(
                  fontSize: 12,
                  color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildStatusCard(bool isDark) {
    final statusConfig = _getStatusConfig();
    
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: statusConfig['bgColor'],
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: statusConfig['borderColor']),
      ),
      child: Row(
        children: [
          Text(statusConfig['icon'], style: const TextStyle(fontSize: 40)),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  statusConfig['title'],
                  style: GoogleFonts.inter(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  statusConfig['description'],
                  style: GoogleFonts.inter(
                    fontSize: 13,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Map<String, dynamic> _getStatusConfig() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    switch (_kycStatus) {
      case 'verified':
        return {
          'icon': '✅',
          'title': 'Compte vérifié',
          'description': 'Votre identité a été confirmée. Accès complet activé.',
          'bgColor': isDark 
              ? const Color(0xFF22C55E).withOpacity(0.1)
              : const Color(0xFFF0FDF4),
          'borderColor': isDark
              ? const Color(0xFF22C55E).withOpacity(0.2)
              : const Color(0xFFBBF7D0),
        };
      case 'rejected':
        return {
          'icon': '❌',
          'title': 'Vérification refusée',
          'description': 'Veuillez soumettre de nouveaux documents valides.',
          'bgColor': isDark 
              ? const Color(0xFFEF4444).withOpacity(0.1)
              : const Color(0xFFFEF2F2),
          'borderColor': isDark
              ? const Color(0xFFEF4444).withOpacity(0.2)
              : const Color(0xFFFECACA),
        };
      case 'submitted':
        return {
          'icon': '📨',
          'title': 'Documents en cours de vérification',
          'description': 'Nous examinons vos documents. Réponse sous 24-48h.',
          'bgColor': isDark 
              ? const Color(0xFF3B82F6).withOpacity(0.1)
              : const Color(0xFFEFF6FF),
          'borderColor': isDark
              ? const Color(0xFF3B82F6).withOpacity(0.2)
              : const Color(0xFFBFDBFE),
        };
      default:
        return {
          'icon': '⏳',
          'title': 'Vérification en attente',
          'description': 'Soumettez vos documents pour vérifier votre identité.',
          'bgColor': isDark 
              ? const Color(0xFFF97316).withOpacity(0.1)
              : const Color(0xFFFFF7ED),
          'borderColor': isDark
              ? const Color(0xFFF97316).withOpacity(0.2)
              : const Color(0xFFFED7AA),
        };
    }
  }

  Widget _buildDocumentsSection(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'DOCUMENTS REQUIS',
          style: GoogleFonts.inter(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.2,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 12),
        _buildDocumentItem('🪪', 'Pièce d\'identité', 
            'Passeport, carte d\'identité ou permis de conduire', 'identity', isDark),
        _buildDocumentItem('🤳', 'Selfie avec document',
            'Photo de vous tenant votre pièce d\'identité', 'selfie', isDark),
        _buildDocumentItem('🏠', 'Justificatif de domicile',
            'Facture de moins de 3 mois', 'address', isDark),
      ],
    );
  }

  Widget _buildDocumentItem(String emoji, String title, String subtitle, String type, bool isDark) {
    final status = _documentStatus[type]!;
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.03) : Colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0),
        ),
      ),
      child: Row(
        children: [
          Text(emoji, style: const TextStyle(fontSize: 28)),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  subtitle,
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              color: _getStatusColor(status['status'], isDark),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              status['label'],
              style: GoogleFonts.inter(
                fontSize: 10,
                fontWeight: FontWeight.bold,
                color: _getStatusTextColor(status['status']),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Color _getStatusColor(String status, bool isDark) {
    switch (status) {
      case 'approved':
        return isDark 
            ? const Color(0xFF22C55E).withOpacity(0.2)
            : const Color(0xFFDCFCE7);
      case 'rejected':
        return isDark 
            ? const Color(0xFFEF4444).withOpacity(0.2)
            : const Color(0xFFFEE2E2);
      case 'pending':
        return isDark 
            ? const Color(0xFFF97316).withOpacity(0.2)
            : const Color(0xFFFED7AA);
      default:
        return isDark 
            ? Colors.white.withOpacity(0.1)
            : const Color(0xFFF1F5F9);
    }
  }

  Color _getStatusTextColor(String status) {
    switch (status) {
      case 'approved':
        return const Color(0xFF22C55E);
      case 'rejected':
        return const Color(0xFFEF4444);
      case 'pending':
        return const Color(0xFFF97316);
      default:
        return const Color(0xFF94A3B8);
    }
  }

  Widget _buildUploadSection(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'UPLOADER UN DOCUMENT',
          style: GoogleFonts.inter(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.2,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 12),
        
        // Document Type Selector
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          decoration: BoxDecoration(
            color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
            ),
          ),
          child: DropdownButtonHideUnderline(
            child: DropdownButton<String>(
              value: _selectedDocType,
              isExpanded: true,
              hint: Text(
                'Choisir le type de document',
                style: GoogleFonts.inter(
                  color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                ),
              ),
              dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
              items: const [
                DropdownMenuItem(value: 'identity', child: Text('Pièce d\'identité')),
                DropdownMenuItem(value: 'selfie', child: Text('Selfie avec document')),
                DropdownMenuItem(value: 'address', child: Text('Justificatif de domicile')),
              ],
              onChanged: (value) => setState(() => _selectedDocType = value),
              style: GoogleFonts.inter(
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
          ),
        ),
        const SizedBox(height: 16),
        
        // Upload Zone
        GestureDetector(
          onTap: _pickDocument,
          child: Container(
            padding: const EdgeInsets.all(32),
            decoration: BoxDecoration(
              color: isDark ? Colors.white.withOpacity(0.02) : Colors.white,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: isDark 
                    ? const Color(0xFF6366F1).withOpacity(0.3)
                    : const Color(0xFF6366F1).withOpacity(0.3),
                width: 2,
                style: BorderStyle.solid,
              ),
            ),
            child: Column(
              children: [
                const Text('📤', style: TextStyle(fontSize: 36)),
                const SizedBox(height: 12),
                Text(
                  _selectedFile != null 
                      ? _selectedFile!.path.split('/').last
                      : 'Appuyez pour sélectionner un fichier',
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    color: _selectedFile != null 
                        ? const Color(0xFF6366F1)
                        : (isDark ? Colors.white : const Color(0xFF1E293B)),
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'JPG, PNG ou PDF • Max 10MB',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        
        // Upload Button
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: (_selectedFile != null && _selectedDocType != null && !_isUploading)
                ? _uploadDocument
                : null,
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF6366F1),
              disabledBackgroundColor: isDark 
                  ? Colors.white.withOpacity(0.1)
                  : const Color(0xFFE2E8F0),
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
            child: _isUploading
                ? const SizedBox(
                    height: 20,
                    width: 20,
                    child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                  )
                : Text(
                    '📤 Envoyer le document',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: Colors.white,
                    ),
                  ),
          ),
        ),
      ],
    );
  }

  Widget _buildInfoBox(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark 
            ? const Color(0xFF3B82F6).withOpacity(0.1)
            : const Color(0xFFEFF6FF),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          const Text('ℹ️', style: TextStyle(fontSize: 24)),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              'La vérification prend généralement 24 à 48 heures. Vous recevrez une notification une fois terminée.',
              style: GoogleFonts.inter(
                fontSize: 13,
                color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
