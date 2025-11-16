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
    this.preferences = const UserPreferences(),
  });

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