import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String email;
  final String? firstName;
  final String? lastName;
  final String? phoneNumber;
  final String? profilePictureUrl;
  final DateTime createdAt;
  final DateTime? lastLoginAt;
  final bool isEmailVerified;
  final bool isPhoneVerified;
  final bool isTwoFactorEnabled;
  final String kycLevel;
  final bool hasPin;
  final UserPreferences preferences;

  const User({
    required this.id,
    required this.email,
    this.firstName,
    this.lastName,
    this.phoneNumber,
    this.profilePictureUrl,
    required this.createdAt,
    this.lastLoginAt,
    this.isEmailVerified = false,
    this.isPhoneVerified = false,
    this.isTwoFactorEnabled = false,
    this.kycLevel = 'none',
    this.hasPin = false,
    this.preferences = const UserPreferences(),
  });
  
  factory User.fromJson(Map<String, dynamic> json) {
    // Handle kycLevel which can be int or String from API
    String kycLevelValue = 'none';
    if (json['kyc_level'] != null) {
      if (json['kyc_level'] is int) {
        kycLevelValue = json['kyc_level'].toString();
      } else {
        kycLevelValue = json['kyc_level'].toString();
      }
    }
    
    return User(
      id: json['id'] ?? json['user_id'] ?? '',
      email: json['email'] ?? '',
      firstName: json['first_name'] ?? json['firstName'],
      lastName: json['last_name'] ?? json['lastName'],
      phoneNumber: json['phone_number'] ?? json['phone'],
      profilePictureUrl: json['profile_picture_url'] ?? json['avatar'],
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : DateTime.now(),
      lastLoginAt: json['last_login_at'] != null 
          ? DateTime.parse(json['last_login_at']) 
          : null,
      isEmailVerified: json['is_email_verified'] ?? json['email_verified'] ?? false,
      isPhoneVerified: json['is_phone_verified'] ?? json['phone_verified'] ?? false,
      isTwoFactorEnabled: json['is_two_factor_enabled'] ?? json['two_fa_enabled'] ?? json['two_factor_enabled'] ?? false,
      kycLevel: kycLevelValue,
      hasPin: json['has_pin'] ?? false,
      preferences: json['preferences'] != null 
          ? UserPreferences.fromJson(json['preferences']) 
          : const UserPreferences(),
    );
  }
  
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'first_name': firstName,
      'last_name': lastName,
      'phone_number': phoneNumber,
      'profile_picture_url': profilePictureUrl,
      'created_at': createdAt.toIso8601String(),
      'last_login_at': lastLoginAt?.toIso8601String(),
      'is_email_verified': isEmailVerified,
      'is_phone_verified': isPhoneVerified,
      'is_two_factor_enabled': isTwoFactorEnabled,
      'kyc_level': kycLevel,
    };
  }

  String get fullName {
    if (firstName != null && lastName != null) {
      return '$firstName $lastName';
    } else if (firstName != null) {
      return firstName!;
    } else if (lastName != null) {
      return lastName!;
    }
    return email.split('@').first;
  }

  String get initials {
    if (firstName != null && lastName != null) {
      return '${firstName!.substring(0, 1)}${lastName!.substring(0, 1)}';
    } else if (firstName != null) {
      return firstName!.substring(0, 1).toUpperCase();
    }
    return email.substring(0, 1).toUpperCase();
  }

  User copyWith({
    String? id,
    String? email,
    String? firstName,
    String? lastName,
    String? phoneNumber,
    String? profilePictureUrl,
    DateTime? createdAt,
    DateTime? lastLoginAt,
    bool? isEmailVerified,
    bool? isPhoneVerified,
    bool? isTwoFactorEnabled,
    String? kycLevel,
    bool? hasPin,
    UserPreferences? preferences,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      phoneNumber: phoneNumber ?? this.phoneNumber,
      profilePictureUrl: profilePictureUrl ?? this.profilePictureUrl,
      createdAt: createdAt ?? this.createdAt,
      lastLoginAt: lastLoginAt ?? this.lastLoginAt,
      isEmailVerified: isEmailVerified ?? this.isEmailVerified,
      isPhoneVerified: isPhoneVerified ?? this.isPhoneVerified,
      isTwoFactorEnabled: isTwoFactorEnabled ?? this.isTwoFactorEnabled,
      kycLevel: kycLevel ?? this.kycLevel,
      hasPin: hasPin ?? this.hasPin,
      preferences: preferences ?? this.preferences,
    );
  }

  @override
  List<Object?> get props => [
        id,
        email,
        firstName,
        lastName,
        phoneNumber,
        profilePictureUrl,
        createdAt,
        lastLoginAt,
        isEmailVerified,
        isPhoneVerified,
        isTwoFactorEnabled,
        kycLevel,
        hasPin,
        preferences,
      ];
}

class UserPreferences extends Equatable {
  final String preferredCurrency;
  final String language;
  final String timezone;
  final bool notificationsEnabled;
  final bool biometricsEnabled;
  final bool marketingEmailsEnabled;
  final bool pushNotificationsEnabled;
  final String theme; // 'light', 'dark', 'system'

  const UserPreferences({
    this.preferredCurrency = 'USD',
    this.language = 'en',
    this.timezone = 'UTC',
    this.notificationsEnabled = true,
    this.biometricsEnabled = false,
    this.marketingEmailsEnabled = false,
    this.pushNotificationsEnabled = true,
    this.theme = 'system',
  });
  
  factory UserPreferences.fromJson(Map<String, dynamic> json) {
    return UserPreferences(
      preferredCurrency: json['preferred_currency'] ?? 'USD',
      language: json['language'] ?? 'en',
      timezone: json['timezone'] ?? 'UTC',
      notificationsEnabled: json['notifications_enabled'] ?? true,
      biometricsEnabled: json['biometrics_enabled'] ?? false,
      marketingEmailsEnabled: json['marketing_emails_enabled'] ?? false,
      pushNotificationsEnabled: json['push_notifications_enabled'] ?? true,
      theme: json['theme'] ?? 'system',
    );
  }

  UserPreferences copyWith({
    String? preferredCurrency,
    String? language,
    String? timezone,
    bool? notificationsEnabled,
    bool? biometricsEnabled,
    bool? marketingEmailsEnabled,
    bool? pushNotificationsEnabled,
    String? theme,
  }) {
    return UserPreferences(
      preferredCurrency: preferredCurrency ?? this.preferredCurrency,
      language: language ?? this.language,
      timezone: timezone ?? this.timezone,
      notificationsEnabled: notificationsEnabled ?? this.notificationsEnabled,
      biometricsEnabled: biometricsEnabled ?? this.biometricsEnabled,
      marketingEmailsEnabled: marketingEmailsEnabled ?? this.marketingEmailsEnabled,
      pushNotificationsEnabled: pushNotificationsEnabled ?? this.pushNotificationsEnabled,
      theme: theme ?? this.theme,
    );
  }

  @override
  List<Object> get props => [
        preferredCurrency,
        language,
        timezone,
        notificationsEnabled,
        biometricsEnabled,
        marketingEmailsEnabled,
        pushNotificationsEnabled,
        theme,
      ];
}

class AuthResult extends Equatable {
  final User? user;
  final String? token;
  final String? refreshToken;
  final bool requires2FA;
  final String? tempToken;
  final DateTime? expiresAt;

  const AuthResult({
    this.user,
    this.token,
    this.refreshToken,
    this.requires2FA = false,
    this.tempToken,
    this.expiresAt,
  });

  @override
  List<Object?> get props => [
        user,
        token,
        refreshToken,
        requires2FA,
        tempToken,
        expiresAt,
      ];
}