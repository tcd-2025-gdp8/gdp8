// src/firebase/useAuth.ts
import { useState, useEffect } from "react";
import { auth } from "./firebase";
import {
    User,
    UserCredential,
    createUserWithEmailAndPassword,
    signInWithEmailAndPassword,
    signOut,
    onAuthStateChanged,
    getIdToken,
} from "firebase/auth";

export const useAuth = () => {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
            if (currentUser) {
                setUser(currentUser);
                getIdToken(currentUser)
                    .then((idToken) => {
                        console.log("Token received from onAuthStateChanged:", idToken);
                        setToken(idToken);
                    })
                    .catch((error) => {
                        console.error("Error getting ID token:", error);
                        setToken(null);
                    });
            } else {
                setUser(null);
                setToken(null);
            }
        });
        return unsubscribe;
    }, []);

    const signup = async (email: string, password: string) => {
        const userCredential: UserCredential = await createUserWithEmailAndPassword(auth, email, password);
        const idToken = await getIdToken(userCredential.user);
        console.log("Token received on signup:", idToken);
        setToken(idToken);
        return userCredential;
    };

    const login = async (email: string, password: string) => {
        const userCredential: UserCredential = await signInWithEmailAndPassword(auth, email, password);
        const idToken = await getIdToken(userCredential.user);
        console.log("Token received on login:", idToken);
        setToken(idToken);
        return userCredential;
    };

    const logout = async () => {
        await signOut(auth);
        console.log("User logged out. Token cleared.");
        setToken(null);
    };

    return { user, token, signup, login, logout };
};
