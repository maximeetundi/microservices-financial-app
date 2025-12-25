import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:dio/dio.dart';

import '../../../../core/services/api_service.dart';
import '../../../../core/services/secure_storage_service.dart';
import '../../domain/entities/user.dart';

// Events
abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class SignInEvent extends AuthEvent {
  final String email;
  final String password;
  final String? totpCode;
  final bool rememberMe;

  const SignInEvent({
    required this.email,
    required this.password,
    this.totpCode,
    this.rememberMe = false,
  });

  @override
  List<Object?> get props => [email, password, totpCode, rememberMe];
}

class SignUpEvent extends AuthEvent {
  final String email;
  final String password;
  final String firstName;
  final String lastName;
  final String? phoneNumber;
  final String? dateOfBirth;
  final String? country;
  final String? currency;

  const SignUpEvent({
    required this.email,
    required this.password,
    required this.firstName,
    required this.lastName,
    this.phoneNumber,
    this.dateOfBirth,
    this.country,
    this.currency,
  });

  @override
  List<Object?> get props => [email, password, firstName, lastName, phoneNumber, dateOfBirth, country, currency];
}

class SignOutEvent extends AuthEvent {
  const SignOutEvent();
}

class CheckAuthStatusEvent extends AuthEvent {
  const CheckAuthStatusEvent();
}

class BiometricSignInEvent extends AuthEvent {
  const BiometricSignInEvent();
}

class ForgotPasswordEvent extends AuthEvent {
  final String email;

  const ForgotPasswordEvent({required this.email});

  @override
  List<Object> get props => [email];
}

class ResetPasswordEvent extends AuthEvent {
  final String token;
  final String newPassword;

  const ResetPasswordEvent({
    required this.token,
    required this.newPassword,
  });

  @override
  List<Object> get props => [token, newPassword];
}

// States
abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitialState extends AuthState {
  const AuthInitialState();
}

class AuthLoadingState extends AuthState {
  const AuthLoadingState();
}

class AuthenticatedState extends AuthState {
  final User user;
  final String token;

  const AuthenticatedState({
    required this.user,
    required this.token,
  });

  @override
  List<Object> get props => [user, token];
}

class UnauthenticatedState extends AuthState {
  const UnauthenticatedState();
}

class Auth2FARequiredState extends AuthState {
  final String tempToken;

  const Auth2FARequiredState({required this.tempToken});

  @override
  List<Object> get props => [tempToken];
}

class AuthErrorState extends AuthState {
  final String message;

  const AuthErrorState({required this.message});

  @override
  List<Object> get props => [message];
}

class PasswordResetEmailSentState extends AuthState {
  final String email;

  const PasswordResetEmailSentState({required this.email});

  @override
  List<Object> get props => [email];
}

class PasswordResetSuccessState extends AuthState {
  const PasswordResetSuccessState();
}

