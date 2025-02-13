// src/firebase/authContextProvider.tsx
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

export function AuthProvider({ children }: AuthProviderProps) {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
            const handleAuthChange = async () => {
                setUser(currentUser);
                if (currentUser) {
                    try {
                        const idToken = await getIdToken(currentUser); // Await the promise
                        setToken(idToken);
                    } catch (error) {
                        console.error("Failed to get ID token:", error);
                        setToken(null);
                    }
                } else {
                    setToken(null);
                }
                setLoading(false);
            };

            void handleAuthChange(); // Explicitly mark the promise as ignored
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
}