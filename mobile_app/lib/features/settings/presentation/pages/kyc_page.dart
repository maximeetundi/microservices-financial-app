import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'dart:io';
import 'package:image_picker/image_picker.dart';

import '../../../../core/widgets/glass_container.dart';
import '../../../../core/services/auth_api_service.dart';

/// KYC Verification Page - Binance Style
class KycPage extends StatefulWidget {
  const KycPage({super.key});

  @override
  State<KycPage> createState() => _KycPageState();
}

class _KycPageState extends State<KycPage> {
  final AuthApiService _authApi = AuthApiService();
  final ImagePicker _picker = ImagePicker();
  
  String _kycStatus = 'none'; // none, pending, submitted, verified, rejected
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

  int get _currentStep {
    if (_documentStatus['identity']!['status'] == 'required') return 1;
    if (_documentStatus['selfie']!['status'] == 'required') return 2;
    if (_documentStatus['address']!['status'] == 'required') return 3;
    return 3;
  }

  Future<void> _loadKycStatus() async {
    try {
      final profile = await _authApi.getProfile();
      
      setState(() {
        // Handle null or empty status as 'none'
        final status = profile['kyc_status'] ?? 'none';
        _kycStatus = (status == '' || status == null) ? 'none' : status;
        _kycLevel = profile['kyc_level'] ?? 0;
        _isLoading = false;
      });
      
      // Try to load documents
      try {
        final kycDocs = await _authApi.getKYCDocuments();
        setState(() {
          _documents = List<Map<String, dynamic>>.from(kycDocs['documents'] ?? []);
          for (var doc in _documents) {
            final type = doc['type'];
            if (_documentStatus.containsKey(type)) {
              _documentStatus[type] = {
                'status': doc['status'] ?? 'pending',
                'label': _getStatusLabel(doc['status']),
              };
            }
          }
        });
      } catch (e) {
        debugPrint('KYC documents not available: $e');
      }
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

  Future<void> _takePhoto() async {
    if (_selectedDocType == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Veuillez sélectionner un type de document')),
      );
      return;
    }

    final XFile? image = await _picker.pickImage(
      source: ImageSource.camera,
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
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Document envoyé avec succès!'),
            backgroundColor: Colors.green,
          ),
        );
      }

      setState(() {
        _documentStatus[_selectedDocType!] = {'status': 'pending', 'label': 'En cours'};
        _kycStatus = 'pending';
        _selectedFile = null;
        _selectedDocType = null;
      });

