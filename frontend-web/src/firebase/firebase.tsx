// src/firebase/firebase.ts
import { initializeApp, FirebaseApp, FirebaseOptions } from "firebase/app";
import { getAnalytics, Analytics } from "firebase/analytics";
import { getAuth, Auth } from "firebase/auth";

const firebaseConfig: FirebaseOptions = {
    apiKey: "AIzaSyCIbysnKfrvwt8gCHi2jydh6iYftPsGYEA",
    authDomain: "studygroups-3380d.firebaseapp.com",
    projectId: "studygroups-3380d",
    storageBucket: "studygroups-3380d.firebasestorage.app",
    messagingSenderId: "247573450343",
    appId: "1:247573450343:web:653066a5bb355c8acb0ec2",
    measurementId: "G-Z6KLPYPV46",
};

// Initialize Firebase
const app: FirebaseApp = initializeApp(firebaseConfig);
const analytics: Analytics = getAnalytics(app);
const auth: Auth = getAuth(app);

export { app, analytics, auth };
