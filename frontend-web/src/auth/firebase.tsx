// src/firebase/firebase.tsx
import { initializeApp, FirebaseApp, FirebaseOptions } from "firebase/app";
import { getAnalytics, Analytics } from "firebase/analytics";
import {
    getAuth,
    setPersistence,
    browserSessionPersistence,
    Auth,
} from "firebase/auth";

const firebaseConfig: FirebaseOptions = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY as string,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN as string,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID as string,
    storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET as string,
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID as string,
    appId: import.meta.env.VITE_FIREBASE_APP_ID as string,
    measurementId: import.meta.env.VITE_FIREBASE_MEASUREMENT_ID as string,
};

// Initialize Firebase
const app: FirebaseApp = initializeApp(firebaseConfig);
const analytics: Analytics = getAnalytics(app);
const auth: Auth = getAuth(app);

// Set persistence to session so the user is signed out when the browser or tab is closed.
setPersistence(auth, browserSessionPersistence)
    .then(() => {
        console.log("Session persistence set successfully.");
    })
    .catch((error) => {
        console.error("Error setting session persistence:", error);
    });

export { app, analytics, auth };
