import { BrowserRouter, Routes, Route } from "react-router-dom";
import { NotFound } from "../pages/not-found/NotFound";
import { Chat } from "../pages/chat/Chat";
import { Auth } from "../pages/auth/Auth";

export function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Chat />} />
        <Route path="/auth" element={<Auth />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}
