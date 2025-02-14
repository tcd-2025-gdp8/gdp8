// src/firebase/authProvider.tsx
import React, { useState, useEffect, ReactNode } from "react";
import { auth } from "./firebase";
import { AuthContext } from "./authContext";
import {
    createUserWithEmailAndPassword,
    signInWithEmailAndPassword,
    signOut,
    onAuthStateChanged,
    getIdToken,
    User,
} from "firebase/auth";

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
            setUser(currentUser);
            if (currentUser) {
                getIdToken(currentUser)
                    .then((idToken) => {
                        setToken(idToken);
                    })
                    .catch((error) => {
                        console.error("Error getting ID token:", error);
                        setToken(null);
                    });
            } else {
                setToken(null);
            }
            setLoading(false);
        });
        return unsubscribe;
    }, []);

    const signup = (email: string, password: string) => {
        return createUserWithEmailAndPassword(auth, email, password);
    };

    const login = (email: string, password: string) => {
        return signInWithEmailAndPassword(auth, email, password);
    };

    const logout = () => {
        return signOut(auth);
    };

    return (
        <AuthContext.Provider value={{ user, token, signup, login, logout }}>
            {!loading && children}
        </AuthContext.Provider>
    );
};
