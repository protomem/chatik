import { BrowserRouter, Route, Routes } from "react-router-dom";
import { NotFound } from "../pages/not-found/NotFound";
import { Auth } from "../pages/auth/Auth";
import { RequireAuth } from "../components/require-auth/RequireAuth";
import { Chat } from "../pages/chat/Chat";

export function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<RequireAuth />}>
          <Route index element={<Chat />} />
        </Route>

        <Route path="/auth" element={<Auth />} />

        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}
