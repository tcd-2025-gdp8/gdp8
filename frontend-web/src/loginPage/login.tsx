// Login.tsx
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../firebase/authContext";

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
                await signup(email, password);
                // Optionally you can auto-login or redirect after successful registration
                navigate("/study-groups");
            } else {
                await login(email, password);
                console.log("User logged in successfully");
                navigate("/study-groups");
            }
        } catch (error) {
            console.error("Failed to authenticate", error);
        }
    };

    return (
        <div>
            <h2>{isRegister ? "Register" : "Login"}</h2>
            <form onSubmit={handleSubmit}>
                <input
                    type="email"
                    value={email}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setEmail(e.target.value)
                    }
                    placeholder="Email"
                />
                <input
                    type="password"
                    value={password}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setPassword(e.target.value)
                    }
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
