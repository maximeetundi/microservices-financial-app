import 'package:flutter_contacts/flutter_contacts.dart';
import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';

/// Service for syncing and matching phone contacts
/// Used to display contact names in messaging (only if the number is saved in device contacts)
class ContactsSyncService {
  static final ContactsSyncService _instance = ContactsSyncService._internal();
  factory ContactsSyncService() => _instance;
  ContactsSyncService._internal();

  // Cache of normalized phone -> contact name
  Map<String, String> _phoneToName = {};
  // Cache of email -> contact name
  Map<String, String> _emailToName = {};
  
  bool _isInitialized = false;
  bool _hasPermission = false;

  /// Check if contacts permission is granted
  Future<bool> checkPermission() async {
    _hasPermission = await FlutterContacts.requestPermission();
    return _hasPermission;
  }

  /// Sync contacts from device
  /// Call this on app start or when user grants permission
  Future<void> syncContacts() async {
    if (!_hasPermission) {
      _hasPermission = await FlutterContacts.requestPermission();
      if (!_hasPermission) {
        debugPrint('ContactsSyncService: Permission denied');
        return;
      }
    }

    try {
      debugPrint('ContactsSyncService: Starting sync...');
      final contacts = await FlutterContacts.getContacts(
        withProperties: true,
        withPhoto: false,
      );

      _phoneToName.clear();
      _emailToName.clear();

      for (final contact in contacts) {
        final displayName = contact.displayName;
        if (displayName.isEmpty) continue;

        // Map all phone numbers for this contact
        for (final phone in contact.phones) {
          final normalized = _normalizePhone(phone.number);
          if (normalized.isNotEmpty) {
            _phoneToName[normalized] = displayName;
          }
        }

        // Map all emails for this contact
        for (final email in contact.emails) {
          final normalizedEmail = email.address.toLowerCase().trim();
          if (normalizedEmail.isNotEmpty) {
            _emailToName[normalizedEmail] = displayName;
          }
        }
      }

      _isInitialized = true;
      debugPrint('ContactsSyncService: Synced ${contacts.length} contacts, ${_phoneToName.length} phones, ${_emailToName.length} emails');
      
      // Cache for next app start
      await _saveToCache();
    } catch (e) {
      debugPrint('ContactsSyncService: Error syncing contacts: $e');
    }
  }

  /// Get contact name by phone number
  /// Returns null if not found in contacts
  String? getNameByPhone(String? phone) {
    if (phone == null || phone.isEmpty) return null;
    final normalized = _normalizePhone(phone);
    return _phoneToName[normalized];
  }

  /// Get contact name by email
  /// Returns null if not found in contacts
  String? getNameByEmail(String? email) {
    if (email == null || email.isEmpty) return null;
    return _emailToName[email.toLowerCase().trim()];
  }

  /// Get display name for a conversation
  /// Priority: contact name (if found) > fallback name > phone/email
  String getDisplayName({
    String? phone,
    String? email,
    String? fallbackName,
  }) {
    // First try phone
    final nameByPhone = getNameByPhone(phone);
    if (nameByPhone != null) return nameByPhone;

    // Then try email
    final nameByEmail = getNameByEmail(email);
    if (nameByEmail != null) return nameByEmail;

    // Fallback to provided name or phone/email
    if (fallbackName != null && fallbackName.isNotEmpty) {
      return fallbackName;
    }
    
    return phone ?? email ?? 'Utilisateur';
  }

  /// Normalize phone number (remove spaces, dashes, and keep only digits)
  String _normalizePhone(String phone) {
    // Remove all non-digit characters except +
    String normalized = phone.replaceAll(RegExp(r'[^\d+]'), '');
    
    // Remove leading zeros after country code
    if (normalized.startsWith('+')) {
      // Keep country code, remove leading zeros from rest
      final countryCodeEnd = normalized.indexOf(RegExp(r'\d'), 1);
      if (countryCodeEnd > 0) {
        final countryCode = normalized.substring(0, 3); // e.g., +33
        final rest = normalized.substring(3).replaceFirst(RegExp(r'^0+'), '');
        normalized = countryCode + rest;
      }
    } else {
      // Just remove leading zeros
      normalized = normalized.replaceFirst(RegExp(r'^0+'), '');
    }
    
    return normalized;
  }

  /// Save cache to SharedPreferences
  Future<void> _saveToCache() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('contacts_phone_cache', jsonEncode(_phoneToName));
      await prefs.setString('contacts_email_cache', jsonEncode(_emailToName));
    } catch (e) {
      debugPrint('ContactsSyncService: Error saving cache: $e');
    }
  }

  /// Load cache from SharedPreferences
  Future<void> loadFromCache() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final phoneCache = prefs.getString('contacts_phone_cache');
      final emailCache = prefs.getString('contacts_email_cache');
      
      if (phoneCache != null) {
        _phoneToName = Map<String, String>.from(jsonDecode(phoneCache));
      }
      if (emailCache != null) {
        _emailToName = Map<String, String>.from(jsonDecode(emailCache));
      }
      
      _isInitialized = _phoneToName.isNotEmpty || _emailToName.isNotEmpty;
      debugPrint('ContactsSyncService: Loaded from cache: ${_phoneToName.length} phones, ${_emailToName.length} emails');
    } catch (e) {
      debugPrint('ContactsSyncService: Error loading cache: $e');
    }
  }

  bool get isInitialized => _isInitialized;
  bool get hasPermission => _hasPermission;
}
