import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

// Events
abstract class NotificationsEvent extends Equatable {
  const NotificationsEvent();
  @override
  List<Object?> get props => [];
}

class LoadNotificationsEvent extends NotificationsEvent {}

class MarkNotificationReadEvent extends NotificationsEvent {
  final String notificationId;
  const MarkNotificationReadEvent(this.notificationId);
  @override
  List<Object?> get props => [notificationId];
}

class MarkAllNotificationsReadEvent extends NotificationsEvent {}

class DeleteNotificationEvent extends NotificationsEvent {
  final String notificationId;
  const DeleteNotificationEvent(this.notificationId);
  @override
  List<Object?> get props => [notificationId];
}

// States
abstract class NotificationsState extends Equatable {
  const NotificationsState();
  @override
  List<Object?> get props => [];
}

class NotificationsInitial extends NotificationsState {}

class NotificationsLoadingState extends NotificationsState {}

class NotificationsLoadedState extends NotificationsState {
  final List<NotificationItem> notifications;
  final int unreadCount;

  const NotificationsLoadedState({
    required this.notifications,
    required this.unreadCount,
  });

  @override
  List<Object?> get props => [notifications, unreadCount];
}

class NotificationsErrorState extends NotificationsState {
  final String message;
  const NotificationsErrorState(this.message);
  @override
  List<Object?> get props => [message];
}

// Notification Model
class NotificationItem {
  final String id;
  final String type;
  final String title;
  final String message;
  final DateTime createdAt;
  final bool isRead;

  NotificationItem({
    required this.id,
    required this.type,
    required this.title,
    required this.message,
    required this.createdAt,
    required this.isRead,
  });
}

// Bloc
class NotificationsBloc extends Bloc<NotificationsEvent, NotificationsState> {
  NotificationsBloc() : super(NotificationsInitial()) {
    on<LoadNotificationsEvent>(_onLoadNotifications);
    on<MarkNotificationReadEvent>(_onMarkRead);
    on<MarkAllNotificationsReadEvent>(_onMarkAllRead);
    on<DeleteNotificationEvent>(_onDelete);
  }

  final List<NotificationItem> _notifications = [
    NotificationItem(
      id: '1',
      type: 'transfer',
      title: 'Transfert reçu',
      message: 'Vous avez reçu 500 USD',
      createdAt: DateTime.now().subtract(const Duration(minutes: 10)),
      isRead: false,
    ),
    NotificationItem(
      id: '2',
      type: 'card',
      title: 'Paiement effectué',
      message: 'Paiement de 25€ chez Amazon',
      createdAt: DateTime.now().subtract(const Duration(hours: 2)),
      isRead: true,
    ),
  ];

  void _onLoadNotifications(LoadNotificationsEvent event, Emitter<NotificationsState> emit) {
    emit(NotificationsLoadingState());
    try {
      final unread = _notifications.where((n) => !n.isRead).length;
      emit(NotificationsLoadedState(notifications: _notifications, unreadCount: unread));
    } catch (e) {
      emit(NotificationsErrorState(e.toString()));
    }
  }

  void _onMarkRead(MarkNotificationReadEvent event, Emitter<NotificationsState> emit) {
    final index = _notifications.indexWhere((n) => n.id == event.notificationId);
    if (index != -1) {
      _notifications[index] = NotificationItem(
        id: _notifications[index].id,
        type: _notifications[index].type,
        title: _notifications[index].title,
        message: _notifications[index].message,
        createdAt: _notifications[index].createdAt,
        isRead: true,
      );
    }
    final unread = _notifications.where((n) => !n.isRead).length;
    emit(NotificationsLoadedState(notifications: _notifications, unreadCount: unread));
  }

  void _onMarkAllRead(MarkAllNotificationsReadEvent event, Emitter<NotificationsState> emit) {
    for (int i = 0; i < _notifications.length; i++) {
      _notifications[i] = NotificationItem(
        id: _notifications[i].id,
        type: _notifications[i].type,
        title: _notifications[i].title,
        message: _notifications[i].message,
        createdAt: _notifications[i].createdAt,
        isRead: true,
      );
    }
    emit(NotificationsLoadedState(notifications: _notifications, unreadCount: 0));
  }

  void _onDelete(DeleteNotificationEvent event, Emitter<NotificationsState> emit) {
    _notifications.removeWhere((n) => n.id == event.notificationId);
    final unread = _notifications.where((n) => !n.isRead).length;
    emit(NotificationsLoadedState(notifications: _notifications, unreadCount: unread));
  }
}