// BLoC
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final ApiService _apiService;
  final SecureStorageService _secureStorage;

  AuthBloc({
    required ApiService apiService,
    required SecureStorageService secureStorage,
  })  : _apiService = apiService,
        _secureStorage = secureStorage,
        super(const AuthInitialState()) {
    on<SignInEvent>(_onSignIn);
    on<SignUpEvent>(_onSignUp);
    on<SignOutEvent>(_onSignOut);
    on<CheckAuthStatusEvent>(_onCheckAuthStatus);
    on<BiometricSignInEvent>(_onBiometricSignIn);
    on<ForgotPasswordEvent>(_onForgotPassword);
    on<ResetPasswordEvent>(_onResetPassword);
  }

  Future<void> _onSignIn(SignInEvent event, Emitter<AuthState> emit) async {
    emit(const AuthLoadingState());

    try {
      final result = await _apiService.auth.login(
        event.email, 
        event.password,
        totpCode: event.totpCode,
      );
      
      // Check if 2FA is required
      if (result['requires_2fa'] == true && event.totpCode == null) {
        emit(Auth2FARequiredState(tempToken: result['temp_token'] ?? ''));
        return;
      }
      
      final userData = result['user'] as Map<String, dynamic>;
      final user = User.fromJson(userData);
      final token = result['access_token'] as String;
      
      // Save user ID
      await _secureStorage.saveUserId(user.id);
      
      emit(AuthenticatedState(user: user, token: token));
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        final data = e.response?.data;
        if (data != null && data is Map<String, dynamic> && data['requires_2fa'] == true) {
          emit(Auth2FARequiredState(tempToken: data['temp_token'] ?? ''));
          return;
        }
      }
      emit(AuthErrorState(message: _getErrorMessage(e)));
    } catch (e) {
      emit(AuthErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onSignUp(SignUpEvent event, Emitter<AuthState> emit) async {
    emit(const AuthLoadingState());

    try {
      final result = await _apiService.auth.register(
        email: event.email,
        password: event.password,
        firstName: event.firstName,
        lastName: event.lastName,
        phone: event.phoneNumber,
        dateOfBirth: event.dateOfBirth,
        country: event.country,
        currency: event.currency,
      );
      
      // Auto-login after registration
      final loginResult = await _apiService.auth.login(event.email, event.password);
      
      final userData = loginResult['user'] as Map<String, dynamic>;
      final user = User.fromJson(userData);
      final token = loginResult['access_token'] as String;
      
      await _secureStorage.saveUserId(user.id);
      
      emit(AuthenticatedState(user: user, token: token));
    } catch (e) {
      emit(AuthErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onSignOut(SignOutEvent event, Emitter<AuthState> emit) async {
    emit(const AuthLoadingState());

    try {
      await _apiService.auth.logout();
      await _secureStorage.clearAll();
      emit(const UnauthenticatedState());
    } catch (e) {
      // Even if logout fails on server, clear local data
      await _secureStorage.clearAll();
      emit(const UnauthenticatedState());
    }
  }

  Future<void> _onCheckAuthStatus(
    CheckAuthStatusEvent event,
    Emitter<AuthState> emit,
  ) async {
    try {
      final isAuthenticated = await _apiService.auth.isAuthenticated();
      
      if (isAuthenticated) {
        final profileData = await _apiService.auth.getProfile();
        final user = User.fromJson(profileData);
        
        emit(AuthenticatedState(user: user, token: ''));
      } else {
        emit(const UnauthenticatedState());
      }
    } catch (e) {
      emit(const UnauthenticatedState());
    }
  }

  Future<void> _onBiometricSignIn(
    BiometricSignInEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoadingState());

    try {
      // Check if biometric is available and enabled
      final isAvailable = await _secureStorage.isBiometricAvailable();
      final isEnabled = await _secureStorage.isBiometricEnabled();
      
      if (!isAvailable || !isEnabled) {
        emit(const AuthErrorState(message: 'L\'authentification biométrique n\'est pas disponible'));
        return;
      }
      
      // Authenticate with biometrics
      final authenticated = await _secureStorage.authenticateWithBiometrics(
        reason: 'Authentifiez-vous pour accéder à votre compte',
      );
      
      if (authenticated) {
        // Get stored session
        final hasSession = await _secureStorage.hasValidSession();
        
        if (hasSession) {
          final profileData = await _apiService.auth.getProfile();
          final user = User.fromJson(profileData);
          
          emit(AuthenticatedState(user: user, token: ''));
        } else {
          emit(const AuthErrorState(message: 'Session expirée, veuillez vous reconnecter'));
        }
      } else {
        emit(const AuthErrorState(message: 'Authentification biométrique échouée'));
      }
    } catch (e) {
      emit(AuthErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onForgotPassword(
    ForgotPasswordEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoadingState());

    try {
      await _apiService.auth.forgotPassword(event.email);
      emit(PasswordResetEmailSentState(email: event.email));
    } catch (e) {
      emit(AuthErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onResetPassword(
    ResetPasswordEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoadingState());

    try {
      // TODO: Implement reset password with token
      emit(const PasswordResetSuccessState());
    } catch (e) {
      emit(AuthErrorState(message: _getErrorMessage(e)));
    }
  }
  
  String _getErrorMessage(dynamic error) {
    if (error is Exception) {
      final message = error.toString();
      
      // Handle DioException messages for 401 Unauthorized
      if (message.contains('401') || message.contains('Unauthorized')) {
        return 'Email ou mot de passe incorrect.';
      }
      
      if (message.contains('Exception: ')) {
        return message.replaceFirst('Exception: ', '');
      }
      return message;
    }
    return 'Une erreur est survenue';
  }
}