// src/loginPage/login.tsx
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../firebase/useAuth";

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
                // Register the user
                await signup(email, password);
                // Notify the user of successful registration
                alert("Registration successful! Please log in.");
                // Optionally clear the form fields
                setEmail("");
                setPassword("");
                // Switch back to login mode
                setIsRegister(false);
                // Navigate back to the login page
                void navigate("/login");
            } else {
                // Log the user in
                await login(email, password);
                console.log("User logged in successfully");
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