      await _loadKycStatus();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isUploading = false);
      }
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
                  const SizedBox(height: 20),
                  _buildProgressSteps(isDark),
                  const SizedBox(height: 20),
                  _buildStatusCard(isDark),
                  const SizedBox(height: 24),
                  _buildDocumentsSection(isDark),
                  if (_kycStatus != 'verified') ...[
                    const SizedBox(height: 24),
                    _buildUploadSection(isDark),
                  ],
                  const SizedBox(height: 24),
                  _buildInfoBox(isDark),
                  const SizedBox(height: 24),
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
            icon: Icon(Icons.arrow_back_ios_new, size: 18,
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
                '🔐 Vérification d\'identité',
                style: GoogleFonts.inter(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                'Validez votre identité pour toutes les fonctionnalités',
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

  Widget _buildProgressSteps(bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          _buildStep(1, 'Identité', _currentStep >= 1, _documentStatus['identity']!['status'] == 'approved', isDark),
          _buildStepLine(_documentStatus['identity']!['status'] == 'approved', isDark),
          _buildStep(2, 'Selfie', _currentStep >= 2, _documentStatus['selfie']!['status'] == 'approved', isDark),
          _buildStepLine(_documentStatus['selfie']!['status'] == 'approved', isDark),
          _buildStep(3, 'Domicile', _currentStep >= 3, _documentStatus['address']!['status'] == 'approved', isDark),
        ],
      ),
    );
  }

  Widget _buildStep(int number, String label, bool isActive, bool isCompleted, bool isDark) {
    return Column(
      children: [
        Container(
          width: 40,
          height: 40,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: isCompleted 
                ? const Color(0xFF22C55E)
                : isActive 
                    ? const Color(0xFF6366F1).withOpacity(0.2)
                    : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
            border: Border.all(
              color: isCompleted 
                  ? const Color(0xFF22C55E)
                  : isActive 
                      ? const Color(0xFF6366F1)
                      : (isDark ? Colors.white.withOpacity(0.2) : const Color(0xFFCBD5E1)),
              width: 2,
            ),
          ),
          child: Center(
            child: isCompleted
                ? const Icon(Icons.check, color: Colors.white, size: 20)
                : Text(
                    '$number',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: isActive 
                          ? const Color(0xFF6366F1)
                          : (isDark ? const Color(0xFF94A3B8) : const Color(0xFF94A3B8)),
                    ),
                  ),
          ),
        ),
        const SizedBox(height: 6),
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 11,
            color: isActive || isCompleted
                ? (isDark ? Colors.white : const Color(0xFF1E293B))
                : const Color(0xFF94A3B8),
          ),
        ),
      ],
    );
  }

  Widget _buildStepLine(bool isActive, bool isDark) {
    return Container(
      width: 40,
      height: 2,
      margin: const EdgeInsets.only(bottom: 20, left: 4, right: 4),
      decoration: BoxDecoration(
        color: isActive 
            ? const Color(0xFF22C55E)
            : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
        borderRadius: BorderRadius.circular(1),
      ),
    );
  }

  Widget _buildStatusCard(bool isDark) {
    final statusConfig = _getStatusConfig();
    
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: statusConfig['bgColor'],
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: statusConfig['borderColor']),
      ),
      child: Row(
        children: [
          Text(statusConfig['icon'], style: const TextStyle(fontSize: 36)),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  statusConfig['title'],
                  style: GoogleFonts.inter(
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  statusConfig['description'],
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              color: statusConfig['badgeColor'],
              borderRadius: BorderRadius.circular(20),
            ),
            child: Text(
              statusConfig['badge'],
              style: GoogleFonts.inter(
                fontSize: 10,
                fontWeight: FontWeight.w600,
                color: statusConfig['badgeTextColor'],
              ),
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
          'title': 'Identité vérifiée',
          'description': 'Accès complet à toutes les fonctionnalités.',
          'badge': 'VÉRIFIÉ',
          'bgColor': isDark 
              ? const Color(0xFF22C55E).withOpacity(0.1)
              : const Color(0xFFF0FDF4),
          'borderColor': isDark
              ? const Color(0xFF22C55E).withOpacity(0.2)
              : const Color(0xFFBBF7D0),
          'badgeColor': const Color(0xFF22C55E).withOpacity(0.2),
          'badgeTextColor': const Color(0xFF22C55E),
        };
      case 'rejected':
        return {
          'icon': '❌',
          'title': 'Vérification refusée',
          'description': 'Veuillez soumettre de nouveaux documents.',
          'badge': 'REFUSÉ',
          'bgColor': isDark 
              ? const Color(0xFFEF4444).withOpacity(0.1)
              : const Color(0xFFFEF2F2),
          'borderColor': isDark
              ? const Color(0xFFEF4444).withOpacity(0.2)
              : const Color(0xFFFECACA),
          'badgeColor': const Color(0xFFEF4444).withOpacity(0.2),
          'badgeTextColor': const Color(0xFFEF4444),
        };
      case 'pending':
      case 'submitted':
        return {
          'icon': '⏳',
          'title': 'Documents en cours de vérification',
          'description': 'Réponse sous 24-48h.',
          'badge': 'EN COURS',
          'bgColor': isDark 
              ? const Color(0xFFF97316).withOpacity(0.1)
              : const Color(0xFFFFF7ED),
          'borderColor': isDark
              ? const Color(0xFFF97316).withOpacity(0.2)
              : const Color(0xFFFED7AA),
          'badgeColor': const Color(0xFFF97316).withOpacity(0.2),
          'badgeTextColor': const Color(0xFFF97316),
        };
      default: // 'none' or empty
        return {
          'icon': '📝',
          'title': 'Vérification non commencée',
          'description': 'Soumettez vos documents pour débloquer toutes les fonctionnalités.',
          'badge': 'NON VÉRIFIÉ',
          'bgColor': isDark 
              ? Colors.white.withOpacity(0.05)
              : const Color(0xFFF8FAFC),
          'borderColor': isDark
              ? Colors.white.withOpacity(0.1)
              : const Color(0xFFE2E8F0),
          'badgeColor': const Color(0xFF94A3B8).withOpacity(0.2),
          'badgeTextColor': const Color(0xFF64748B),
        };
    }
  }

  Widget _buildDocumentsSection(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '📄 DOCUMENTS REQUIS',
          style: GoogleFonts.inter(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.2,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 12),
        _buildDocumentCard(
          '🪪',
          'Pièce d\'identité',
          'Passeport, carte d\'identité ou permis',
          'identity',
          ['Photo claire', 'Non expiré', 'Coins visibles'],
          isDark,
        ),
        _buildDocumentCard(
          '🤳',
          'Selfie avec document',
          'Photo de vous tenant votre pièce d\'identité',
          'selfie',
          ['Visage visible', 'Document lisible', 'Bonne lumière'],
          isDark,
        ),
        _buildDocumentCard(
          '🏠',
          'Justificatif de domicile',
          'Facture ou relevé bancaire récent',
          'address',
          ['Moins de 3 mois', 'Adresse complète', 'Nom visible'],
          isDark,
        ),
      ],
    );
  }

  Widget _buildDocumentCard(String emoji, String title, String subtitle, String type, List<String> requirements, bool isDark) {
    final status = _documentStatus[type]!;
    final isSelected = _selectedDocType == type;
    final isApproved = status['status'] == 'approved';
    
    return GestureDetector(
      onTap: isApproved ? null : () {
        setState(() {
          _selectedDocType = type;
          _selectedFile = null;
        });
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected
              ? const Color(0xFF6366F1).withOpacity(0.1)
              : (isDark ? Colors.white.withOpacity(0.03) : Colors.white),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isSelected
                ? const Color(0xFF6366F1)
                : (isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0)),
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(emoji, style: const TextStyle(fontSize: 28)),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title,
                        style: GoogleFonts.inter(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
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
                    color: _getDocStatusColor(status['status'], isDark),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    status['label'],
                    style: GoogleFonts.inter(
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                      color: _getDocStatusTextColor(status['status']),
                    ),
                  ),
                ),
              ],
            ),
            if (isSelected) ...[
              const SizedBox(height: 12),
              Wrap(
                spacing: 8,
                runSpacing: 4,
                children: requirements.map((req) => Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: const Color(0xFF22C55E).withOpacity(0.1),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Text(
                    '✓ $req',
                    style: GoogleFonts.inter(
                      fontSize: 10,
                      color: const Color(0xFF22C55E),
                    ),
                  ),
                )).toList(),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Color _getDocStatusColor(String status, bool isDark) {
    switch (status) {
      case 'approved':
        return const Color(0xFF22C55E).withOpacity(0.2);
      case 'rejected':
        return const Color(0xFFEF4444).withOpacity(0.2);
      case 'pending':
        return const Color(0xFFF97316).withOpacity(0.2);
      default:
        return isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFF1F5F9);
    }
  }

  Color _getDocStatusTextColor(String status) {
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
    if (_selectedDocType == null) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF6366F1).withOpacity(0.05),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFF6366F1).withOpacity(0.1)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '📤 TÉLÉCHARGER ${_getDocTypeName(_selectedDocType!).toUpperCase()}',
            style: GoogleFonts.inter(
              fontSize: 12,
              fontWeight: FontWeight.w600,
              letterSpacing: 1.2,
              color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 16),
          
          // Preview or Upload Area
          if (_selectedFile != null) ...[
            Stack(
              children: [
                ClipRRect(
                  borderRadius: BorderRadius.circular(12),
                  child: Image.file(
                    _selectedFile!,
                    width: double.infinity,
                    height: 200,
                    fit: BoxFit.cover,
                  ),
                ),
                Positioned(
                  top: 8,
                  right: 8,
                  child: GestureDetector(
                    onTap: () => setState(() => _selectedFile = null),
                    child: Container(
                      padding: const EdgeInsets.all(6),
                      decoration: const BoxDecoration(
                        color: Colors.red,
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(Icons.close, color: Colors.white, size: 16),
                    ),
                  ),
                ),
              ],
            ),
          ] else ...[
            Row(
              children: [
                Expanded(
                  child: GestureDetector(
                    onTap: _pickDocument,
                    child: Container(
                      padding: const EdgeInsets.all(24),
                      decoration: BoxDecoration(
                        color: isDark ? Colors.white.withOpacity(0.02) : Colors.white,
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: const Color(0xFF6366F1).withOpacity(0.3),
                          width: 2,
                          style: BorderStyle.solid,
                        ),
                      ),
                      child: Column(
                        children: [
                          const Icon(Icons.photo_library, size: 32, color: Color(0xFF6366F1)),
                          const SizedBox(height: 8),
                          Text(
                            'Galerie',
                            style: GoogleFonts.inter(
                              fontSize: 12,
                              color: isDark ? Colors.white : const Color(0xFF1E293B),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: GestureDetector(
                    onTap: _takePhoto,
                    child: Container(
                      padding: const EdgeInsets.all(24),
                      decoration: BoxDecoration(
                        color: isDark ? Colors.white.withOpacity(0.02) : Colors.white,
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: const Color(0xFF6366F1).withOpacity(0.3),
                          width: 2,
                          style: BorderStyle.solid,
                        ),
                      ),
                      child: Column(
                        children: [
                          const Icon(Icons.camera_alt, size: 32, color: Color(0xFF6366F1)),
                          const SizedBox(height: 8),
                          Text(
                            'Caméra',
                            style: GoogleFonts.inter(
                              fontSize: 12,
                              color: isDark ? Colors.white : const Color(0xFF1E293B),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ],
          const SizedBox(height: 16),
          
          // Upload Button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: (_selectedFile != null && !_isUploading) ? _uploadDocument : null,
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
      ),
    );
  }

  String _getDocTypeName(String type) {
    switch (type) {
      case 'identity': return "Pièce d'identité";
      case 'selfie': return 'Selfie';
      case 'address': return 'Justificatif';
      default: return 'Document';
    }
  }

  Widget _buildInfoBox(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark 
            ? const Color(0xFF3B82F6).withOpacity(0.1)
            : const Color(0xFFEFF6FF),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isDark 
              ? const Color(0xFF3B82F6).withOpacity(0.2)
              : const Color(0xFFBFDBFE),
        ),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('💡', style: TextStyle(fontSize: 24)),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Conseils pour une vérification rapide',
                  style: GoogleFonts.inter(
                    fontSize: 13,
                    fontWeight: FontWeight.w600,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                const SizedBox(height: 6),
                Text(
                  '• Documents lisibles et non flous\n• Évitez les reflets et ombres\n• Vérification sous 24-48 heures',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    height: 1.4,
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
}
