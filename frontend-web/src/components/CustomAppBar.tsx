import { useEffect, useState } from "react";
import {
    Box,
    Typography,
    Drawer,
    IconButton,
    AppBar,
    Toolbar,
    Badge,
} from "@mui/material";
import NotificationsIcon from "@mui/icons-material/Notifications";
import CloseIcon from "@mui/icons-material/Close";

import { Notification, getNotifications, markNotificationAsRead } from "../api/notifications";
import { useAuth } from "../auth/useAuth";


export default function CustomAppBar() {
    const { token, user } = useAuth();

    const [openNotifications, setOpenNotifications] = useState(false);
    const [notifications, setNotifications] = useState<Notification[]>([]);


    useEffect(() => {
        const fetchNotifications = async () => {
            if(!user) return;
            try {
                const fetchedNotifications = await getNotifications(token, user.uid);
                fetchedNotifications.reverse();
                setNotifications(fetchedNotifications);
            } catch (error) {
                console.error("Error fetching notifications:", error);
                alert("Error fetching notifications");
            }
        };
        void fetchNotifications();
    }, [token, user]);

    const handleDeleteNotification = (notificationId: number) => {
        void (async () => {
            try {
                await markNotificationAsRead(token, notificationId);
                setNotifications((prev) =>
                    prev.filter((notification) => notification.id !== notificationId)
                );
            } catch (error) {
                console.error("Error deleting notification:", error);
                alert("Error deleting notification");
            }
        })();
    };


    return (
        <>
            <AppBar
                position="fixed"
                sx={{
                    backgroundColor: "#ffffff",
                    color: "#000000",
                    zIndex: (theme) => theme.zIndex.drawer - 1,
                    width: "calc(100% - 240px)",
                    marginLeft: "240px",
                }}
            >
                <Toolbar sx={{ flexDirection: "row", alignItems: "center" }}>
                    <Box sx={{ flexGrow: 1, display: "flex", alignItems: "center" }}>
                        <Typography variant="h6">
                            PeerSphere
                        </Typography>
                    </Box>
                    {user && (
                        <Typography variant="subtitle1" sx={{ mr: 2 }}>
                            {user.email}
                        </Typography>
                    )}
                    <IconButton color="inherit" onClick={() => setOpenNotifications(true)}>
                        <Badge badgeContent={notifications.length} color="error">
                            <NotificationsIcon />
                        </Badge>
                    </IconButton>
                </Toolbar>
            </AppBar>

            {/* Notifications Drawer */}
            <Drawer
                anchor="right"
                open={openNotifications}
                onClose={() => setOpenNotifications(false)}
                sx={{
                    "& .MuiDrawer-paper": {
                        width: 350,
                        backgroundColor: "#ffffff",
                        borderLeft: "1px solid #e0e0e0",
                    },
                }}
            >
                <Box sx={{ padding: "1.5rem" }}>
                    <Box
                        sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "center",
                            marginBottom: "1.5rem",
                            borderBottom: "1px solid #e0e0e0",
                            paddingBottom: "1rem"
                        }}
                    >
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>Notifications</Typography>
                        <IconButton onClick={() => setOpenNotifications(false)}>
                            <CloseIcon />
                        </IconButton>
                    </Box>
                    {notifications.length === 0 ? (
                        <Box sx={{ 
                            textAlign: "center", 
                            py: 4,
                            color: "text.secondary"
                        }}>
                            <NotificationsIcon sx={{ fontSize: 48, opacity: 0.5, mb: 2 }} />
                            <Typography variant="body1">
                                No notifications
                            </Typography>
                        </Box>
                    ) : (
                        notifications.map((notification) => (
                            <Box
                                key={notification.id}
                                sx={{
                                    padding: "1rem",
                                    borderRadius: "8px",
                                    backgroundColor: "#ffffff",
                                    boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
                                    marginBottom: "0.75rem",
                                    display: "flex",
                                    justifyContent: "space-between",
                                    alignItems: "center",
                                    transition: "all 0.2s",
                                    "&:hover": {
                                        boxShadow: "0 2px 5px rgba(0,0,0,0.15)",
                                    }
                                }}
                            >
                                <Typography sx={{ pr: 2 }}>{notification.content}</Typography>
                                <IconButton
                                    size="small"
                                    onClick={() => handleDeleteNotification(notification.id)}
                                    sx={{ 
                                        opacity: 0.6,
                                        "&:hover": { opacity: 1 }
                                    }}
                                >
                                    <CloseIcon fontSize="small" />
                                </IconButton>
                            </Box>
                        ))
                    )}
                </Box>
            </Drawer>
        </>
    )
}
