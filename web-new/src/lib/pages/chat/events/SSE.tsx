import React from "react";

interface SSEProps {
  children: React.ReactNode;
}

export const SSE: React.FC<SSEProps> = ({ children }) => {
  return <>{children}</>;
};
