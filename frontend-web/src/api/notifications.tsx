import { fetchApi, fetchApiToJson } from "../utils/apiFetch";


interface BaseNotificationDto {
    id: number;
    createdAt: Date;
}

interface User {
    id: string;
    name: string;
}

interface StudyGroup {
    id: number;
    name: string;
}

interface StudySession {
    id: number;
    title: string;
    startTime: string;
    endTime: string;
}

interface StudySessionAvailabilityRequest {
    id: number;
    studyGroup: StudyGroup;
    title: string;
}

interface NotificationPayloadMap {
    "study-group-joined": { studyGroup: StudyGroup; triggeringUser: User };
    "study-group-requested-to-join": { studyGroup: StudyGroup; triggeringUser: User };
    "study-group-left": { studyGroup: StudyGroup; triggeringUser: User };
    "study-group-accepted-invite": { studyGroup: StudyGroup; triggeringUser: User };
    "study-group-rejected-invite": { studyGroup: StudyGroup; triggeringUser: User };

    "study-group-invited": { studyGroup: StudyGroup; triggeringUser: User; targetUser: User };
    "study-group-accepted-join-request": { studyGroup: StudyGroup; triggeringUser: User; targetUser: User };
    "study-group-rejected-join-request": { studyGroup: StudyGroup; triggeringUser: User; targetUser: User };
    "study-group-removed-member": { studyGroup: StudyGroup; triggeringUser: User; targetUser: User };

    //"study-group-chat-message": { studyGroup: StudyGroup; triggeringUser: User }; // Currently not implemented in the backend

    "study-session-reminder": { studySession: StudySession };

    "study-session-scheduled": { studySession: StudySession };
    "study-session-updated": { studySession: StudySession };
    "study-session-cancelled": { studySession: StudySession };

    "study-session-availability-request-created": { studySessionAvailabilityRequest: StudySessionAvailabilityRequest };
    "study-session-availability-request-updated-entries": { studySessionAvailabilityRequest: StudySessionAvailabilityRequest; triggeringUser: User};
};

type NotificationDto = {
    [K in keyof NotificationPayloadMap]: BaseNotificationDto & {
      type: K;
      payload: NotificationPayloadMap[K];
    };
}[keyof NotificationPayloadMap];

export type Notification = NotificationDto & {
    content: string;
};


export async function getNotifications(token: string | null, userId: string): Promise<Notification[]> {
    const notifications = await fetchApiToJson<NotificationDto[]>("/notifications", token);
    const augmentedNotifications: Notification[] = notifications.map(notification => augmentNotification(notification, userId));
    return augmentedNotifications;
}

export async function markNotificationAsRead(token: string | null, notificationId: number): Promise<void> {
    await fetchApi(`/notifications/${notificationId}`, token, {
        method: "DELETE",
    });
}

function augmentNotification(notification: NotificationDto, userId: string): Notification {
    let content: string;

    switch (notification.type) {
        case "study-group-joined":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have joined the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has joined the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-requested-to-join":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have requested to join the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has requested to join the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-left":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have left the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has left the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-accepted-invite":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have accepted the invitation to join the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has accepted the invitation to join the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-rejected-invite":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have rejected the invitation to join the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has rejected the invitation to join the study group "${notification.payload.studyGroup.name}".`;
            break;
        
        case "study-group-invited":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have invited ${notification.payload.targetUser?.name} to join the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has invited ${notification.payload.targetUser?.name} to join the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-accepted-join-request":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have accepted ${notification.payload.targetUser?.name}'s request to join the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has accepted ${notification.payload.targetUser?.name}'s request to join the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-rejected-join-request":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have rejected ${notification.payload.targetUser?.name}'s request to join the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has rejected ${notification.payload.targetUser?.name}'s request to join the study group "${notification.payload.studyGroup.name}".`;
            break;
        case "study-group-removed-member":
            content = (userId === notification.payload.triggeringUser.id) ? 
                        `You have removed ${notification.payload.targetUser?.name} from the study group "${notification.payload.studyGroup.name}".` : 
                        `${notification.payload.triggeringUser.name} has removed ${notification.payload.targetUser?.name} from the study group "${notification.payload.studyGroup.name}".`;
            break;

        // case "study-group-chat-message":
        //     content = `You have a new message in the study group "${notification.payload.studyGroup.name}" from ${notification.payload.triggeringUser.name}.`;
        //     break;

        case "study-session-reminder":
            content = `Reminder: Your study session for "${notification.payload.studySession.title}" is coming up (${formatTimestamp(notification.payload.studySession.startTime)} - ${formatTimestamp(notification.payload.studySession.endTime)}).`;
            break;
        
        case "study-session-scheduled":
            content = `A new study session "${notification.payload.studySession.title}" has been scheduled (${formatTimestamp(notification.payload.studySession.startTime)} - ${formatTimestamp(notification.payload.studySession.endTime)}).`;
            break;
        case "study-session-updated":
            content = `Study session "${notification.payload.studySession.title}" has been updated (${formatTimestamp(notification.payload.studySession.startTime)} - ${formatTimestamp(notification.payload.studySession.endTime)}).`;
            break;
        case "study-session-cancelled":
            content = `Study session "${notification.payload.studySession.title}" (${formatTimestamp(notification.payload.studySession.startTime)} - ${formatTimestamp(notification.payload.studySession.endTime)}) has been cancelled.`;
            break;
        
        case "study-session-availability-request-created":
            content = `A new availability request "${notification.payload.studySessionAvailabilityRequest.title}" for study group "${notification.payload.studySessionAvailabilityRequest.studyGroup.name}" has been created.`;
            break;
        case "study-session-availability-request-updated-entries":
            content = `${notification.payload.triggeringUser.name} has updated the availability request "${notification.payload.studySessionAvailabilityRequest.title}" for study group "${notification.payload.studySessionAvailabilityRequest.studyGroup.name}" with their availability.`;
            break;
        
        default:
            content = "You have a new notification.";
    }

    return { ...notification, content };
}

function formatTimestamp(timestamp: string) {
    const date = new Date(timestamp);

    const pad = (n: number) => n.toString().padStart(2, '0');

    const hours = pad(date.getHours());
    const minutes = pad(date.getMinutes());
    const day = pad(date.getDate());
    const month = pad(date.getMonth() + 1); // Months are 0-based
    const year = date.getFullYear();

    return `${hours}:${minutes} ${day}/${month}/${year}`;
}
