import { Button } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useAppDispatch } from "../../../store/hooks";
import { authActions } from "../../../store/auth/auth.slice";

export function LogoutButton() {
  const nav = useNavigate();
  const dispatch = useAppDispatch();

  const handleLogout = () => {
    dispatch(authActions.clearCredentials());
    nav("/", { replace: true });
  };

  return (
    <Button variant="outlined" size="large" onClick={handleLogout}>
      Logout
    </Button>
  );
}
