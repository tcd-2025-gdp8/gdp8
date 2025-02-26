import { Navigate, Outlet } from "react-router-dom";
import {
    Box
} from "@mui/material";

import { useAuth } from "./auth/useAuth";
import Sidebar from "./components/Sidebar";
import CustomAppBar from "./components/CustomAppBar";

export default function Layout() {
    const { token, user } = useAuth();

    if(!token || !user) {
        return <Navigate to="/login" />;
    }
    
    return (
        <Box sx={{ display: "flex", width: "100vw", height: "100vh" }}>
            <Sidebar />
            <Box sx={{ flexGrow: 1 }}>
                <CustomAppBar />

                {/* Main Content Area (Below the AppBar) */}
                <Box 
                    sx={{
                        padding: "2rem",
                        marginTop: "64px", // Height of the AppBar
                        height: "calc(100vh - 64px)",
                        boxSizing: "border-box",
                        overflow: "auto",
                    }}
                >
                    <Outlet />
                </Box>
            </Box>
        </Box>
    );
}
