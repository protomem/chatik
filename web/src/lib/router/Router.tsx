import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { NotFoundPage } from "../pages/not-found/NotFoundPage";
import { ChatPage } from "../pages/chat/ChatPage";
import { SignInPage } from "../pages/auth/SignInPage";
import { RequiredAuth } from "../components/required-auth/RequiredAuth";
import { SignUpPage } from "../pages/auth/SingUpPage";

export const Router: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/sign-in" element={<SignInPage />} />
        <Route path="/sign-up" element={<SignUpPage />} />

        <Route path="/" element={<RequiredAuth />}>
          <Route index element={<ChatPage />} />
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  );
};
