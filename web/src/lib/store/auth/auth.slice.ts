import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { IUser } from "../../domain/entity/entities";

interface AuthState {
  token: string;
  user: IUser | null;
}

const initialState: AuthState = {
  token: "",
  user: null,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload }: PayloadAction<AuthState>) => {
      if (!payload.token || !payload.user) return;

      state.token = payload.token;
      state.user = payload.user;
    },

    clearCredentials: (state) => {
      state.token = "";
      state.user = null;
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
