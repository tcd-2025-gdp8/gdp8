// src/firebase/authContext.ts
import { createContext } from "react";
import { User, UserCredential } from "firebase/auth";

export interface AuthContextType {
    user: User | null;
    token: string | null;
    signup: (email: string, password: string) => Promise<UserCredential>;
    login: (email: string, password: string) => Promise<UserCredential>;
    logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);