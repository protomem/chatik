import React from "react";
import { useAppSelector } from "../../../feature/hooks";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";

interface SSEProps {
  children: React.ReactNode;
}

export const SSE: React.FC<SSEProps> = ({ children }) => {
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const sse = new EventSource(
    `http://${
      import.meta.env.VITE_API_URL
    }/api/v1/stream/sse?token=${accessToken}`,
  );

  return <>{children}</>;
};
