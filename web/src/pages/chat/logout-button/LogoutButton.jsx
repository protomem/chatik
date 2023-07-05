import { Button } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { authActions } from "../../../store/auth/auth.slice";

export function LogoutButton() {
  const nav = useNavigate();

  const dispatch = useDispatch();

  const handleClick = () => {
    dispatch(authActions.logout());
    nav("/auth", { replace: true });
  };

  return (
    <Button onClick={handleClick} variant="outlined" size="medium">
      Logout
    </Button>
  );
}
