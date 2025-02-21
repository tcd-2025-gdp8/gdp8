import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth/useAuth";
import {
    Typography,
    TextField,
    Button,
    Box,
    Link,
    Paper,
    Grid,
} from "@mui/material";

import "./login.css";

const Login: React.FC = () => {
    const [isRegister, setIsRegister] = useState<boolean>(false);
    const [firstName, setFirstName] = useState<string>("");
    const [lastName, setLastName] = useState<string>("");
    const { login, signup } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        const formData = new FormData(e.currentTarget);
        const email = formData.get('email') as string;
        const password = formData.get('password') as string;

        try {
            if (isRegister) {
                await signup(email, password);
                const userCredential = await login(email, password);
                const newToken = await userCredential.user.getIdToken();
                const firebaseUID = userCredential.user.uid;

                const createUserResponse = await fetchWithToken<{ id: string }>(
                    newToken,
                    "http://localhost:8080/api/user",
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                            id: firebaseUID,
                            name: firstName + " " + lastName,
                            modules: []
                        }),
                    }
                );
                console.log(createUserResponse);
                alert("Registration successful! Please log in.");

                e.currentTarget.reset();
                setFirstName("");
                setLastName("");
                setIsRegister(false);

                // 3) Navigate to login page
                void navigate("/landing");
            } else {
                await login(email, password);

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
                            {isRegister && (
                            <>
                                <Grid item xs={6}>
                                    <TextField
                                        required
                                        fullWidth
                                        id="firstName"
                                        label="First Name"
                                        name="firstName"
                                        value={firstName}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFirstName(e.target.value)}
                                    />
                                </Grid>
                                <Grid item xs={6}>
                                    <TextField
                                        required
                                        fullWidth
                                        id="lastName"
                                        label="Last Name"
                                        name="lastName"
                                        value={lastName}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setLastName(e.target.value)}
                                    />
                                </Grid>
                            </>
                        )}
                        <Grid item xs={12}>
                            <TextField
                                required
                                fullWidth
                                id="email"
                                label="Email Address"
                                name="email"
                                autoComplete="email"
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
