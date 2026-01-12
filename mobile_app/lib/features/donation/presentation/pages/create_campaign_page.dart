import 'dart:io';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';
import '../../../../core/services/donation_api_service.dart';

class CreateCampaignPage extends StatefulWidget {
  const CreateCampaignPage({super.key});

  @override
  State<CreateCampaignPage> createState() => _CreateCampaignPageState();
}

class _CreateCampaignPageState extends State<CreateCampaignPage> {
  final _formKey = GlobalKey<FormState>();
  final DonationApiService _api = DonationApiService();
  final ImagePicker _picker = ImagePicker();

  // Form Fields
  String _title = '';
  String _description = '';
  double? _targetAmount;
  String _currency = 'XOF';
  bool _allowRecurring = true;
  bool _allowAnonymous = true;

  // Media
  File? _imageFile;
  File? _videoFile;
  VideoPlayerController? _videoController;

  // Dynamic Schema
  final List<Map<String, dynamic>> _dynamicFields = [
    {'name': 'full_name', 'label': 'Nom complet', 'type': 'text', 'required': true, 'options': <String>[]},
    {'name': 'email', 'label': 'Email', 'type': 'email', 'required': true, 'options': <String>[]},
  ];

  bool _loading = false;

  @override
  void dispose() {
    _videoController?.dispose();
    super.dispose();
  }

  // === Media Handling ===

  Future<void> _pickImage() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      setState(() {
        _imageFile = File(image.path);
      });
    }
  }

  Future<void> _pickVideo() async {
    final XFile? video = await _picker.pickVideo(source: ImageSource.gallery);
    if (video != null) {
      final file = File(video.path);
      // Basic size check (50MB)
      if (await file.length() > 50 * 1024 * 1024) {
        if (mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Vidéo trop volumineuse (Max 50MB)')));
        return;
      }

      setState(() {
        _videoFile = file;
        _videoController?.dispose();
        _videoController = VideoPlayerController.file(file)..initialize().then((_) => setState(() {}));
      });
    }
  }

  // === Dynamic Form Logic ===

  void _addField() {
    showDialog(
      context: context,
      builder: (context) => _AddFieldDialog(
        onAdd: (field) {
          setState(() {
            _dynamicFields.add(field);
          });
        },
      ),
    );
  }

  void _removeField(int index) {
    setState(() {
      _dynamicFields.removeAt(index);
    });
  }

  // === Submit ===

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    _formKey.currentState!.save();

    setState(() => _loading = true);

    try {
      String imageUrl = '';
      String videoUrl = '';

      // Upload Media
      if (_imageFile != null) {
        imageUrl = await _api.uploadMedia(_imageFile);
      }
      if (_videoFile != null) {
        videoUrl = await _api.uploadMedia(_videoFile);
      }

      // Build Schema
      // Ensure names are clean
      final schema = _dynamicFields.map((f) {
        return {
          'name': f['label'].toString().toLowerCase().replaceAll(RegExp(r'\s+'), '_'),
          'label': f['label'],
          'type': f['type'],
          'required': f['required'],
          'options': f['options'],
        };
      }).toList();

      final payload = {
        'title': _title,
        'description': _description,
        'target_amount': _targetAmount ?? 0,
        'currency': _currency,
        'image_url': imageUrl,
        'video_url': videoUrl,
        'allow_recurring': _allowRecurring,
        'allow_anonymous': _allowAnonymous,
        'form_schema': schema,
      };

      await _api.createCampaign(payload);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Campagne créée avec succès! 🚀')));
        Navigator.pop(context, true); // Return success
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Erreur: $e')));
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1a1a2e),
      appBar: AppBar(
        title: const Text('Créer une Cagnotte'),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Title
              TextFormField(
                style: const TextStyle(color: Colors.white),
                decoration: _inputDecoration('Titre de la campagne'),
                validator: (v) => v!.isEmpty ? 'Requis' : null,
                onSaved: (v) => _title = v!,
              ),
              const SizedBox(height: 16),
              
              // Description
              TextFormField(
                style: const TextStyle(color: Colors.white),
                decoration: _inputDecoration('Description & Histoire'),
                maxLines: 4,
                validator: (v) => v!.isEmpty ? 'Requis' : null,
                onSaved: (v) => _description = v!,
              ),
              const SizedBox(height: 16),

              // Target & Currency
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      style: const TextStyle(color: Colors.white),
                      decoration: _inputDecoration('Objectif (XOF)'),
                      keyboardType: TextInputType.number,
                      onSaved: (v) => _targetAmount = v!.isEmpty ? 0 : double.tryParse(v),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.05),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: Colors.white.withOpacity(0.1)),
                    ),
                    child: DropdownButton<String>(
                      value: _currency,
                      dropdownColor: const Color(0xFF1a1a2e),
                      underline: const SizedBox(),
                      style: const TextStyle(color: Colors.white),
                      items: ['XOF', 'USD', 'EUR'].map((c) => DropdownMenuItem(value: c, child: Text(c))).toList(),
                      onChanged: (v) => setState(() => _currency = v!),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // Media Section
              const Text('Médias', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
              const SizedBox(height: 12),
              
              // Image Picker
              GestureDetector(
                onTap: _pickImage,
                child: Container(
                  height: 150,
                  width: double.infinity,
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.white.withOpacity(0.1), style: BorderStyle.solid),
                    image: _imageFile != null ? DecorationImage(image: FileImage(_imageFile!), fit: BoxFit.cover) : null,
                  ),
                  child: _imageFile == null
                      ? const Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(Icons.image, color: Colors.white54, size: 40),
                            SizedBox(height: 8),
                            Text('Ajouter une image de couverture', style: TextStyle(color: Colors.white54)),
                          ],
                        )
                      : null,
                ),
              ),
              const SizedBox(height: 12),

              // Video Picker
              GestureDetector(
                onTap: _pickVideo,
                child: Container(
                  height: 150,
                  width: double.infinity,
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.white.withOpacity(0.1)),
                  ),
                  child: _videoFile != null
                      ? Stack(
                          alignment: Alignment.center,
                          children: [
                            if (_videoController != null && _videoController!.value.isInitialized)
                              AspectRatio(
                                aspectRatio: _videoController!.value.aspectRatio,
                                child: VideoPlayer(_videoController!),
                              ),
                            const Icon(Icons.play_circle_fill, color: Colors.white, size: 50),
                            Positioned(
                                top: 8,
                                right: 8,
                                child: IconButton(
                                  icon: const Icon(Icons.close, color: Colors.red),
                                  onPressed: () => setState(() => _videoFile = null),
                                ))
                          ],
                        )
                      : const Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(Icons.videocam, color: Colors.white54, size: 40),
                            SizedBox(height: 8),
                            Text('Ajouter une vidéo (Optionnel)', style: TextStyle(color: Colors.white54)),
                          ],
                        ),
                ),
              ),
              const SizedBox(height: 24),

              // Dynamic Form Builder
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text('Formulaire Donateur', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
                  TextButton.icon(
                    onPressed: _addField,
                    icon: const Icon(Icons.add, size: 18),
                    label: const Text('Ajouter champ'),
                  ),
                ],
              ),
              
              ListView.separated(
                physics: const NeverScrollableScrollPhysics(),
                shrinkWrap: true,
                itemCount: _dynamicFields.length,
                separatorBuilder: (c, i) => const SizedBox(height: 8),
                itemBuilder: (context, index) {
                  final field = _dynamicFields[index];
                  return Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.05),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      children: [
                        Expanded(child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(field['label'], style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                            Text('${field['type']} • ${field['required'] ? "Obligatoire" : "Optionnel"}', style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12)),
                          ],
                        )),
                        IconButton(
                          icon: const Icon(Icons.delete, color: Colors.red),
                          onPressed: () => _removeField(index),
                        ),
                      ],
                    ),
                  );
                },
              ),
              const SizedBox(height: 32),

              // Submit Button
              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: _loading ? null : _submit,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366f1),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  ),
                  child: _loading 
                      ? const CircularProgressIndicator(color: Colors.white)
                      : const Text('Lancer la campagne 🚀', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  InputDecoration _inputDecoration(String label) {
    return InputDecoration(
      labelText: label,
      labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
      filled: true,
      fillColor: Colors.white.withOpacity(0.05),
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
      enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
      focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFF6366f1))),
    );
  }
}

