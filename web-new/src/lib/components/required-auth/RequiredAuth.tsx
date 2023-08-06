import React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAppSelector } from "../../feature/hooks";
import { selectIsAuthenticated } from "../../feature/auth/auth.selectors";

export const RequiredAuth: React.FC = () => {
  const location = useLocation();
  const isAuth = useAppSelector((state) => selectIsAuthenticated(state));

  return isAuth ? (
    <Outlet />
  ) : (
    <Navigate to="/sign-in" state={{ from: location }} replace />
  );
};
