import React from "react";
import { useAppSelector } from "../../../feature/hooks";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";

interface WSProps {
  children: React.ReactNode;
}

export const WS: React.FC<WSProps> = ({ children }) => {
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const ws = new WebSocket(
    `ws://${
      import.meta.env.VITE_API_URL
    }/api/v1/stream/ws?token=${accessToken}`,
  );

  return <>{children}</>;
};
