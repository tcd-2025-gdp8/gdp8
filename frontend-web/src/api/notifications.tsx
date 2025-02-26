import { fetchApi, fetchApiToJson } from "../utils/apiFetch";


export type NotificationType =
    | "study-group-joined"
    | "study-group-requested-to-join"
    | "study-group-left"
    | "study-group-accepted-invite"
    | "study-group-rejected-invite"
    | "study-group-invited"
    | "study-group-accepted-join-request"
    | "study-group-rejected-join-request"
    | "study-group-removed-member"
    | "study-group-chat-message";


export interface NotificationDto {
    id: number;
    type: NotificationType;
    triggeringUser: {
        id: string;
        name: string;
    };
    targetUser: {
        id: string;
        name: string;
    } | null;
    studyGroup: {
        id: number;
        name: string;
    };
    messageId: number | null;
    createdAt: Date;
}

export interface Notification extends NotificationDto {
    content: string;
}


export async function getNotifications(token: string | null, userId: string): Promise<Notification[]> {
    const notifications = await fetchApiToJson<NotificationDto[]>("/notifications", token);
    console.log(notifications);
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
            content = (userId === notification.triggeringUser.id) ? 
                        `You have joined the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has joined the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-requested-to-join":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have requested to join the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has requested to join the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-left":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have left the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has left the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-accepted-invite":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have accepted the invitation to join the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has accepted the invitation to join the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-rejected-invite":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have rejected the invitation to join the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has rejected the invitation to join the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-invited":
           content = (userId === notification.triggeringUser.id) ? 
                       `You have invited ${notification.targetUser?.name} to join the study group "${notification.studyGroup.name}".` : 
                       `${notification.triggeringUser.name} has invited ${notification.targetUser?.name} to join the study group "${notification.studyGroup.name}".`;
           break;
        case "study-group-accepted-join-request":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have accepted ${notification.targetUser?.name}'s request to join the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has accepted ${notification.targetUser?.name}'s request to join the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-rejected-join-request":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have rejected ${notification.targetUser?.name}'s request to join the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has rejected ${notification.targetUser?.name}'s request to join the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-removed-member":
            content = (userId === notification.triggeringUser.id) ? 
                        `You have removed ${notification.targetUser?.name} from the study group "${notification.studyGroup.name}".` : 
                        `${notification.triggeringUser.name} has removed ${notification.targetUser?.name} from the study group "${notification.studyGroup.name}".`;
            break;
        case "study-group-chat-message":
            content = `You have a new message in the study group "${notification.studyGroup.name}" from ${notification.triggeringUser.name}.`;
            break;
        default:
            content = "You have a new notification.";
    }

    return { ...notification, content };
}

