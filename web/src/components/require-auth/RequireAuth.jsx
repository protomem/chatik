import { Outlet, useLocation, Navigate } from "react-router-dom";
import { selectIsLoggedIn } from "../../store/auth/auth.selectors";
import { useSelector } from "react-redux";

export function RequireAuth() {
  const isLoggedIn = useSelector((state) => selectIsLoggedIn(state));
  const location = useLocation();

  return isLoggedIn ? (
    <Outlet />
  ) : (
    <Navigate to="/auth" state={{ from: location }} replace />
  );
}
