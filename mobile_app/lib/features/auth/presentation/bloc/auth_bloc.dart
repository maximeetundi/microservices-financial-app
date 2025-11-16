import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../domain/entities/user.dart';
import '../../domain/usecases/sign_in_usecase.dart';
import '../../domain/usecases/sign_up_usecase.dart';
import '../../domain/usecases/sign_out_usecase.dart';
import '../../domain/usecases/check_auth_status_usecase.dart';
import '../../domain/usecases/biometric_auth_usecase.dart';

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

  const SignUpEvent({
    required this.email,
    required this.password,
    required this.firstName,
    required this.lastName,
    this.phoneNumber,
  });

  @override
  List<Object?> get props => [email, password, firstName, lastName, phoneNumber];
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
  final SignInUseCase _signInUseCase;
  final SignUpUseCase _signUpUseCase;
  final SignOutUseCase _signOutUseCase;
  final CheckAuthStatusUseCase _checkAuthStatusUseCase;
  final BiometricAuthUseCase _biometricAuthUseCase;

  AuthBloc({
    required SignInUseCase signInUseCase,
    required SignUpUseCase signUpUseCase,
    required SignOutUseCase signOutUseCase,
    required CheckAuthStatusUseCase checkAuthStatusUseCase,
    required BiometricAuthUseCase biometricAuthUseCase,
  })  : _signInUseCase = signInUseCase,
        _signUpUseCase = signUpUseCase,
        _signOutUseCase = signOutUseCase,
        _checkAuthStatusUseCase = checkAuthStatusUseCase,
        _biometricAuthUseCase = biometricAuthUseCase,
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
      final result = await _signInUseCase(SignInParams(
        email: event.email,
        password: event.password,
        totpCode: event.totpCode,
        rememberMe: event.rememberMe,
      ));

      result.fold(
        (failure) => emit(AuthErrorState(message: failure.message)),
        (authResult) {
          if (authResult.requires2FA && event.totpCode == null) {
            emit(Auth2FARequiredState(tempToken: authResult.tempToken!));
          } else {
            emit(AuthenticatedState(
              user: authResult.user!,
              token: authResult.token!,
            ));
          }
        },
      );
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onSignUp(SignUpEvent event, Emitter<AuthState> emit) async {
    emit(const AuthLoadingState());

    try {
      final result = await _signUpUseCase(SignUpParams(
        email: event.email,
        password: event.password,
        firstName: event.firstName,
        lastName: event.lastName,
        phoneNumber: event.phoneNumber,
      ));

      result.fold(
        (failure) => emit(AuthErrorState(message: failure.message)),
        (authResult) => emit(AuthenticatedState(
          user: authResult.user!,
          token: authResult.token!,
        )),
      );
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onSignOut(SignOutEvent event, Emitter<AuthState> emit) async {
    emit(const AuthLoadingState());

    try {
      await _signOutUseCase(NoParams());
      emit(const UnauthenticatedState());
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onCheckAuthStatus(
    CheckAuthStatusEvent event,
    Emitter<AuthState> emit,
  ) async {
    try {
      final result = await _checkAuthStatusUseCase(NoParams());

      result.fold(
        (failure) => emit(const UnauthenticatedState()),
        (authResult) {
          if (authResult.user != null && authResult.token != null) {
            emit(AuthenticatedState(
              user: authResult.user!,
              token: authResult.token!,
            ));
          } else {
            emit(const UnauthenticatedState());
          }
        },
      );
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
      final result = await _biometricAuthUseCase(NoParams());

      result.fold(
        (failure) => emit(AuthErrorState(message: failure.message)),
        (authResult) => emit(AuthenticatedState(
          user: authResult.user!,
          token: authResult.token!,
        )),
      );
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onForgotPassword(
    ForgotPasswordEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoadingState());

    try {
      // Implement forgot password logic
      // For now, just emit success
      emit(PasswordResetEmailSentState(email: event.email));
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onResetPassword(
    ResetPasswordEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoadingState());

    try {
      // Implement reset password logic
      // For now, just emit success
      emit(const PasswordResetSuccessState());
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }
}

// Use Case Parameters
class SignInParams extends Equatable {
  final String email;
  final String password;
  final String? totpCode;
  final bool rememberMe;

  const SignInParams({
    required this.email,
    required this.password,
    this.totpCode,
    this.rememberMe = false,
  });

  @override
  List<Object?> get props => [email, password, totpCode, rememberMe];
}

class SignUpParams extends Equatable {
  final String email;
  final String password;
  final String firstName;
  final String lastName;
  final String? phoneNumber;

  const SignUpParams({
    required this.email,
    required this.password,
    required this.firstName,
    required this.lastName,
    this.phoneNumber,
  });

  @override
  List<Object?> get props => [email, password, firstName, lastName, phoneNumber];
}

class NoParams extends Equatable {
  @override
  List<Object> get props => [];
}