class _AddFieldDialog extends StatefulWidget {
  final Function(Map<String, dynamic>) onAdd;
  const _AddFieldDialog({required this.onAdd});
  @override
  State<_AddFieldDialog> createState() => _AddFieldDialogState();
}

class _AddFieldDialogState extends State<_AddFieldDialog> {
  String _label = '';
  String _type = 'text';
  bool _required = false;
  String _optionsText = '';

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: const Color(0xFF1a1a2e),
      title: const Text('Nouveau champ', style: TextStyle(color: Colors.white)),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextField(
            style: const TextStyle(color: Colors.white),
            decoration: const InputDecoration(labelText: 'Nom du champ', labelStyle: TextStyle(color: Colors.white54)),
            onChanged: (v) => _label = v,
          ),
          const SizedBox(height: 12),
          DropdownButtonFormField<String>(
            value: _type,
            dropdownColor: const Color(0xFF16213e),
            style: const TextStyle(color: Colors.white),
            decoration: const InputDecoration(labelText: 'Type', labelStyle: TextStyle(color: Colors.white54)),
            items: ['text', 'email', 'phone', 'number', 'select', 'checkbox'].map((t) => DropdownMenuItem(value: t, child: Text(t))).toList(),
            onChanged: (v) => setState(() => _type = v!),
          ),
          if (_type == 'select')
             TextField(
              style: const TextStyle(color: Colors.white),
              decoration: const InputDecoration(labelText: 'Options (virgules)', labelStyle: TextStyle(color: Colors.white54)),
              onChanged: (v) => _optionsText = v,
            ),
          const SizedBox(height: 12),
          SwitchListTile(
            title: const Text('Obligatoire', style: TextStyle(color: Colors.white)),
            value: _required,
            onChanged: (v) => setState(() => _required = v),
          )
        ],
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: const Text('Annuler')),
        TextButton(
          onPressed: () {
            if (_label.isNotEmpty) {
              widget.onAdd({
                'label': _label,
                'type': _type,
                'required': _required,
                'options': _optionsText.split(',').map((e) => e.trim()).where((e) => e.isNotEmpty).toList(),
              });
              Navigator.pop(context);
            }
          },
          child: const Text('Ajouter'),
        ),
      ],
    );
  }
}
