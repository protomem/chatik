import { Navigate, Outlet, useLocation } from "react-router-dom";
import { selectIsAuth } from "../../store/auth/auth.selectors";
import { useAppSelector } from "../../store/hooks";

export function RequireAuth() {
  const isAuth = useAppSelector((state) => selectIsAuth(state));
  const location = useLocation();

  return isAuth ? (
    <Outlet />
  ) : (
    <Navigate to="/auth" state={{ from: location }} replace />
  );
}
