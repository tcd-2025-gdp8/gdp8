// src/loginPage/login.tsx
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../firebase/useAuth";
import { fetchWithToken } from "../api";  // <-- Import your helper

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
                void navigate("/login");
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
                void navigate("/study-groups");
            }
        } catch (error) {
            console.error("Failed to authenticate", error);
            const errorMessage = error instanceof Error ? error.message : String(error);
            alert("Authentication failed: " + errorMessage);
        }
    };

    return (
        <div>
            <h2>{isRegister ? "Register" : "Login"}</h2>
            <form onSubmit={(e) => void handleSubmit(e)}>
                <input
                    type="email"
                    value={email}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                    placeholder="Email"
                />
                <input
                    type="password"
                    value={password}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                    placeholder="Password"
                />
                <button type="submit">{isRegister ? "Register" : "Login"}</button>
            </form>
            <p>
                {isRegister ? "Already have an account?" : "Don't have an account?"}{" "}
                <button onClick={() => setIsRegister(!isRegister)}>
                    {isRegister ? "Login" : "Register"}
                </button>
            </p>
        </div>
    );
};

export default Login;
