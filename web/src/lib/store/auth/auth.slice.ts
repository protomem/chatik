import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { IUser } from "../../entity/entities";

interface AuthState {
  token: string;
  user: IUser | null;
}

const getTokenFromLocalStorage = (): string => {
  const token = localStorage.getItem("token");
  if (!token) return "";

  return token;
};

const getUserFromLocalStorage = (): IUser | null => {
  const user = localStorage.getItem("user");
  if (!user) return null;

  try {
    return JSON.parse(user) as IUser;
  } catch {
    return null;
  }
};

const saveTokenToLocalStorage = (token: string) => {
  localStorage.setItem("token", token);
};

const saveUserToLocalStorage = (user: IUser) => {
  localStorage.setItem("user", JSON.stringify(user));
};

const removeTokenFromLocalStorage = () => {
  localStorage.removeItem("token");
};

const removeUserFromLocalStorage = () => {
  localStorage.removeItem("user");
};

const initialState: AuthState = {
  token: getTokenFromLocalStorage(),
  user: getUserFromLocalStorage(),
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload }: PayloadAction<AuthState>) => {
      if (!payload.token || !payload.user) return;

      state.token = payload.token;
      state.user = payload.user;

      saveTokenToLocalStorage(payload.token);
      saveUserToLocalStorage(payload.user);
    },

    clearCredentials: (state) => {
      state.token = "";
      state.user = null;

      removeTokenFromLocalStorage();
      removeUserFromLocalStorage();
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
