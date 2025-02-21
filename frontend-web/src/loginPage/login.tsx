// src/loginPage/login.tsx
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../firebase/useAuth";
import { fetchWithToken } from "../api";
import {
    Typography,
    TextField,
    Button,
    Box,
    Link,
    Paper,
    Grid,
} from "@mui/material";
import "./login.css"; // Import the CSS file for custom styles

const Login: React.FC = () => {
    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [isRegister, setIsRegister] = useState<boolean>(false);
    const { login, signup } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log("Form submitted. Mode:", isRegister ? "Register" : "Login", "Email:", email);

        try {
            if (isRegister) {
                // 1) Register the user with Firebase
                await signup(email, password);
                alert("Registration successful! Please log in.");

                // 2) Clear form fields and switch back to login mode
                setEmail("");
                setPassword("");
                setIsRegister(false);

                // 3) Navigate to login page
                void navigate("/landing");
            } else {
                // 1) Log the user in (Firebase side)
                const userCredential = await login(email, password);
                console.log("User logged in successfully");

                // 2) Get a fresh token from the userCredential
                const newToken = await userCredential.user.getIdToken();

                // 3) Call your backend to verify the token
                const verifyResponse = await fetchWithToken<{ status: string; uid: string }>(
                    newToken,
                    "http://localhost:8080/api/auth/verify",
                    { method: "POST" }
                );
                console.log("Verify response from backend:", verifyResponse);

                // 4) Navigate to your protected page
                void navigate("/landing");
            }
        } catch (error) {
            console.error("Failed to authenticate", error);
            const errorMessage = error instanceof Error ? error.message : String(error);
            alert("Authentication failed: " + errorMessage);
        }
    };

    return (
        <Box
            sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                minHeight: "100vh", // Full viewport height
                width: "100vw", // Full viewport width
                backgroundColor: "#f0f2f5", // Light grey background
            }}
        >
            <Paper
                elevation={3}
                sx={{
                    padding: "2rem",
                    width: "100%",
                    maxWidth: "400px", // Limit the width of the white box
                    borderRadius: "8px", // Rounded corners
                    backgroundColor: "#ffffff", // White background
                }}
            >
                {/* Blackboard + StudyWise Header */}
                <Typography
                    variant="h4"
                    component="h1"
                    sx={{
                        color: "#0056b3", // Blue color
                        fontWeight: "bold",
                        textAlign: "center",
                        marginBottom: "2rem",
                    }}
                >
                    Blackboard + StudyWise
                </Typography>

                <Typography component="h2" variant="h6" align="center" sx={{ marginBottom: "1.5rem" }}>
                    {isRegister ? "Create an Account" : "Sign In"}
                </Typography>
                <Box component="form" onSubmit={(e) => void handleSubmit(e)}>
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <TextField
                                required
                                fullWidth
                                id="email"
                                label="Email Address"
                                name="email"
                                autoComplete="email"
                                value={email}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                required
                                fullWidth
                                name="password"
                                label="Password"
                                type="password"
                                id="password"
                                autoComplete="current-password"
                                value={password}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                            />
                        </Grid>
                    </Grid>
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        sx={{
                            mt: 3,
                            mb: 2,
                            backgroundColor: "#0056b3", // Blue color
                            "&:hover": {
                                backgroundColor: "#004494", // Darker blue on hover
                            },
                        }}
                    >
                        {isRegister ? "Register" : "Login"}
                    </Button>
                    <Grid container justifyContent="center">
                        <Grid item>
                            <Link
                                href="#"
                                variant="body2"
                                onClick={() => setIsRegister(!isRegister)}
                                sx={{ color: "#0056b3" }} // Blue color
                            >
                                {isRegister
                                    ? "Already have an account? Sign In"
                                    : "Don't have an account? Register"}
                            </Link>
                        </Grid>
                    </Grid>
                </Box>
            </Paper>
        </Box>
    );
};

export default Login;