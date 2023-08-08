import React from "react";

interface WSProps {
  children: React.ReactNode;
}

export const WS: React.FC<WSProps> = ({ children }) => {
  return <>{children}</>;
};
